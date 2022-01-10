package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/peterbourgon/fastly-exporter/pkg/rt"
	"github.com/peterbourgon/ff/v3"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var programVersion = "dev"

func main() {
	var (
		token             string
		listen            string
		namespace         string
		subsystem         string
		serviceShard      string
		userAgent         string
		serviceIDs        stringslice
		serviceAllowlist  stringslice
		serviceBlocklist  stringslice
		metricAllowlist   stringslice
		metricBlocklist   stringslice
		datacenterRefresh time.Duration
		serviceRefresh    time.Duration
		apiTimeout        time.Duration
		rtTimeout         time.Duration
		debug             bool
		versionFlag       bool
		configFileExample bool
	)

	fs := flag.NewFlagSet("fastly-exporter", flag.ContinueOnError)
	{
		fs.StringVar(&token, "token", "", "Fastly API token (required)")
		fs.StringVar(&listen, "listen", "127.0.0.1:8080", "listen address for Prometheus metrics")
		fs.StringVar(&namespace, "namespace", "fastly", "Prometheus namespace")
		fs.StringVar(&subsystem, "subsystem", "rt", "Prometheus subsystem")
		fs.StringVar(&serviceShard, "service-shard", "", "if set, only include services whose hashed IDs modulo m equal n-1 (format 'n/m')")
		fs.StringVar(&userAgent, "api-user-agent", `Fastly-Exporter (`+programVersion+`)`, "if set, override the reporting user-agent HTTP header")
		fs.Var(&serviceIDs, "service", "if set, only include this service ID (repeatable)")
		fs.Var(&serviceAllowlist, "service-allowlist", "if set, only include services whose names match this regex (repeatable)")
		fs.Var(&serviceBlocklist, "service-blocklist", "if set, don't include services whose names match this regex (repeatable)")
		fs.Var(&metricAllowlist, "metric-allowlist", "if set, only export metrics whose names match this regex (repeatable)")
		fs.Var(&metricBlocklist, "metric-blocklist", "if set, don't export metrics whose names match this regex (repeatable)")
		fs.DurationVar(&datacenterRefresh, "datacenter-refresh", 10*time.Minute, "how often to poll api.fastly.com for updated datacenter metadata (10m–1h)")
		fs.DurationVar(&serviceRefresh, "service-refresh", 1*time.Minute, "how often to poll api.fastly.com for updated service metadata (15s–10m)")
		fs.DurationVar(&serviceRefresh, "api-refresh", 1*time.Minute, "DEPRECATED -- use service-refresh instead")
		fs.DurationVar(&apiTimeout, "api-timeout", 15*time.Second, "HTTP client timeout for api.fastly.com requests (5–60s)")
		fs.DurationVar(&rtTimeout, "rt-timeout", 45*time.Second, "HTTP client timeout for rt.fastly.com requests (45–120s)")
		fs.BoolVar(&debug, "debug", false, "log debug information")
		fs.BoolVar(&versionFlag, "version", false, "print version information and exit")
		fs.String("config-file", "", "config file (optional)")
		fs.BoolVar(&configFileExample, "config-file-example", false, "print example config file to stdout and exit")
		fs.Usage = usageFor(fs)
	}
	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("FASTLY_EXPORTER"),
		ff.WithConfigFileFlag("config-file"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if versionFlag {
		fmt.Fprintf(os.Stdout, "fastly-exporter v%s\n", programVersion)
		os.Exit(0)
	}

	if configFileExample {
		fmt.Fprintln(os.Stdout, exampleConfigFile)
		os.Exit(0)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		loglevel := level.AllowInfo()
		if debug {
			loglevel = level.AllowDebug()
		}
		logger = level.NewFilter(logger, loglevel)
	}

	if token == "" {
		if token = os.Getenv("FASTLY_API_TOKEN"); token == "" {
			level.Error(logger).Log("err", "-token or FASTLY_API_TOKEN is required")
			os.Exit(1)
		}
	}

	fs.Visit(func(f *flag.Flag) {
		if f.Name == "api-refresh" {
			level.Warn(logger).Log("msg", "-api-refresh is deprecated and will be removed in a future version, please use -service-refresh instead")
		}
	})

	{
		if datacenterRefresh < 10*time.Minute {
			level.Warn(logger).Log("msg", "-datacenter-refresh cannot be shorter than 10m; setting it to 10m")
			datacenterRefresh = 10 * time.Minute
		}
		if datacenterRefresh > 1*time.Hour {
			level.Warn(logger).Log("msg", "-datacenter-refresh cannot be longer than 1h; setting it to 1h")
			datacenterRefresh = 1 * time.Hour
		}
		if serviceRefresh < 15*time.Second {
			level.Warn(logger).Log("msg", "-service-refresh cannot be shorter than 15s; setting it to 15s")
			serviceRefresh = 15 * time.Second
		}
		if serviceRefresh > 10*time.Minute {
			level.Warn(logger).Log("msg", "-service-refresh cannot be longer than 10m; setting it to 10m")
			serviceRefresh = 10 * time.Minute
		}
		if apiTimeout < 5*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be shorter than 5s; setting it to 5s")
			apiTimeout = 5 * time.Second
		}
		if apiTimeout > 60*time.Second {
			level.Warn(logger).Log("msg", "-api-timeout cannot be longer than 60s; setting it to 60s")
			apiTimeout = 60 * time.Second
		}
		if rtTimeout < 45*time.Second {
			level.Warn(logger).Log("msg", "-rt-timeout cannot be shorter than 45s; setting it to 45s")
			rtTimeout = 45 * time.Second
		}
		if rtTimeout > 120*time.Second {
			level.Warn(logger).Log("msg", "-rt-timeout cannot be longer than 120s; setting it to 120s")
			rtTimeout = 120 * time.Second
		}
	}

	var serviceNameFilter filter.Filter
	{
		for _, expr := range serviceAllowlist {
			if err := serviceNameFilter.Allow(expr); err != nil {
				level.Error(logger).Log("err", "invalid -service-allowlist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "services", "type", "name allowlist", "expr", expr)
		}
		for _, expr := range serviceBlocklist {
			if err := serviceNameFilter.Block(expr); err != nil {
				level.Error(logger).Log("err", "invalid -service-blocklist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "services", "type", "name blocklist", "expr", expr)
		}
	}

	var metricNameFilter filter.Filter
	{
		for _, expr := range metricAllowlist {
			if err := metricNameFilter.Allow(expr); err != nil {
				level.Error(logger).Log("err", "invalid -metric-allowlist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "metrics", "type", "name allowlist", "expr", expr)

		}
		for _, expr := range metricBlocklist {
			if err := metricNameFilter.Block(expr); err != nil {
				level.Error(logger).Log("err", "invalid -metric-blocklist", "msg", err)
				os.Exit(1)
			}
			level.Info(logger).Log("filter", "metrics", "type", "name blocklist", "expr", expr)
		}
	}

	var shardN, shardM uint64
	{
		if serviceShard != "" {
			toks := strings.SplitN(serviceShard, "/", 2)
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

	var apiLogger log.Logger
	{
		apiLogger = log.With(logger, "component", "api.fastly.com")
	}

	var apiClient *http.Client
	{
		apiClient = &http.Client{
			Timeout: apiTimeout,
		}
	}

	apiClient.Transport = api.NewCustomUserAgent(apiClient.Transport, userAgent)

	var serviceCache *api.ServiceCache
	{
		serviceCacheOptions := []api.ServiceCacheOption{
			api.WithLogger(apiLogger),
			api.WithNameFilter(serviceNameFilter),
		}

		if len(serviceIDs) > 0 {
			level.Info(logger).Log("filter", "services", "type", "explicit service IDs", "count", len(serviceIDs))
			serviceCacheOptions = append(serviceCacheOptions, api.WithExplicitServiceIDs(serviceIDs...))
		}

		if shardM > 0 {
			level.Info(logger).Log("filter", "services", "type", "shard", "shard", fmt.Sprintf("%d/%d", shardN, shardM))
			serviceCacheOptions = append(serviceCacheOptions, api.WithShard(shardN, shardM))
		}

		serviceCache = api.NewServiceCache(apiClient, token, serviceCacheOptions...)
	}

	var datacenterCache *api.DatacenterCache
	{
		datacenterCache = api.NewDatacenterCache(apiClient, token)
	}

	{
		var g errgroup.Group
		g.Go(func() error {
			if err := serviceCache.Refresh(context.Background()); err != nil {
				level.Warn(logger).Log("during", "initial fetch of service IDs", "err", err, "msg", "service metrics unavailable, will retry")
			}
			return nil
		})
		g.Go(func() error {
			if err := datacenterCache.Refresh(context.Background()); err != nil {
				level.Warn(logger).Log("during", "initial fetch of datacenters", "err", err, "msg", "datacenter labels unavailable, will retry")
			}
			return nil
		})
		g.Wait()
	}

	var defaultGatherers prometheus.Gatherers
	{
		dcs, err := datacenterCache.Gatherer(namespace, subsystem)
		if err != nil {
			level.Error(apiLogger).Log("during", "create datacenter gatherer", "err", err)
			os.Exit(1)
		}
		defaultGatherers = append(defaultGatherers, dcs)
	}

	var registry *prom.Registry
	{
		registry = prom.NewRegistry(programVersion, namespace, subsystem, metricNameFilter, defaultGatherers)
	}

	var manager *rt.Manager
	{
		var (
			rtLogger          = log.With(logger, "component", "rt.fastly.com")
			rtClient          = &http.Client{Timeout: rtTimeout}
			subscriberOptions = []rt.SubscriberOption{
				rt.WithLogger(rtLogger),
				rt.WithMetadataProvider(serviceCache),
			}
		)
		rtClient.Transport = api.NewCustomUserAgent(rtClient.Transport, userAgent)
		manager = rt.NewManager(serviceCache, rtClient, token, registry, subscriberOptions, rtLogger)
		manager.Refresh() // populate initial subscribers, based on the initial cache refresh
	}

	var g run.Group
	{
		// Every datacenterRefresh, ask the api.DatacenterCache to refresh
		// metadata from the api.fastly.com/datacenters endpoint.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(datacenterRefresh)
		)
		g.Add(func() error {
			for {
				select {
				case <-ticker.C:
					if err := datacenterCache.Refresh(ctx); err != nil {
						level.Warn(apiLogger).Log("during", "datacenter refresh", "err", err, "msg", "the datacenter info metrics may be stale")
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}, func(error) {
			ticker.Stop()
			cancel()
		})
	}
	{
		// Every serviceRefresh, ask the api.ServiceCache to refresh the set of
		// services we should be exporting data for. Then, ask the rt.Manager to
		// refresh its set of rt.Subscribers, based on those latest services.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(serviceRefresh)
		)
		g.Add(func() error {
			for {
				select {
				case <-ticker.C:
					if err := serviceCache.Refresh(ctx); err != nil {
						level.Warn(apiLogger).Log("during", "service refresh", "err", err, "msg", "the set of exported services and their metadata may be stale")
					}
					manager.Refresh() // safe to do with stale data in the cache
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}, func(error) {
			ticker.Stop()
			cancel()
		})
	}
	{
		// A pseudo-actor for the rt.Manager, which waits for interrupt and then
		// tears down all of the managed subscribers.
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
		// The HTTP server that Prometheus will scrape.
		serverLogger := log.With(logger, "component", "server")
		server := http.Server{
			Addr:    listen,
			Handler: registry,
		}
		g.Add(func() error {
			level.Info(serverLogger).Log("listen", listen)
			return server.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			level.Debug(serverLogger).Log("msg", "shutting down")
			server.Shutdown(ctx)
		})
	}
	{
		// Catch ctrl-C.
		var (
			ctx     = context.Background()
			signals = os.Interrupt
		)
		g.Add(run.SignalHandler(ctx, signals))
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

func usageFor(fs *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  fastly-exporter [flags]\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")

		tw := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			def := f.DefValue
			if def == "" {
				def = "..."
			}
			fmt.Fprintf(tw, "  -%s %s\t%s%s\n", f.Name, def, f.Usage, envVarSuffix(f))
		})
		tw.Flush()

		fmt.Fprintf(os.Stderr, "\n")
	}
}

func envVarSuffix(f *flag.Flag) string {
	if _, ok := f.Value.(*stringslice); ok {
		return "" // no repeatable flags as env vars
	}

	switch f.Name {
	case "version", "config-file-example":
		return ""

	case "token":
		return " (or via FASTLY_API_TOKEN)"

	default:
		return " (or via FASTLY_EXPORTER_" + strings.Replace(strings.ToUpper(f.Name), "-", "_", -1) + ")"
	}
}

var exampleConfigFile = strings.TrimSpace(`
token ABC123

service-refresh 30s
api-timeout 60s

service-allowlist Prod
service-blocklist Staging
service-blocklist Dev

metric-blocklist imgopto
`)
