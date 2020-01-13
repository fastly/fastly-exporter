package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
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
		token            = fs.String("token", "", "Fastly API token (required; also via FASTLY_API_TOKEN)")
		addr             = fs.String("endpoint", "http://127.0.0.1:8080/metrics", "Prometheus /metrics endpoint")
		namespace        = fs.String("namespace", "fastly", "Prometheus namespace")
		subsystem        = fs.String("subsystem", "rt", "Prometheus subsystem")
		serviceShard     = fs.String("service-shard", "", "if set, only include services whose hashed IDs modulo m equal n-1 (format 'n/m')")
		serviceIDs       = stringslice{}
		serviceWhitelist = stringslice{}
		serviceBlacklist = stringslice{}
		metricWhitelist  = stringslice{}
		metricBlacklist  = stringslice{}

		apiRefresh  = fs.Duration("api-refresh", time.Minute, "how often to poll api.fastly.com for updated service metadata (15s–10m)")
		apiTimeout  = fs.Duration("api-timeout", 15*time.Second, "HTTP client timeout for api.fastly.com requests (5–60s)")
		rtTimeout   = fs.Duration("rt-timeout", 45*time.Second, "HTTP client timeout for rt.fastly.com requests (45–120s)")
		debug       = fs.Bool("debug", false, "Log debug information")
		versionFlag = fs.Bool("version", false, "print version information and exit")
	)
	fs.Var(&serviceIDs, "service", "if set, only include this service ID (repeatable)")
	fs.Var(&serviceWhitelist, "service-whitelist", "if set, only include services whose names match this regex (repeatable)")
	fs.Var(&serviceBlacklist, "service-blacklist", "if set, don't include services whose names match this regex (repeatable)")
	fs.Var(&metricWhitelist, "metric-whitelist", "if set, only export metrics whose names match this regex (repeatable)")
	fs.Var(&metricBlacklist, "metric-blacklist", "if set, don't export metrics whose names match this regex (repeatable)")
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
		if *apiRefresh < 15*time.Second {
			level.Warn(logger).Log("msg", "-api-refresh cannot be shorter than 15s; setting it to 15s")
			*apiRefresh = 15 * time.Second
		}
		if *apiRefresh > 10*time.Minute {
			level.Warn(logger).Log("msg", "-api-refresh cannot be longer than 10m; setting it to 10m")
			*apiRefresh = 10 * time.Minute
		}
		if *apiTimeout < 5*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be shorter than 5s; setting it to 5s")
			*apiTimeout = 5 * time.Second
		}
		if *apiTimeout > 60*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be longer than 60s; setting it to 60s")
			*apiTimeout = 60 * time.Second
		}
		if *rtTimeout < 45*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be shorter than 45s; setting it to 45s")
			*rtTimeout = 45 * time.Second
		}
		if *rtTimeout > 120*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be longer than 120s; setting it to 120s")
			*rtTimeout = 120 * time.Second
		}
	}

	var serviceNameFilter filter.Filter
	{
		for _, expr := range serviceWhitelist {
			if err := serviceNameFilter.Whitelist(expr); err != nil {
				level.Error(logger).Log("err", "invalid -service-whitelist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "services", "type", "name whitelist", "expr", expr)
		}
		for _, expr := range serviceBlacklist {
			if err := serviceNameFilter.Blacklist(expr); err != nil {
				level.Error(logger).Log("err", "invalid -service-blacklist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "services", "type", "name blacklist", "expr", expr)
		}
	}

	var metricNameFilter filter.Filter
	{
		for _, expr := range metricWhitelist {
			if err := metricNameFilter.Whitelist(expr); err != nil {
				level.Error(logger).Log("err", "invalid -metric-whitelist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "metrics", "type", "name whitelist", "expr", expr)

		}
		for _, expr := range metricBlacklist {
			if err := metricNameFilter.Blacklist(expr); err != nil {
				level.Error(logger).Log("err", "invalid -metric-blacklist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "metrics", "type", "name blacklist", "expr", expr)
		}
	}

	var shardN, shardM uint64
	{
		if *serviceShard != "" {
			toks := strings.SplitN(*serviceShard, "/", 2)
			if len(toks) != 2 {
				level.Error(logger).Log("err", "-service-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			var err error
			shardN, err = strconv.ParseUint(toks[0], 10, 64)
			if err != nil {
				level.Error(logger).Log("err", "-service-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			if shardN <= 0 {
				level.Error(logger).Log("err", "first part of -service-shard flag should be greater than zero")
				os.Exit(1)
			}
			shardM, err = strconv.ParseUint(toks[1], 10, 64)
			if err != nil {
				level.Error(logger).Log("err", "-service-shard must be of the format 'n/m'")
				os.Exit(1)
			}
			if shardN > shardM {
				level.Error(logger).Log("err", fmt.Sprintf("-service-shard with n=%d m=%d is invalid: n must be less than or equal to m", shardN, shardM))
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "services", "type", "by shard", "n", shardN, "m", shardM)

		}
	}

	var registry *prometheus.Registry
	{
		registry = prometheus.NewRegistry()
	}

	var metrics *prom.Metrics
	{
		var err error
		metrics, err = prom.NewMetrics(*namespace, *subsystem, metricNameFilter, registry)
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
		apiCacheOptions = append(apiCacheOptions,
			api.WithLogger(apiLogger),
			api.WithNameFilter(serviceNameFilter),
		)

		if len(serviceIDs) > 0 {
			level.Info(logger).Log("filter", "services", "type", "explicit service IDs", "count", len(serviceIDs))
			apiCacheOptions = append(apiCacheOptions, api.WithExplicitServiceIDs(serviceIDs...))
		}

		if shardM > 0 {
			apiCacheOptions = append(apiCacheOptions, api.WithShard(shardN, shardM))
		}
	}

	var apiClient *http.Client
	{
		apiClient = &http.Client{
			Timeout: *apiTimeout,
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
			Timeout: *rtTimeout,
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
		// Serve Prometheus metrics (and /debug/pprof/...) over HTTP.
		http.Handle(promURL.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		server := http.Server{
			Addr:    promURL.Host,
			Handler: http.DefaultServeMux,
		}
		g.Add(func() error {
			return server.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			level.Debug(logger).Log("msg", "shutting down HTTP server")
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
