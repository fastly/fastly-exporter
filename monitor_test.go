package main

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMonitorFixture(t *testing.T) {
	// Bunch of setup.
	var (
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
		namespace   = "testspace"
		subsystem   = "testsystem"
		client      = &mockRealtimeClient{response: rtResponseFixture}
		token       = "irrelevant-token"
		serviceID   = "my-service-id"
		serviceName = "my-service-name"
		cache       = newNameCache()
		metrics     = prometheusMetrics{}
		logger      = log.NewNopLogger()
	)

	// Make sure the monitor goroutine terminates.
	defer func() {
		cancel()
		<-done
	}()

	// Set up the service name mapping, and register the Prometheus metrics.
	cache.update(map[string]string{serviceID: serviceName})
	metrics.register(namespace, subsystem)

	// We're going to wait until the first call to process is done.
	var processed uint64
	postprocess := func() { atomic.AddUint64(&processed, 1) }

	// Launch the monitor goroutine against our fixture data.
	go func() {
		monitor(ctx, client, token, serviceID, cache, metrics, postprocess, logger)
		close(done)
	}()

	// We're gonna read the Prometheus metrics output, a true end-to-end test.
	server := httptest.NewServer(promhttp.Handler())
	defer server.Close()

	// Spin until we've processed at least one mock response.
	for atomic.LoadUint64(&processed) == 0 {
		time.Sleep(time.Millisecond)
	}

	// Collect output from Prometheus.
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
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

	// Deterministic order.
	sort.Strings(selected)
	sort.Strings(expectedMetricsOutputSlice)

	// Compare output to expected output.
	if want, have := expectedMetricsOutputSlice, selected; !cmp.Equal(want, have) {
		t.Error(cmp.Diff(want, have))
	}
}

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
testspace_testsystem_recv_sub_count_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 6
testspace_testsystem_recv_sub_time_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1.872661e+06
testspace_testsystem_req_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 294
testspace_testsystem_requests_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_resp_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 39
testspace_testsystem_resp_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 441
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
