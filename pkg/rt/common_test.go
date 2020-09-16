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
                "SJC": {
                    "requests": 1,
                    "tls": 1,
                    "resp_header_bytes": 450,
                    "header_size": 450,
                    "resp_body_bytes": 1494,
                    "body_size": 1494,
                    "req_header_bytes": 210,
                    "bereq_header_bytes": 675,
                    "billed_header_bytes": 450,
                    "billed_body_bytes": 1494,
                    "status_2xx": 1,
                    "status_200": 1,
                    "hits": 0,
                    "miss": 1,
                    "pass": 0,
                    "synth": 0,
                    "errors": 0,
                    "hits_time": 0,
                    "miss_time": 0.069132,
                    "miss_histogram": {
                        "70": 1
                    },
                    "tls_v12": 1,
                    "object_size_10k": 1,
                    "recv_sub_time": 860249,
                    "recv_sub_count": 2,
                    "hash_sub_time": 4540,
                    "hash_sub_count": 1,
                    "miss_sub_time": 994,
                    "miss_sub_count": 1,
                    "deliver_sub_time": 19082,
                    "deliver_sub_count": 1,
                    "prehash_sub_time": 802,
                    "prehash_sub_count": 1,
                    "predeliver_sub_time": 835,
                    "predeliver_sub_count": 1,
                    "miss_resp_body_bytes": 1494,
                    "edge_requests": 1,
                    "edge_resp_header_bytes": 450,
                    "edge_resp_body_bytes": 1494,
                    "origin_fetches": 1,
                    "origin_fetch_header_bytes": 675,
                    "origin_fetch_resp_header_bytes": 250,
                    "origin_fetch_resp_body_bytes": 2907
                }
            },
            "aggregated": {
                "requests": 1,
                "tls": 1,
                "resp_header_bytes": 450,
                "header_size": 450,
                "resp_body_bytes": 1494,
                "body_size": 1494,
                "req_header_bytes": 210,
                "bereq_header_bytes": 675,
                "billed_header_bytes": 450,
                "billed_body_bytes": 1494,
                "status_2xx": 1,
                "status_200": 1,
                "hits": 0,
                "miss": 1,
                "pass": 0,
                "synth": 0,
                "errors": 0,
                "hits_time": 0,
                "miss_time": 0.069132,
                "miss_histogram": {
                    "70": 1
                },
                "tls_v12": 1,
                "object_size_10k": 1,
                "recv_sub_time": 860249,
                "recv_sub_count": 2,
                "hash_sub_time": 4540,
                "hash_sub_count": 1,
                "miss_sub_time": 994,
                "miss_sub_count": 1,
                "deliver_sub_time": 19082,
                "deliver_sub_count": 1,
                "prehash_sub_time": 802,
                "prehash_sub_count": 1,
                "predeliver_sub_time": 835,
                "predeliver_sub_count": 1,
                "miss_resp_body_bytes": 1494,
                "edge_requests": 1,
                "edge_resp_header_bytes": 450,
                "edge_resp_body_bytes": 1494,
                "origin_fetches": 1,
                "origin_fetch_header_bytes": 675,
                "origin_fetch_resp_header_bytes": 250,
                "origin_fetch_resp_body_bytes": 2907
            },
            "recorded": 1600223685
        }
    ],
    "Timestamp": 1600223694,
    "AggregateDelay": 9
}`

var expectedMetricsOutputSlice = strings.Split(strings.TrimSpace(`
testspace_testsystem_attack_blocked_req_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_blocked_req_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_logged_req_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_logged_req_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_passed_req_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_passed_req_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_req_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_req_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_attack_resp_synth_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_bereq_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_bereq_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 675
testspace_testsystem_billed_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1494
testspace_testsystem_billed_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 450
testspace_testsystem_billed_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_blacklist_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_body_size_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1494
testspace_testsystem_deliver_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_deliver_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 19082
testspace_testsystem_edge_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1494
testspace_testsystem_edge_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 450
testspace_testsystem_edge_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_error_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_error_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_errors_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_fetch_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_fetch_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hash_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_hash_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 4540
testspace_testsystem_header_size_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 450
testspace_testsystem_hit_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hit_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hit_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hits_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_hits_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_http2_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_shield_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transform_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transform_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgopto_transforms_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_frames_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_shield_frames_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_shield_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_shield_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_shield_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_imgvideo_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_ipv6_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_log_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_logging_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="+Inf"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.005"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.01"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.025"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.05"} 0
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.1"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.25"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.5"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="1"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="10"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="2.5"} 1
testspace_testsystem_miss_duration_seconds_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="5"} 1
testspace_testsystem_miss_duration_seconds_count{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_miss_duration_seconds_sum{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0.07
testspace_testsystem_miss_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1494
testspace_testsystem_miss_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_miss_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 994
testspace_testsystem_miss_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0.069132
testspace_testsystem_miss_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="+Inf"} 1
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.005"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.01"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.025"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.05"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.1"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.25"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="0.5"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="1"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="10"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="2.5"} 0
testspace_testsystem_object_size_bytes_bucket{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",le="5"} 0
testspace_testsystem_object_size_bytes_count{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_object_size_bytes_sum{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 10240
testspace_testsystem_origin_fetch_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_origin_fetch_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 675
testspace_testsystem_origin_fetch_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 2907
testspace_testsystem_origin_fetch_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 250
testspace_testsystem_origin_fetches_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_origin_revalidations_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_deliver_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_manifests_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_shield_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transform_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_otfp_transforms_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pass_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pci_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pipe_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_pipe_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_predeliver_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_predeliver_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 835
testspace_testsystem_prehash_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_prehash_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 802
testspace_testsystem_realtime_api_requests_total{result="success",service_id="my-service-id",service_name="my-service-name"} 1
testspace_testsystem_recv_sub_count_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 2
testspace_testsystem_recv_sub_time_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 860249
testspace_testsystem_req_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_req_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 210
testspace_testsystem_requests_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 2
testspace_testsystem_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 1494
testspace_testsystem_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 450
testspace_testsystem_restarts_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_segblock_origin_fetches_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_segblock_shield_fetches_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_service_info{service_id="my-service-id",service_name="my-service-name",service_version="123"} 1
testspace_testsystem_shield_fetch_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_fetch_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_fetch_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_fetch_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_fetches_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_resp_body_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_resp_header_bytes_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_revalidations_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_shield_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="200"} 1
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="204"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="206"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="301"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="302"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="304"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="400"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="401"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="403"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="404"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="416"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="429"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="500"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="501"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="502"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="503"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="504"} 0
testspace_testsystem_status_code_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_code="505"} 0
testspace_testsystem_status_group_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_group="1xx"} 0
testspace_testsystem_status_group_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_group="2xx"} 1
testspace_testsystem_status_group_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_group="3xx"} 0
testspace_testsystem_status_group_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_group="4xx"} 0
testspace_testsystem_status_group_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",status_group="5xx"} 0
testspace_testsystem_synth_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_tls_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",tls_version="any"} 1
testspace_testsystem_tls_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",tls_version="v10"} 0
testspace_testsystem_tls_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",tls_version="v11"} 0
testspace_testsystem_tls_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",tls_version="v12"} 1
testspace_testsystem_tls_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name",tls_version="v13"} 0
testspace_testsystem_uncacheable_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_video_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_blocked_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_logged_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
testspace_testsystem_waf_passed_total{datacenter="SJC",service_id="my-service-id",service_name="my-service-name"} 0
`), "\n")
