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
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		token     = flag.String("token", "", "Fastly API token")
		service   = flag.String("service", "", "Fastly service")
		addr      = flag.String("endpoint", "http://127.0.0.1:8080/metrics", "Prometheus /metrics endpoint")
		namespace = flag.String("namespace", "", "Prometheus namespace")
		subsystem = flag.String("subsystem", "", "Prometheus subsystem")
		debug     = flag.Bool("debug", false, "log debug information")
	)
	flag.Parse()

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
