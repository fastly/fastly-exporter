package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "dev"

func main() {
	fs := flag.NewFlagSet("fastly-exporter", flag.ExitOnError)
	var (
		token      = fs.String("token", "", "Fastly API token (required; also via FASTLY_API_TOKEN)")
		serviceIDs = stringslice{}
		addr       = fs.String("endpoint", "http://127.0.0.1:8080/metrics", "Prometheus /metrics endpoint")
		namespace  = fs.String("namespace", "fastly", "Prometheus namespace")
		subsystem  = fs.String("subsystem", "rt", "Prometheus subsystem")
		debug      = fs.Bool("debug", false, "Log debug information")
	)
	fs.Var(&serviceIDs, "service", "Specific Fastly service ID (optional, repeatable)")
	fs.Usage = usageFor(fs, "fastly-exporter [flags]")
	fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		loglevel := level.AllowInfo()
		if *debug {
			loglevel = level.AllowDebug()
		}
		logger = level.NewFilter(logger, loglevel)
	}

	if *token == "" {
		if *token = os.Getenv("FASTLY_API_TOKEN"); *token == "" {
			level.Error(logger).Log("err", "-token or FASTLY_API_TOKEN is required")
			os.Exit(1)
		}
	}

	var metrics prometheusMetrics
	{
		metrics.register(*namespace, *subsystem)
	}

	var promURL *url.URL
	{
		var err error
		promURL, err = url.Parse(*addr)
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
		level.Info(logger).Log("prometheus_addr", promURL.Host, "path", promURL.Path, "namespace", *namespace, "subsystem", *subsystem)
	}

	var cache *nameCache
	{
		cache = newNameCache()
	}

	var manager *monitorManager
	{
		var (
			postprocess = func() {}                          // only used for tests
			rtClient    = &http.Client{Timeout: time.Minute} // rt.fastly.com blocks awhile by design
		)
		manager = newMonitorManager(rtClient, *token, cache, metrics, postprocess, log.With(logger, "component", "monitors"))
	}

	var apiClient *http.Client
	{
		apiClient = &http.Client{
			Timeout: 10 * time.Second, // api.fastly.com should be fast
		}
	}

	var queryer *serviceQueryer
	{
		queryer = newServiceQueryer(*token, serviceIDs, cache, manager)
		if err := queryer.refresh(apiClient); err != nil { // first refresh must succeed
			level.Error(logger).Log("during", "initial service refresh", "err", err)
			os.Exit(1)
		}
	}

	var g run.Group
	{
		// Every minute, query Fastly for new services and their names.
		// Update our name cache and managed monitors accordingly.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(time.Minute)
		)
		g.Add(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case <-ticker.C:
					if err := queryer.refresh(apiClient); err != nil {
						level.Warn(logger).Log("during", "service refresh", "err", err)
					}
				}
			}
		}, func(error) {
			ticker.Stop()
			cancel()
		})
	}
	{
		// A pseudo-actor that exists only to tear down all managed monitors.
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			<-ctx.Done()
			manager.stopAll()
			return ctx.Err()
		}, func(error) {
			cancel()
		})
	}
	{
		// Serve Prometheus metrics over HTTP.
		mux := http.NewServeMux()
		mux.Handle(promURL.Path, promhttp.Handler())
		server := http.Server{
			Addr:    promURL.Host,
			Handler: mux,
		}
		g.Add(func() error {
			return server.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	}
	{
		// Catch ctrl-C.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			c           = make(chan os.Signal, 1)
		)
		signal.Notify(c, os.Interrupt)
		g.Add(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			}
		}, func(error) {
			cancel()
		})
	}
	level.Info(logger).Log("exit", g.Run())
}

type stringslice []string

func (ss *stringslice) Set(s string) error {
	(*ss) = append(*ss, s)
	return nil
}

func (ss *stringslice) String() string {
	if len(*ss) <= 0 {
		return "..."
	}
	return strings.Join(*ss, ", ")
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			def := f.DefValue
			if def == "" {
				def = "..."
			}
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, def, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "VERSION\n")
		fmt.Fprintf(os.Stderr, "  %s\n", version)
		fmt.Fprintf(os.Stderr, "\n")
	}
}
