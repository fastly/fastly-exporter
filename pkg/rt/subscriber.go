package rt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/pkg/errors"
)

// HTTPClient is a consumer contract for the subscriber.
// It models a concrete http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// MetadataProvider is a consumer contract for the subscriber.
// It models the service lookup method of an api.Cache.
type MetadataProvider interface {
	Metadata(id string) (name string, version int, found bool)
}

// Subscriber polls rt.fastly.com for a single service ID.
// It emits the received real-time stats data to Prometheus.
type Subscriber struct {
	client      HTTPClient
	userAgent   string
	token       string
	serviceID   string
	provider    MetadataProvider
	metrics     *prom.Metrics
	postprocess func()
	logger      log.Logger
}

// SubscriberOption provides some additional behavior to a subscriber.
type SubscriberOption func(*Subscriber)

// WithUserAgent sets the User-Agent supplied to rt.fastly.com.
// By default, the DefaultUserAgent is used.
func WithUserAgent(ua string) SubscriberOption {
	return func(s *Subscriber) { s.userAgent = ua }
}

// WithMetadataProvider sets the resolver used to look up service names and
// versions. By default, a no-op metadata resolver is used, which causes each
// service to have its name set to its service ID, and its version set to
// "unknown".
func WithMetadataProvider(p MetadataProvider) SubscriberOption {
	return func(s *Subscriber) { s.provider = p }
}

// WithLogger sets the logger used by the subscriber while running.
// By default, no log events are emitted.
func WithLogger(logger log.Logger) SubscriberOption {
	return func(s *Subscriber) { s.logger = logger }
}

// WithPostprocess sets the postprocess function for the subscriber, which is
// invoked after each successful call to the real-time stats API. By default, a
// no-op postprocess function is invoked. This option is only useful for tests.
func WithPostprocess(f func()) SubscriberOption {
	return func(s *Subscriber) { s.postprocess = f }
}

// DefaultUserAgent passed to rt.fastly.com.
// To change, use the WithUserAgent option.
const DefaultUserAgent = "Fastly-Exporter (unknown version)"

// NewSubscriber returns a ready-to-use subscriber.
// Run must be called to update the metrics.
func NewSubscriber(client HTTPClient, token, serviceID string, metrics *prom.Metrics, options ...SubscriberOption) *Subscriber {
	s := &Subscriber{
		client:      client,
		userAgent:   DefaultUserAgent,
		token:       token,
		serviceID:   serviceID,
		metrics:     metrics,
		provider:    nopMetadataProvider{},
		postprocess: func() {},
		logger:      log.NewNopLogger(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// Run polls rt.fastly.com in a hot loop, collecting real-time stats information
// and emitting it to the Prometheus metrics provided to the constructor. The
// method returns when the context is canceled, or a non-recoverable error
// occurs.
func (s *Subscriber) Run(ctx context.Context) error {
	var ts uint64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			name, result, delay, newts, fatal := s.query(ctx, ts)
			s.metrics.RealtimeAPIRequestsTotal.WithLabelValues(s.serviceID, name, string(result)).Inc()
			if fatal != nil {
				return fatal
			}
			if delay > 0 {
				contextSleep(ctx, delay)
			}
			ts = newts
		}
	}
}

// query rt.fastly.com for the service ID represented by the subscriber, and
// with the provided starting timestamp. The function may block for several
// seconds; cancel the context to provoke early termination. On success, the
// received real-time data is processed, and the Prometheus metrics related to
// the Fastly service are updated.
//
// Returns the current name of the service, the broad class of result of the API
// request, any delay that should pass before query is invoked again, the new
// timestamp that should be provided to the next call to query, and an error.
// Recoverable errors are logged internally and not returned, so any non-nil
// error returned by this method should be considered fatal to the subscriber.
func (s *Subscriber) query(ctx context.Context, ts uint64) (currentName string, result rtResult, delay time.Duration, newts uint64, fatal error) {
	name, ver, found := s.provider.Metadata(s.serviceID)
	version := strconv.Itoa(ver)
	if !found {
		name, version = s.serviceID, "unknown"
	}

	// rt.fastly.com blocks until it has data to return.
	// It's safe to call in a (single-threaded!) hot loop.
	u := fmt.Sprintf("https://rt.fastly.com/v1/channel/%s/ts/%d", url.QueryEscape(s.serviceID), ts)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return name, rtResultError, 0, ts, errors.Wrap(err, "error constructing real-time stats API request")
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Fastly-Key", s.token)
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		level.Error(s.logger).Log("during", "execute request", "err", err)
		return name, rtResultError, time.Second, ts, nil
	}

	var rt realtimeResponse
	if err := jsoniterAPI.NewDecoder(resp.Body).Decode(&rt); err != nil {
		resp.Body.Close()
		level.Error(s.logger).Log("during", "decode response", "err", err)
		return name, rtResultError, time.Second, ts, nil
	}
	resp.Body.Close()

	rterr := rt.Error
	if rterr == "" {
		rterr = "<none>"
	}

	switch resp.StatusCode {
	case http.StatusOK:
		level.Debug(s.logger).Log("status_code", resp.StatusCode, "response_ts", rt.Timestamp, "err", rterr)
		if strings.Contains(rterr, "No data available") {
			result = rtResultNoData
		} else {
			result = rtResultSuccess
		}
		process(rt, s.serviceID, name, version, s.metrics)
		s.postprocess()

	case http.StatusUnauthorized, http.StatusForbidden:
		result = rtResultError
		level.Error(s.logger).Log("status_code", resp.StatusCode, "response_ts", rt.Timestamp, "err", rterr, "msg", "token may be invalid")
		delay = 15 * time.Second

	default:
		result = rtResultUnknown
		level.Error(s.logger).Log("status_code", resp.StatusCode, "response_ts", rt.Timestamp, "err", rterr)
		delay = 5 * time.Second
	}

	return name, result, delay, rt.Timestamp, nil
}

//
//
//

type rtResult string

const (
	rtResultUnknown rtResult = "unknown"
	rtResultError   rtResult = "error"
	rtResultNoData  rtResult = "no data"
	rtResultSuccess rtResult = "success"
)

//
//
//

type nopMetadataProvider struct{}

func (nopMetadataProvider) Metadata(string) (string, int, bool) { return "", 0, false }

// realtimeResponse models the response from rt.fastly.com. It can get quite
// large; when there are lots of services being monitored, unmarshaling to this
// type is the CPU bottleneck of the program.
type realtimeResponse struct {
	Timestamp uint64 `json:"Timestamp"`
	Data      []struct {
		Datacenter map[string]datacenter `json:"datacenter"`
		Aggregated datacenter            `json:"aggregated"`
		Recorded   uint64                `json:"recorded"`
	} `json:"Data"`
	Error string `json:"error"`
}

func (resp *realtimeResponse) unmarshalStdlib(data []byte) error {
	return json.Unmarshal(data, resp)
}

var jsoniterAPI = jsoniter.ConfigFastest

func (resp *realtimeResponse) unmarshalJsoniter(data []byte) error {
	return jsoniterAPI.Unmarshal(data, resp)
}

type datacenter struct {
	Requests                        uint64            `json:"requests"`
	TLS                             uint64            `json:"tls"`
	Shield                          uint64            `json:"shield"`
	IPv6                            uint64            `json:"ipv6"`
	ImgOpto                         uint64            `json:"imgopto"`
	ImgOptoShield                   uint64            `json:"imgopto_shield"`
	ImgOptoTransform                uint64            `json:"imgopto_transforms"`
	OTFP                            uint64            `json:"otfp"`
	OTFPShield                      uint64            `json:"otfp_shield"`
	OTFPTransform                   uint64            `json:"otfp_transforms"`
	OTFPManifest                    uint64            `json:"otfp_manifests"`
	Video                           uint64            `json:"video"`
	PCI                             uint64            `json:"pci"`
	Logging                         uint64            `json:"logging"`
	HTTP2                           uint64            `json:"http2"`
	RespHeaderBytes                 uint64            `json:"resp_header_bytes"`
	HeaderSize                      uint64            `json:"header_size"`
	RespBodyBytes                   uint64            `json:"resp_body_bytes"`
	BodySize                        uint64            `json:"body_size"`
	ReqHeaderBytes                  uint64            `json:"req_header_bytes"`
	ReqBodyBytes                    uint64            `json:"req_body_bytes"`
	BackendReqHeaderBytes           uint64            `json:"bereq_header_bytes"`
	BackendReqBodyBytes             uint64            `json:"bereq_body_bytes"`
	BilledHeaderBytes               uint64            `json:"billed_header_bytes"`
	BilledBodyBytes                 uint64            `json:"billed_body_bytes"`
	WAFBlocked                      uint64            `json:"waf_blocked"`
	WAFLogged                       uint64            `json:"waf_logged"`
	WAFPassed                       uint64            `json:"waf_passed"`
	AttackReqHeaderBytes            uint64            `json:"attack_req_header_bytes"`
	AttackReqBodyBytes              uint64            `json:"attack_req_body_bytes"`
	AttackRespSynthBytes            uint64            `json:"attack_resp_synth_bytes"`
	AttackLoggedReqHeaderBytes      uint64            `json:"attack_logged_req_header_bytes"`
	AttackLoggedReqBodyBytes        uint64            `json:"attack_logged_req_body_bytes"`
	AttackBlockedReqHeaderBytes     uint64            `json:"attack_blocked_req_header_bytes"`
	AttackBlockedReqBodyBytes       uint64            `json:"attack_blocked_req_body_bytes"`
	AttackPassedReqHeaderBytes      uint64            `json:"attack_passed_req_header_bytes"`
	AttackPassedReqBodyBytes        uint64            `json:"attack_passed_req_body_bytes"`
	ShieldRespHeaderBytes           uint64            `json:"shield_resp_header_bytes"`
	ShieldRespBodyBytes             uint64            `json:"shield_resp_body_bytes"`
	OTFPRespHeaderBytes             uint64            `json:"otfp_resp_header_bytes"`
	OTFPRespBodyBytes               uint64            `json:"otfp_resp_body_bytes"`
	OTFPShieldRespHeaderBytes       uint64            `json:"otfp_shield_resp_header_bytes"`
	OTFPShieldRespBodyBytes         uint64            `json:"otfp_shield_resp_body_bytes"`
	OTFPTransformRespHeaderBytes    uint64            `json:"otfp_transform_resp_header_bytes"`
	OTFPTransformRespBodyBytes      uint64            `json:"otfp_transform_resp_body_bytes"`
	OTFPShieldTime                  uint64            `json:"otfp_shield_time"`
	OTFPTransformTime               uint64            `json:"otfp_transform_time"`
	OTFPDeliverTime                 uint64            `json:"otfp_deliver_time"`
	ImgOptoRespHeaderBytes          uint64            `json:"imgopto_resp_header_bytes"`
	ImgOptoRespBodyBytes            uint64            `json:"imgopto_resp_body_bytes"`
	ImgOptoShieldRespHeaderBytes    uint64            `json:"imgopto_shield_resp_header_bytes"`
	ImgOptoShieldRespBodyBytes      uint64            `json:"imgopto_shield_resp_body_bytes"`
	ImgOptoTransformRespHeaderBytes uint64            `json:"imgopto_transform_resp_header_bytes"`
	ImgOptoTransformRespBodyBytes   uint64            `json:"imgopto_transform_resp_body_bytes"`
	Status1xx                       uint64            `json:"status_1xx"`
	Status2xx                       uint64            `json:"status_2xx"`
	Status3xx                       uint64            `json:"status_3xx"`
	Status4xx                       uint64            `json:"status_4xx"`
	Status5xx                       uint64            `json:"status_5xx"`
	Status200                       uint64            `json:"status_200"`
	Status204                       uint64            `json:"status_204"`
	Status301                       uint64            `json:"status_301"`
	Status302                       uint64            `json:"status_302"`
	Status304                       uint64            `json:"status_304"`
	Status400                       uint64            `json:"status_400"`
	Status401                       uint64            `json:"status_401"`
	Status403                       uint64            `json:"status_403"`
	Status404                       uint64            `json:"status_404"`
	Status416                       uint64            `json:"status_416"`
	Status500                       uint64            `json:"status_500"`
	Status501                       uint64            `json:"status_501"`
	Status502                       uint64            `json:"status_502"`
	Status503                       uint64            `json:"status_503"`
	Status504                       uint64            `json:"status_504"`
	Status505                       uint64            `json:"status_505"`
	Hits                            uint64            `json:"hits"`
	Misses                          uint64            `json:"miss"`
	Passes                          uint64            `json:"pass"`
	Synths                          uint64            `json:"synth"`
	Errors                          uint64            `json:"errors"`
	Uncacheable                     uint64            `json:"uncacheable"`
	HitsTime                        float64           `json:"hits_time"`
	MissTime                        float64           `json:"miss_time"`
	PassTime                        float64           `json:"pass_time"`
	MissHistogram                   map[string]uint64 `json:"miss_histogram"`
	TLSv10                          uint64            `json:"tls_v10"`
	TLSv11                          uint64            `json:"tls_v11"`
	TLSv12                          uint64            `json:"tls_v12"`
	TLSv13                          uint64            `json:"tls_v13"`
	ObjectSize1k                    uint64            `json:"object_size_1k"`
	ObjectSize10k                   uint64            `json:"object_size_10k"`
	ObjectSize100k                  uint64            `json:"object_size_100k"`
	ObjectSize1m                    uint64            `json:"object_size_1m"`
	ObjectSize10m                   uint64            `json:"object_size_10m"`
	ObjectSize100m                  uint64            `json:"object_size_100m"`
	ObjectSize1g                    uint64            `json:"object_size_1g"`
	RecvSubTime                     uint64            `json:"recv_sub_time"`
	RecvSubCount                    uint64            `json:"recv_sub_count"`
	HashSubTime                     uint64            `json:"hash_sub_time"`
	HashSubCount                    uint64            `json:"hash_sub_count"`
	MissSubTime                     uint64            `json:"miss_sub_time"`
	MissSubCount                    uint64            `json:"miss_sub_count"`
	FetchSubTime                    uint64            `json:"fetch_sub_time"`
	FetchSubCount                   uint64            `json:"fetch_sub_count"`
	DeliverSubTime                  uint64            `json:"deliver_sub_time"`
	DeliverSubCount                 uint64            `json:"deliver_sub_count"`
	HitSubTime                      uint64            `json:"hit_sub_time"`
	HitSubCount                     uint64            `json:"hit_sub_count"`
	PrehashSubTime                  uint64            `json:"prehash_sub_time"`
	PrehashSubCount                 uint64            `json:"prehash_sub_count"`
	PredeliverSubTime               uint64            `json:"predeliver_sub_time"`
	PredeliverSubCount              uint64            `json:"predeliver_sub_count"`
}

func contextSleep(ctx context.Context, d time.Duration) {
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}
