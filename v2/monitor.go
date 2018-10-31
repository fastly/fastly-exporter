package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func monitor(ctx context.Context, token string, serviceID string, resolver nameResolver, metrics prometheusMetrics, logger log.Logger) error {
	var ts uint64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			// rt.fastly.com blocks until it has data to return.
			// It's safe to call in a (single-threaded!) hot loop.
			u := fmt.Sprintf("https://rt.fastly.com/v1/channel/%s/ts/%d", serviceID, ts)
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				return err // fatal for sure
			}
			req.Header.Set("Fastly-Key", token)
			req.Header.Set("Accept", "application/json")
			resp, err := http.DefaultClient.Do(req.WithContext(ctx))
			if err != nil {
				level.Error(logger).Log("err", err)
				contextSleep(ctx, time.Second)
				continue
			}
			var rt realtimeResponse
			if err := json.NewDecoder(resp.Body).Decode(&rt); err != nil {
				level.Error(logger).Log("err", err)
				contextSleep(ctx, time.Second)
				continue
			}
			rterr := rt.Error
			if rterr == "" {
				rterr = "<none>"
			}
			level.Debug(logger).Log("response_ts", rt.Timestamp, "err", rterr)
			process(rt, serviceID, resolver.resolve(serviceID), metrics)
			ts = rt.Timestamp
		}
	}
}

type realtimeResponse struct {
	Timestamp uint64 `json:"Timestamp"`
	Data      []struct {
		Datacenter map[string]datacenter `json:"datacenter"`
		Aggregated datacenter            `json:"aggregated"`
		Recorded   uint64                `json:"recorded"`
	} `json:"Data"`
	Error string `json:"error"`
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
	BackendReqHeaderBytes           uint64            `json:"bereq_header_bytes"`
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
