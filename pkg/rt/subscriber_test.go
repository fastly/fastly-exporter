package rt_test

import (
	"context"
	"sync/atomic"
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
		metrics, _ = gen.NewMetrics(namespace, subsystem, nameFilter, registry)
	)

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
	cache.update([]api.Service{{ID: serviceID, Name: serviceName, Version: serviceVersion}})

	var (
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
	)
	go func() {
		subscriber.Run(ctx)
		close(done)
	}()

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
		metrics, _  = gen.NewMetrics("ns", "ss", filter.Filter{}, registry)
		processed   = make(chan struct{}, 100)
		postprocess = func() { processed <- struct{}{} }
		options     = []rt.SubscriberOption{rt.WithPostprocess(postprocess)}
		subscriber  = rt.NewSubscriber(client, "token", "service_id", metrics, options...)
	)
	go subscriber.Run(context.Background())

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

func TestUserAgent(t *testing.T) {
	var (
		client      = newMockRealtimeClient(`{}`)
		userAgent   = "Some user agent string"
		metrics, _  = gen.NewMetrics("ns", "ss", filter.Filter{}, prometheus.NewRegistry())
		processed   = make(chan struct{})
		postprocess = func() { close(processed) }
		options     = []rt.SubscriberOption{rt.WithUserAgent(userAgent), rt.WithPostprocess(postprocess)}
		subscriber  = rt.NewSubscriber(client, "token", "service_id", metrics, options...)
	)
	go subscriber.Run(context.Background())

	<-processed

	if want, have := userAgent, client.lastUserAgent; want != have {
		t.Errorf("User-Agent: want %q, have %q", want, have)
	}
}

func TestBadTokenNoSpam(t *testing.T) {
	var (
		client     = &countingRealtimeClient{code: 403, response: `{"Error": "unauthorized"}`}
		metrics, _ = gen.NewMetrics("namespace", "subsystem", filter.Filter{}, prometheus.NewRegistry())
		subscriber = rt.NewSubscriber(client, "presumably bad token", "service ID", metrics)
	)
	go subscriber.Run(context.Background())

	time.Sleep(time.Second)

	if want, have := uint64(1), atomic.LoadUint64(&client.served); want != have {
		t.Fatalf("Unauthorized rt.fastly.com request count: want %d, have %d", want, have)
	}
}
