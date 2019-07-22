package rt_test

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/peterbourgon/fastly-exporter/pkg/rt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestSubscriberFixture(t *testing.T) {
	var (
		namespace  = "testspace"
		subsystem  = "testsystem"
		registry   = prometheus.NewRegistry()
		metrics, _ = prom.NewMetrics(namespace, subsystem, registry)
	)

	var (
		client         = &mockRealtimeClient{response: rtResponseFixture}
		serviceID      = "my-service-id"
		serviceName    = "my-service-name"
		serviceVersion = 123
		provider       = mockMetadataProvider{serviceID: api.Service{ID: serviceID, Name: serviceName, Version: serviceVersion}}
		processed      uint64                                       // we'll spin on this, to wait for the first process
		postprocess    = func() { atomic.AddUint64(&processed, 1) } // we'll update the processed var with this function
		options        = []rt.SubscriberOption{rt.WithMetadataProvider(provider), rt.WithPostprocess(postprocess)}
		subscriber     = rt.NewSubscriber(client, "irrelevant token", serviceID, metrics, options...)
	)

	var (
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
	)
	go func() {
		subscriber.Run(ctx)
		close(done)
	}()

	for atomic.LoadUint64(&processed) == 0 {
		time.Sleep(time.Millisecond)
	}

	server := httptest.NewServer(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	defer server.Close()
	resp, err := http.Get(server.URL)
	assertNoErr(t, err)
	var (
		s        = bufio.NewScanner(resp.Body)
		prefix   = namespace + "_" + subsystem + "_"
		selected []string
	)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), prefix) {
			selected = append(selected, s.Text())
		}
	}

	sort.Strings(selected)
	sort.Strings(expectedMetricsOutputSlice)
	if want, have := expectedMetricsOutputSlice, selected; !cmp.Equal(want, have) {
		t.Error(cmp.Diff(want, have))
	}

	cancel()
	<-done
}

func TestBadTokenNoSpam(t *testing.T) {
	var (
		client     = &countingRealtimeClient{code: 403, response: `{"Error": "unauthorized"}`}
		metrics, _ = prom.NewMetrics("namespace", "subsystem", prometheus.NewRegistry())
		subscriber = rt.NewSubscriber(client, "presumably bad token", "service ID", metrics)
	)
	go subscriber.Run(context.Background())

	time.Sleep(time.Second)

	if want, have := uint64(1), atomic.LoadUint64(&client.served); want != have {
		t.Fatalf("mock rt.fastly.com request count: want %d, have %d", want, have)
	}
}

func TestUserAgent(t *testing.T) {
	var (
		client     = &userAgentCapturingClient{}
		userAgent  = "Some user agent string"
		metrics, _ = prom.NewMetrics("ns", "ss", prometheus.NewRegistry())
		subscriber = rt.NewSubscriber(client, "token", "serviceid", metrics, rt.WithUserAgent(userAgent))
	)
	go subscriber.Run(context.Background())

	want, have := userAgent, ""
	if !within(time.Second, func() bool {
		have, _ = client.userAgent.Load().(string)
		return want == have
	}) {
		t.Fatalf("timeout waiting for correct User-Agent: want %q, have %q", want, have)
	}
}

//
//
//

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

type mockMetadataProvider map[string]api.Service

func (p mockMetadataProvider) Metadata(id string) (name string, version int, found bool) {
	if s, ok := p[id]; ok {
		name, version, found = s.Name, s.Version, true
	}
	return name, version, found
}

type mockRealtimeClient struct {
	response string
	served   uint64
}

func (c *mockRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	// First request immediately returns real data.
	if atomic.AddUint64(&(c.served), 1) == 1 {
		return fixedResponseClient{200, c.response}.Do(req)
	}

	// Subsequent requests block a bit and then return empty JSON.
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case <-time.After(time.Second):
		return fixedResponseClient{200, "{}"}.Do(req)
	}
}

type countingRealtimeClient struct {
	code     int
	response string
	served   uint64
}

func (c *countingRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	atomic.AddUint64(&(c.served), 1)
	return fixedResponseClient{c.code, c.response}.Do(req)
}

type userAgentCapturingClient struct {
	userAgent atomic.Value
}

func (c *userAgentCapturingClient) Do(req *http.Request) (*http.Response, error) {
	c.userAgent.Store(req.Header.Get("User-Agent"))
	return fixedResponseClient{200, "{}"}.Do(req)
}

type fixedResponseClient struct {
	code     int
	response string
}

func (c fixedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c.code)
		fmt.Fprint(w, c.response)
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

func within(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() { // 🔥
			return true
		}
	}
	return false
}

//
//
//

const rtResponseFixture = `{
	"Data": [
		{
			"datacenter": {
				"BWI": {
					"requests": 1,
					"resp_header_bytes": 441,
					"header_size": 441,
					"resp_body_bytes": 39,
					"body_size": 39,
					"req_header_bytes": 294,
					"bereq_header_bytes": 599,
					"billed_header_bytes": 441,
					"billed_body_bytes": 39,
					"status_4xx": 1,
					"status_404": 1,
					"hits": 0,
					"miss": 1,
					"pass": 0,
					"synth": 0,
					"errors": 0,
					"hits_time": 0,
					"miss_time": 0.008408,
					"miss_histogram": {
						"9": 1
					},
					"object_size_1k": 1,
					"recv_sub_time": 1872661,
					"recv_sub_count": 6,
					"hash_sub_time": 1716,
					"hash_sub_count": 1,
					"miss_sub_time": 60677,
					"miss_sub_count": 9,
					"fetch_sub_time": 70428,
					"fetch_sub_count": 4,
					"deliver_sub_time": 26291,
					"deliver_sub_count": 1,
					"prehash_sub_time": 689,
					"prehash_sub_count": 1,
					"predeliver_sub_time": 1251,
					"predeliver_sub_count": 1
				}
			},
			"aggregated": {
				"requests": 1,
				"resp_header_bytes": 441,
				"header_size": 441,
				"resp_body_bytes": 39,
				"body_size": 39,
				"req_header_bytes": 294,
				"bereq_header_bytes": 599,
				"billed_header_bytes": 441,
				"billed_body_bytes": 39,
				"status_4xx": 1,
				"status_404": 1,
				"hits": 0,
				"miss": 1,
				"pass": 0,
				"synth": 0,
				"errors": 0,
				"hits_time": 0,
				"miss_time": 0.008408,
				"miss_histogram": {
					"9": 1
				},
				"object_size_1k": 1,
				"recv_sub_time": 1872661,
				"recv_sub_count": 6,
				"hash_sub_time": 1716,
				"hash_sub_count": 1,
				"miss_sub_time": 60677,
				"miss_sub_count": 9,
				"fetch_sub_time": 70428,
				"fetch_sub_count": 4,
				"deliver_sub_time": 26291,
				"deliver_sub_count": 1,
				"prehash_sub_time": 689,
				"prehash_sub_count": 1,
				"predeliver_sub_time": 1251,
				"predeliver_sub_count": 1
			},
			"recorded": 1541085726
		}
	],
	"Timestamp": 1541085735,
	"AggregateDelay": 9
}`

var expectedMetricsOutputSlice = strings.Split(strings.TrimSpace(`
testspace_testsystem_attack_blocked_req_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_blocked_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_logged_req_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_logged_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_passed_req_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_passed_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_req_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_resp_synth_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_bereq_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 599
testspace_testsystem_billed_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 39
testspace_testsystem_billed_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 441
testspace_testsystem_body_size_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 39
testspace_testsystem_deliver_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_deliver_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 26291
testspace_testsystem_errors_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_fetch_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 4
testspace_testsystem_fetch_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 70428
testspace_testsystem_hash_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_hash_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1716
testspace_testsystem_header_size_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 441
testspace_testsystem_hit_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hit_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hits_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hits_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_http2_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transform_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transform_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transforms_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_ipv6_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_logging_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="+Inf"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.005"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.01"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.025"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.05"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.1"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.25"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="0.5"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="16"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="2"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="32"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="4"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="60"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="8"} 1
testspace_testsystem_miss_duration_seconds_count{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_miss_duration_seconds_sum{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0.009
testspace_testsystem_miss_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 9
testspace_testsystem_miss_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 60677
testspace_testsystem_miss_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0.008408
testspace_testsystem_miss_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="+Inf"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1.024e+06"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1.024e+07"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1.024e+08"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1.024e+09"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="1024"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="10240"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",le="102400"} 1
testspace_testsystem_object_size_bytes_count{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_object_size_bytes_sum{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1024
testspace_testsystem_otfp_deliver_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_manifests_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transforms_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pci_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_predeliver_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_predeliver_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1251
testspace_testsystem_prehash_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_prehash_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 689
testspace_testsystem_realtime_api_requests_total{result="success",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_recv_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 6
testspace_testsystem_recv_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1.872661e+06
testspace_testsystem_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 294
testspace_testsystem_requests_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 39
testspace_testsystem_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 441
testspace_testsystem_service_info{service_id="my-service-id",service_name="my-service-name",service_version="123"} 1
testspace_testsystem_shield_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="200"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="204"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="301"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="302"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="304"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="400"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="401"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="403"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="404"} 1
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="416"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="500"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="501"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="502"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="503"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="504"} 0
testspace_testsystem_status_code_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_code="505"} 0
testspace_testsystem_status_group_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_group="1xx"} 0
testspace_testsystem_status_group_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_group="2xx"} 0
testspace_testsystem_status_group_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_group="3xx"} 0
testspace_testsystem_status_group_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_group="4xx"} 1
testspace_testsystem_status_group_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",status_group="5xx"} 0
testspace_testsystem_synth_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_tls_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",tls_version="any"} 0
testspace_testsystem_tls_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",tls_version="v10"} 0
testspace_testsystem_tls_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",tls_version="v11"} 0
testspace_testsystem_tls_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",tls_version="v12"} 0
testspace_testsystem_tls_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name",tls_version="v13"} 0
testspace_testsystem_uncacheable_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_video_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_blocked_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_logged_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_passed_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
`), "\n")
