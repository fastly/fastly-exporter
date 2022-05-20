package rt_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/fastly/fastly-exporter/pkg/prom"
	"github.com/fastly/fastly-exporter/pkg/rt"
	"github.com/prometheus/client_golang/prometheus"
)

func TestSubscriberFixture(t *testing.T) {
	var (
		namespace  = "testspace"
		subsystem  = "testsystem"
		registry   = prometheus.NewRegistry()
		nameFilter = filter.Filter{}
		metrics    = prom.NewMetrics(namespace, subsystem, nameFilter, registry)
	)

	// Set up a subscriber.
	var (
		client         = newMockRealtimeClient(rtResponseFixture, `{}`)
		serviceID      = "my-service-id"
		serviceName    = "my-service-name"
		serviceVersion = 123
		cache          = &mockCache{}
		processed      = make(chan struct{})
		postprocess    = func() { close(processed) }
		options        = []rt.SubscriberOption{rt.WithMetadataProvider(cache), rt.WithPostprocess(postprocess)}
		subscriber     = rt.NewSubscriber(client, "irrelevant token", serviceID, metrics, options...)
	)

	// Prep the mock cache.
	cache.update([]api.Service{{ID: serviceID, Name: serviceName, Version: serviceVersion}})

	// Tell the subscriber to fetch real-time stats.
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() { errc <- subscriber.RunRealtime(ctx) }()

	// Block until the subscriber does finishes one fetch
	<-processed

	// Assert the Prometheus metrics.
	output := prometheusOutput(t, registry, namespace+"_"+subsystem+"_")
	assertMetricOutput(t, expetedMetricsOutputMap, output)

	// Kill the subscriber's goroutine, and wait for it to finish.
	cancel()
	err := <-errc
	switch {
	case err == nil:
	case errors.Is(err, context.Canceled):
	case err != nil:
		t.Fatal(err)
	}
}

func TestSubscriberNoData(t *testing.T) {
	var (
		client      = newMockRealtimeClient(`{"Error": "No data available, please retry"}`, `{}`)
		registry    = prometheus.NewRegistry()
		metrics     = prom.NewMetrics("ns", "ss", filter.Filter{}, registry)
		processed   = make(chan struct{}, 100)
		postprocess = func() { processed <- struct{}{} }
		options     = []rt.SubscriberOption{rt.WithPostprocess(postprocess)}
		subscriber  = rt.NewSubscriber(client, "token", "service_id", metrics, options...)
	)
	go subscriber.RunRealtime(context.Background())

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
		client     = &countingRealtimeClient{code: 403, response: `{"Error": "unauthorized"}`}
		metrics    = prom.NewMetrics("namespace", "subsystem", filter.Filter{}, prometheus.NewRegistry())
		subscriber = rt.NewSubscriber(client, "presumably bad token", "service ID", metrics)
	)
	go subscriber.RunRealtime(context.Background())

	time.Sleep(time.Second)

	if want, have := uint64(1), atomic.LoadUint64(&client.served); want != have {
		t.Fatalf("Unauthorized rt.fastly.com request count: want %d, have %d", want, have)
	}
}
