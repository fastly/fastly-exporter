package rt_test

import (
	"context"
	"testing"
	"time"

	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/gen"
	"github.com/peterbourgon/fastly-exporter/pkg/rt"
	"github.com/prometheus/client_golang/prometheus"
)

func TestSubscriberFixture(t *testing.T) {
	var (
		namespace  = "testspace"
		subsystem  = "testsystem"
		registry   = prometheus.NewRegistry()
		nameFilter = filter.Filter{}
		metrics    = gen.NewMetrics(namespace, subsystem, nameFilter, registry)
	)

	var (
		client         = newMockRealtimeClient(rtResponseFixture, `{}`)
		serviceID      = "my-service-id"
		serviceName    = "my-service-name"
		serviceVersion = 123
		cache          = &mockCache{}
		processed      = make(chan struct{})
		postprocess    = func() { close(processed) }
		config         = rt.SubscriberConfig{Client: client, ServiceID: serviceID, Metrics: metrics, Metadata: cache, Postprocess: postprocess}
	)

	cache.update([]api.Service{{ID: serviceID, Name: serviceName, Version: serviceVersion}})

	subscriber, err := rt.NewSubscriber(config)
	assertNoErr(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { subscriber.Run(ctx); close(done) }()

	<-processed

	output := prometheusOutput(t, registry, namespace+"_"+subsystem+"_")
	assertMetricOutput(t, expetedMetricsOutputMap, output)

	cancel()
	<-done
}

func TestSubscriberNoData(t *testing.T) {
	var (
		client      = newMockRealtimeClient(`{"Error": "No data available, please retry"}`, `{}`)
		registry    = prometheus.NewRegistry()
		metrics     = gen.NewMetrics("ns", "ss", filter.Filter{}, registry)
		processed   = make(chan struct{}, 100)
		postprocess = func() { processed <- struct{}{} }
		config      = rt.SubscriberConfig{Client: client, ServiceID: "service_id", Metrics: metrics, Postprocess: postprocess}
	)

	subscriber, err := rt.NewSubscriber(config)
	assertNoErr(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { subscriber.Run(ctx); close(done) }()
	defer func() { cancel(); <-done }()

	<-processed // No data
	client.advance()
	<-processed // OK
	client.advance()
	<-processed // OK

	want := map[string]float64{
		`ns_ss_realtime_api_requests_total{result="no data",service_id="service_id",service_name="service_id"}`: 1,
		`ns_ss_realtime_api_requests_total{result="success",service_id="service_id",service_name="service_id"}`: 2,
	}
	have := prometheusOutput(t, registry, "ns_ss_realtime_api_requests_total")
	assertMetricOutput(t, want, have)
}

func TestBadTokenNoSpam(t *testing.T) {
	var (
		respc    = make(chan struct{})
		client   = &forwardingClient{code: 403, response: `{"Error": "unauthorized"}`, done: respc}
		metrics  = gen.NewMetrics("namespace", "subsystem", filter.Filter{}, prometheus.NewRegistry())
		delayc   = make(chan time.Duration, 100)
		proceedc = make(chan struct{})
		delay    = func(c context.Context, d time.Duration) { delayc <- d; <-proceedc }
		config   = rt.SubscriberConfig{Client: client, Token: "presumably bad token", ServiceID: "service ID", Metrics: metrics, Delay: delay}
	)

	subscriber, err := rt.NewSubscriber(config)
	assertNoErr(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { subscriber.Run(ctx); close(done) }()

	select {
	case <-respc:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for first query")
	}

	var delayFor time.Duration
	select {
	case delayFor = <-delayc:
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for post-query delay")
	}
	if delayFor < time.Second {
		t.Fatalf("inter-query delay (%s) too short", delayFor)
	}

	close(proceedc)
	cancel()
	<-done
}
