package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/peterbourgon/fastly-exporter/pkg/rt"
	"github.com/peterbourgon/usage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var programVersion = "dev"

func main() {
	fs := flag.NewFlagSet("fastly-exporter", flag.ExitOnError)
	var (
		token        = fs.String("token", "", "Fastly API token (required; also via FASTLY_API_TOKEN)")
		addr         = fs.String("endpoint", "http://127.0.0.1:8080/metrics", "Prometheus /metrics endpoint")
		namespace    = fs.String("namespace", "fastly", "Prometheus namespace")
		subsystem    = fs.String("subsystem", "rt", "Prometheus subsystem")
		serviceIDs   = stringslice{}
		nameRegexStr = fs.String("name-regex", "", "if provided, only export services whose names match this regex")
		shard        = fs.String("shard", "", "if provided, only export services whose hashed IDs modulo m equal (n-1) (format 'n/m')")
		apiRefresh   = fs.Duration("api-refresh", time.Minute, "how often to poll api.fastly.com for updated service metadata")
		debug        = fs.Bool("debug", false, "Log debug information")
		versionFlag  = fs.Bool("version", false, "print version information and exit")
	)
	fs.Var(&serviceIDs, "service", "if provided, only export services with this service ID (repeatable)")
	fs.Usage = usage.For(fs, "fastly-exporter [flags]")
	fs.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "fastly-exporter v%s\n", programVersion)
		os.Exit(0)
	}

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

	{
		if *apiRefresh < 30*time.Second {
			level.Warn(logger).Log("msg", "-api-refresh cannot be shorter than 15s; setting it to 15s")
			*apiRefresh = 15 * time.Second
		}
		if *apiRefresh > 10*time.Minute {
			level.Warn(logger).Log("msg", "-api-refresh cannot be longer than 10m; setting it to 10m")
			*apiRefresh = 10 * time.Minute
		}
	}

	var nameRegex *regexp.Regexp
	{
		if *nameRegexStr != "" {
			var err error
			if nameRegex, err = regexp.Compile(*nameRegexStr); err != nil {
				level.Error(logger).Log("err", "-service-name-regex invalid", "msg", err)
				os.Exit(1)
			}
		}
	}

	var shardN, shardM uint64
	{
		if *shard != "" {
			toks := strings.SplitN(*shard, "/", 2)
			if len(toks) != 2 {
				level.Error(logger).Log("err", "-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			var err error
			shardN, err = strconv.ParseUint(toks[0], 10, 64)
			if err != nil {
				level.Error(logger).Log("err", "-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			if shardN <= 0 {
				level.Error(logger).Log("err", "first part of -shard flag should be greater than zero")
				os.Exit(1)
			}
			shardM, err = strconv.ParseUint(toks[1], 10, 64)
			if err != nil {
				level.Error(logger).Log("err", "-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			if shardN > shardM {
				level.Error(logger).Log("err", fmt.Sprintf("-shard with n=%d m=%d is invalid: n must be less than or equal to m", shardN, shardM))
				os.Exit(1)
			}
		}
	}

	var registry *prometheus.Registry
	{
		registry = prometheus.NewRegistry()
	}

	var metrics *prom.Metrics
	{
		var err error
		metrics, err = prom.NewMetrics(*namespace, *subsystem, registry)
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
	}

	var apiLogger log.Logger
	{
		apiLogger = log.With(logger, "component", "api.fastly.com")
	}

	var apiCacheOptions []api.CacheOption
	{
		apiCacheOptions = append(apiCacheOptions, api.WithLogger(apiLogger))

		if len(serviceIDs) > 0 {
			level.Info(apiLogger).Log("filtering_on", "explicit service IDs", "count", len(serviceIDs))
			apiCacheOptions = append(apiCacheOptions, api.WithExplicitServiceIDs(serviceIDs...))
		}

		if nameRegex != nil {
			level.Info(apiLogger).Log("filtering_on", "service name regex", "regex", *nameRegexStr)
			apiCacheOptions = append(apiCacheOptions, api.WithNameMatching(nameRegex))
		}

		if shardM > 0 {
			level.Info(apiLogger).Log("filtering_on", "shard allocation", "shard", *shard, "n", shardN, "m", shardM)
			apiCacheOptions = append(apiCacheOptions, api.WithShard(shardN, shardM))
		}
	}

	var apiClient *http.Client
	{
		apiClient = &http.Client{
			Timeout: 15 * time.Second, // api.fastly.com should be fast
		}
	}

	var cache *api.Cache
	{
		cache = api.NewCache(*token, apiCacheOptions...)

		if err := cache.Refresh(apiClient); err != nil {
			level.Error(apiLogger).Log("during", "initial service refresh", "err", err)
			os.Exit(1)
		}
	}

	var rtLogger log.Logger
	{
		rtLogger = log.With(logger, "component", "rt.fastly.com")
	}

	var manager *rt.Manager
	{
		rtClient := &http.Client{
			Timeout: 45 * time.Second, // rt.fastly.com may block for up to ~30 seconds
		}
		subscriberOptions := []rt.SubscriberOption{
			rt.WithLogger(rtLogger),
			rt.WithMetadataProvider(cache),
			rt.WithUserAgent(`Fastly-Exporter (` + programVersion + `)`),
		}

		manager = rt.NewManager(cache, rtClient, *token, metrics, subscriberOptions, rtLogger)
		manager.Refresh() // populate initial subscribers, based on the initial cache refresh
	}

	var g run.Group
	{
		// Every *apiRefresh, ask the api.Cache to refresh the set of services
		// we should be exporting data for. Then, ask the rt.Manager to refresh
		// its set of rt.Subscribers, based on those latest services.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(*apiRefresh)
		)
		g.Add(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case <-ticker.C:
					err := cache.Refresh(apiClient)
					if err != nil {
						level.Warn(apiLogger).Log("during", "service refresh", "err", err, "msg", "the set of exported services and their metadata may be stale")
					}
					manager.Refresh() // safe to do with stale data in the cache
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
			manager.StopAll()
			return ctx.Err()
		}, func(error) {
			cancel()
		})
	}
	{
		// Serve Prometheus metrics over HTTP.
		mux := http.NewServeMux()
		mux.Handle(promURL.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		server := http.Server{
			Addr:    promURL.Host,
			Handler: mux,
		}
		g.Add(func() error {
			return server.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			level.Debug(logger).Log("msg", "shutting down Prometheus HTTP server")
			server.Shutdown(ctx)
		})
	}
	{
		// Catch ctrl-C.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			sigchan     = make(chan os.Signal, 1)
		)
		signal.Notify(sigchan, os.Interrupt)
		g.Add(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-sigchan:
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