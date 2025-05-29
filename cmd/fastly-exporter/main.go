// Package main is the entry point for the fastly-exporter.
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

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/fastly/fastly-exporter/pkg/prom"
	"github.com/fastly/fastly-exporter/pkg/rt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var programVersion = "dev"

func main() {
	var (
		token               string
		listen              string
		namespace           string
		deprecatedSubsystem string
		serviceShard        string
		serviceIDs          stringslice
		serviceAllowlist    stringslice
		serviceBlocklist    stringslice
		metricAllowlist     stringslice
		metricBlocklist     stringslice
		certificateRefresh  time.Duration
		datacenterRefresh   time.Duration
		productRefresh      time.Duration
		serviceRefresh      time.Duration
		apiTimeout          time.Duration
		rtTimeout           time.Duration
		aggregateOnly       bool
		debug               bool
		versionFlag         bool
		configFileExample   bool
	)

	fs := flag.NewFlagSet("fastly-exporter", flag.ContinueOnError)
	{
		fs.StringVar(&token, "token", "", "Fastly API token (required)")
		fs.StringVar(&listen, "listen", "127.0.0.1:8080", "listen address for Prometheus metrics")
		fs.StringVar(&namespace, "namespace", "fastly", "Prometheus namespace")
		fs.StringVar(&deprecatedSubsystem, "subsystem", "rt", "DEPRECATED -- will be fixed to 'rt' in a future version")
		fs.StringVar(&serviceShard, "service-shard", "", "if set, only include services whose hashed IDs modulo m equal n-1 (format 'n/m')")
		fs.Var(&serviceIDs, "service", "if set, only include this service ID (repeatable)")
		fs.Var(&serviceAllowlist, "service-allowlist", "if set, only include services whose names match this regex (repeatable)")
		fs.Var(&serviceBlocklist, "service-blocklist", "if set, don't include services whose names match this regex (repeatable)")
		fs.Var(&metricAllowlist, "metric-allowlist", "if set, only export metrics whose names match this regex (repeatable)")
		fs.Var(&metricBlocklist, "metric-blocklist", "if set, don't export metrics whose names match this regex (repeatable)")
		fs.DurationVar(&certificateRefresh, "certificate-refresh", 6*time.Hour, "how often to poll api.fastly.com for updated custom TLS certificate metadata (10m–24h); a value of 0 will disable certificate refresh")
		fs.DurationVar(&datacenterRefresh, "datacenter-refresh", 10*time.Minute, "how often to poll api.fastly.com for updated datacenter metadata (10m–1h)")
		fs.DurationVar(&productRefresh, "product-refresh", 10*time.Minute, "how often to poll api.fastly.com for updated product metadata (10m–24h)")
		fs.DurationVar(&serviceRefresh, "service-refresh", 1*time.Minute, "how often to poll api.fastly.com for updated service metadata (15s–10m)")
		fs.DurationVar(&serviceRefresh, "api-refresh", 1*time.Minute, "DEPRECATED -- use service-refresh instead")
		fs.DurationVar(&apiTimeout, "api-timeout", 15*time.Second, "HTTP client timeout for api.fastly.com requests (5–60s)")
		fs.DurationVar(&rtTimeout, "rt-timeout", 45*time.Second, "HTTP client timeout for rt.fastly.com requests (45–120s)")
		fs.BoolVar(&aggregateOnly, "aggregate-only", false, "Use aggregated data rather than per-datacenter")
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
		logger = level.NewFilter(logger, getLogLevel(debug))
	}

	if token == "" {
		if token = os.Getenv("FASTLY_API_TOKEN"); token == "" {
			level.Error(logger).Log("err", "-token or FASTLY_API_TOKEN is required")
			os.Exit(1)
		}
	}

	switch deprecatedSubsystem {
	case "rt":
		// good
	case "origin":
		level.Error(logger).Log("err", "-subsystem cannot be 'origin'")
		os.Exit(1)
	default:
		level.Warn(logger).Log("msg", "-subsystem is DEPRECATED and will be fixed to 'rt' in a future version")
	}

	fs.Visit(func(f *flag.Flag) {
		if f.Name == "api-refresh" {
			level.Warn(logger).Log("msg", "-api-refresh is deprecated and will be removed in a future version, please use -service-refresh instead")
		}
	})

	{
		if certificateRefresh == 0 {
			level.Info(logger).Log("msg", "-certificate-refresh is disabled; set to a duration between 10m-24h to enable")
		} else {
			if certificateRefresh < 10*time.Minute {
				level.Warn(logger).Log("msg", "-certificate-refresh cannot be shorter than 10m; setting it to 10m")
				certificateRefresh = 10 * time.Minute
			}
			if certificateRefresh > 24*time.Hour {
				level.Warn(logger).Log("msg", "-certificate-refresh cannot be longer than 24h; setting it to 24h")
				certificateRefresh = 24 * time.Hour
			}
		}
		if datacenterRefresh < 10*time.Minute {
			level.Warn(logger).Log("msg", "-datacenter-refresh cannot be shorter than 10m; setting it to 10m")
			datacenterRefresh = 10 * time.Minute
		}
		if datacenterRefresh > 1*time.Hour {
			level.Warn(logger).Log("msg", "-datacenter-refresh cannot be longer than 1h; setting it to 1h")
			datacenterRefresh = 1 * time.Hour
		}
		if productRefresh < 10*time.Minute {
			level.Warn(logger).Log("msg", "-product-refresh cannot be shorter than 10m; setting it to 10m")
			productRefresh = 10 * time.Minute
		}
		if productRefresh > 24*time.Hour {
			level.Warn(logger).Log("msg", "-product-refresh cannot be longer than 24h; setting it to 24h")
			productRefresh = 24 * time.Hour
		}
		if serviceRefresh < 15*time.Second {
			level.Warn(logger).Log("msg", "-service-refresh cannot be shorter than 15s; setting it to 15s")
			serviceRefresh = 15 * time.Second
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

	var userAgent string
	{
		userAgent = `Fastly-Exporter (` + programVersion + `)`
	}

	var apiClient *http.Client
	{
		apiClient = &http.Client{
			Timeout:   apiTimeout,
			Transport: userAgentTransport(http.DefaultTransport, userAgent),
		}
	}

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

	var certificateCache *api.CertificateCache
	{
		enabled := certificateRefresh != 0 && !metricNameFilter.Blocked(prometheus.BuildFQName(namespace, deprecatedSubsystem, "cert_expiry_timestamp_seconds"))
		certificateCache = api.NewCertificateCache(apiClient, token, enabled)
	}
	var datacenterCache *api.DatacenterCache
	{
		enabled := !metricNameFilter.Blocked(prometheus.BuildFQName(namespace, deprecatedSubsystem, "datacenter_info"))
		datacenterCache = api.NewDatacenterCache(apiClient, token, enabled)
	}

	var productCache *api.ProductCache
	{
		productCache = api.NewProductCache(apiClient, token, apiLogger)
	}

	{
		var g errgroup.Group
		g.Go(func() error {
			if err := serviceCache.Refresh(context.Background()); err != nil {
				level.Warn(logger).Log("during", "initial fetch of service IDs", "err", err, "msg", "service metrics unavailable, will retry")
			}
			return nil
		})
		if certificateCache.Enabled() {
			g.Go(func() error {
				if err := certificateCache.Refresh(context.Background()); err != nil {
					level.Warn(logger).Log("during", "initial fetch of certificates", "err", err, "msg", "certificate labels unavailable, will retry")
				}
				return nil
			})
		}
		if datacenterCache.Enabled() {
			g.Go(func() error {
				if err := datacenterCache.Refresh(context.Background()); err != nil {
					level.Warn(logger).Log("during", "initial fetch of datacenters", "err", err, "msg", "datacenter labels unavailable, will retry")
				}
				return nil
			})
		}
		g.Go(func() error {
			if err := productCache.Refresh(context.Background()); err != nil {
				level.Warn(logger).Log("during", "initial fetch of products", "err", err, "msg", "products API unavailable, will retry")
			}
			return nil
		})

		g.Wait()
	}

	var defaultGatherers prometheus.Gatherers
	if certificateCache.Enabled() {
		certs, err := certificateCache.Gatherer(namespace, deprecatedSubsystem)
		if err != nil {
			level.Error(apiLogger).Log("during", "create certificate gatherer", "err", err)
			os.Exit(1)
		}
		defaultGatherers = append(defaultGatherers, certs)
	}

	if datacenterCache.Enabled() {
		dcs, err := datacenterCache.Gatherer(namespace, deprecatedSubsystem)
		if err != nil {
			level.Error(apiLogger).Log("during", "create datacenter gatherer", "err", err)
			os.Exit(1)
		}
		defaultGatherers = append(defaultGatherers, dcs)
	}

	if !metricNameFilter.Blocked(prometheus.BuildFQName(namespace, deprecatedSubsystem, "token_expiration")) {
		tokenRecorder := api.NewTokenRecorder(apiClient, token)
		tg, err := tokenRecorder.Gatherer(namespace, deprecatedSubsystem)
		if err != nil {
			level.Error(apiLogger).Log("during", "create token gatherer", "err", err)
		} else {
			err = tokenRecorder.Set(context.Background())
			if err != nil {
				level.Error(apiLogger).Log("during", "set token gauge metric", "err", err)
			}
			defaultGatherers = append(defaultGatherers, tg)
		}
	}

	var registry *prom.Registry
	{
		registry = prom.NewRegistry(programVersion, namespace, deprecatedSubsystem, metricNameFilter, defaultGatherers)
	}

	var manager *rt.Manager
	{
		var (
			rtLogger          = log.With(logger, "component", "rt.fastly.com")
			rtClient          = &http.Client{Timeout: rtTimeout, Transport: userAgentTransport(http.DefaultTransport, userAgent)}
			subscriberOptions = []rt.SubscriberOption{
				rt.WithLogger(rtLogger),
				rt.WithMetadataProvider(serviceCache),
				rt.WithAggregateOnly(aggregateOnly),
			}
		)
		manager = rt.NewManager(serviceCache, rtClient, token, registry, subscriberOptions, productCache, rtLogger)
		manager.Refresh() // populate initial subscribers, based on the initial cache refresh
	}

	var g run.Group
	// only setup the ticker if the certificateCache is enabled.
	if certificateCache.Enabled() {

		// Every certificateRefresh, ask the api.CertificateCache to refresh
		// metadata from the api.fastly.com/tls/certificates endpoint.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(certificateRefresh)
		)
		g.Add(func() error {
			for {
				select {
				case <-ticker.C:
					if err := certificateCache.Refresh(ctx); err != nil {
						level.Warn(apiLogger).Log("during", "certificate refresh", "err", err, "msg", "the certificate info metrics may be stale")
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
	// only setup the ticker if the datacenterCache is enabled.
	if datacenterCache.Enabled() {

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
		// Every productRefresh, ask the api.ProductCache to refresh
		// data from the product entitlement endpoint.
		var (
			ctx, cancel = context.WithCancel(context.Background())
			ticker      = time.NewTicker(productRefresh)
		)
		g.Add(func() error {
			for {
				select {
				case <-ticker.C:
					if err := productCache.Refresh(ctx); err != nil {
						level.Warn(apiLogger).Log("during", "product refresh", "err", err, "msg", "the product entitlement data may be stale")
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

func getLogLevel(debug bool) level.Option {
	switch {
	case debug:
		return level.AllowDebug()
	default:
		return level.AllowInfo()
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func userAgentTransport(next http.RoundTripper, userAgent string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.Header.Set("User-Agent", userAgent)
		return next.RoundTrip(req)
	})
}
