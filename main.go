package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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
		token     = fs.String("token", "", "Fastly API token")
		service   = fs.String("service", "", "Fastly service")
		addr      = fs.String("endpoint", "http://127.0.0.1:8080/metrics", "Prometheus /metrics endpoint")
		namespace = fs.String("namespace", "", "Prometheus namespace")
		subsystem = fs.String("subsystem", "", "Prometheus subsystem")
		debug     = fs.Bool("debug", false, "log debug information")
	)
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
		level.Error(logger).Log("err", "-token is required")
		os.Exit(1)
	}
	if *service == "" {
		level.Error(logger).Log("err", "-service is required")
		os.Exit(1)
	}
	level.Info(logger).Log("fastly_service", *service)

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

	var m prometheusMetrics
	{
		m.register(*namespace, *subsystem)
	}

	var g run.Group
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return queryLoop(ctx, *token, *service, &m, log.With(logger, "query", "rt.fastly.com"))
		}, func(error) {
			cancel()
		})
	}
	{
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
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
		})
	}
	level.Info(logger).Log("exit", g.Run())
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
