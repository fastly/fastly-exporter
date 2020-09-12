package rt_test

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertStringSliceEqual(t *testing.T, want, have []string) {
	t.Helper()

	sort.Strings(want)
	sort.Strings(have)

	if !cmp.Equal(want, have) {
		t.Error(cmp.Diff(want, have))
	}
}

//
//
//

type mockCache struct {
	mtx      sync.RWMutex
	services []api.Service
}

func (c *mockCache) update(services []api.Service) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.services = services
}

func (c *mockCache) ServiceIDs() (ids []string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, s := range c.services {
		ids = append(ids, s.ID)
	}
	return ids
}

func (c *mockCache) Metadata(id string) (name string, version int, found bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, s := range c.services {
		if s.ID == id {
			return s.Name, s.Version, true
		}
	}
	return name, version, false
}

//
//
//

type mockRealtimeClient struct {
	responses     []string
	next          chan struct{}
	lastUserAgent string
}

func newMockRealtimeClient(responses ...string) *mockRealtimeClient {
	next := make(chan struct{}, 1)
	next <- struct{}{} // first request works out-of-the-box
	return &mockRealtimeClient{
		responses: responses,
		next:      next,
	}
}

func (c *mockRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	select {
	case <-c.next:
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}

	var response string
	if len(c.responses) <= 1 {
		response = c.responses[0]
	} else {
		response, c.responses = c.responses[0], c.responses[1:]
	}

	c.lastUserAgent = req.Header.Get("User-Agent")

	return fixedResponseClient{200, response}.Do(req)
}

func (c *mockRealtimeClient) advance() {
	c.next <- struct{}{}
}

//
//
//

type countingRealtimeClient struct {
	code     int
	response string
	served   uint64
}

func (c *countingRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	atomic.AddUint64(&(c.served), 1)
	return fixedResponseClient{c.code, c.response}.Do(req)
}

//
//
//

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

//
//
//

func prometheusOutput(t *testing.T, registry *prometheus.Registry, prefix string) []string {
	t.Helper()

	server := httptest.NewServer(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	defer server.Close()
	resp, err := http.Get(server.URL)
	assertNoErr(t, err)

	var selected []string
	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), prefix) {
			selected = append(selected, s.Text())
		}
	}

	return selected
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
testspace_testsystem_bereq_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_bereq_header_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 599
testspace_testsystem_billed_body_bytes_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 39
testspace_testsystem_billed_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_blacklist_total{datacenter="BWI",service_id="my-service-id",service_name="my-service-name"} 0
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
