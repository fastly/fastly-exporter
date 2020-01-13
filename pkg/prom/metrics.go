package prom

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics are the concrete Prometheus metrics that get updated with data
// retrieved from the real-time stats API. The same set of metrics are updated
// for all service IDs, only the labels change.
type Metrics struct {
	// These metrics concern the exporter itself.
	RealtimeAPIRequestsTotal *prometheus.CounterVec

	// These metrics concern the Fastly service.
	ServiceInfo                          *prometheus.GaugeVec
	RequestsTotal                        *prometheus.CounterVec
	TLSTotal                             *prometheus.CounterVec
	ShieldTotal                          *prometheus.CounterVec
	IPv6Total                            *prometheus.CounterVec
	ImgOptoTotal                         *prometheus.CounterVec
	ImgOptoShieldTotal                   *prometheus.CounterVec
	ImgOptoTransformTotal                *prometheus.CounterVec
	OTFPTotal                            *prometheus.CounterVec
	OTFPShieldTotal                      *prometheus.CounterVec
	OTFPTransformTotal                   *prometheus.CounterVec
	OTFPManifestTotal                    *prometheus.CounterVec
	VideoTotal                           *prometheus.CounterVec
	PCITotal                             *prometheus.CounterVec
	LoggingTotal                         *prometheus.CounterVec
	HTTP2Total                           *prometheus.CounterVec
	RespHeaderBytesTotal                 *prometheus.CounterVec
	HeaderSizeTotal                      *prometheus.CounterVec
	RespBodyBytesTotal                   *prometheus.CounterVec
	BodySizeTotal                        *prometheus.CounterVec
	ReqHeaderBytesTotal                  *prometheus.CounterVec
	BackendReqHeaderBytesTotal           *prometheus.CounterVec
	BilledHeaderBytesTotal               *prometheus.CounterVec
	BilledBodyBytesTotal                 *prometheus.CounterVec
	WAFBlockedTotal                      *prometheus.CounterVec
	WAFLoggedTotal                       *prometheus.CounterVec
	WAFPassedTotal                       *prometheus.CounterVec
	AttackReqHeaderBytesTotal            *prometheus.CounterVec
	AttackReqBodyBytesTotal              *prometheus.CounterVec
	AttackRespSynthBytesTotal            *prometheus.CounterVec
	AttackLoggedReqHeaderBytesTotal      *prometheus.CounterVec
	AttackLoggedReqBodyBytesTotal        *prometheus.CounterVec
	AttackBlockedReqHeaderBytesTotal     *prometheus.CounterVec
	AttackBlockedReqBodyBytesTotal       *prometheus.CounterVec
	AttackPassedReqHeaderBytesTotal      *prometheus.CounterVec
	AttackPassedReqBodyBytesTotal        *prometheus.CounterVec
	ShieldRespHeaderBytesTotal           *prometheus.CounterVec
	ShieldRespBodyBytesTotal             *prometheus.CounterVec
	OTFPRespHeaderBytesTotal             *prometheus.CounterVec
	OTFPRespBodyBytesTotal               *prometheus.CounterVec
	OTFPShieldRespHeaderBytesTotal       *prometheus.CounterVec
	OTFPShieldRespBodyBytesTotal         *prometheus.CounterVec
	OTFPTransformRespHeaderBytesTotal    *prometheus.CounterVec
	OTFPTransformRespBodyBytesTotal      *prometheus.CounterVec
	OTFPShieldTimeTotal                  *prometheus.CounterVec
	OTFPTransformTimeTotal               *prometheus.CounterVec
	OTFPDeliverTimeTotal                 *prometheus.CounterVec
	ImgOptoRespHeaderBytesTotal          *prometheus.CounterVec
	ImgOptoRespBodyBytesTotal            *prometheus.CounterVec
	ImgOptoShieldRespHeaderBytesTotal    *prometheus.CounterVec
	ImgOptoShieldRespBodyBytesTotal      *prometheus.CounterVec
	ImgOptoTransformRespHeaderBytesTotal *prometheus.CounterVec
	ImgOptoTransformRespBodyBytesTotal   *prometheus.CounterVec
	StatusGroupTotal                     *prometheus.CounterVec
	StatusCodeTotal                      *prometheus.CounterVec
	HitsTotal                            *prometheus.CounterVec
	MissesTotal                          *prometheus.CounterVec
	PassesTotal                          *prometheus.CounterVec
	SynthsTotal                          *prometheus.CounterVec
	ErrorsTotal                          *prometheus.CounterVec
	UncacheableTotal                     *prometheus.CounterVec
	HitsTimeTotal                        *prometheus.CounterVec
	MissTimeTotal                        *prometheus.CounterVec
	PassTimeTotal                        *prometheus.CounterVec
	MissDurationSeconds                  *prometheus.HistogramVec
	ObjectSizeBytes                      *prometheus.HistogramVec
	RecvSubTimeTotal                     *prometheus.CounterVec
	RecvSubCountTotal                    *prometheus.CounterVec
	HashSubTimeTotal                     *prometheus.CounterVec
	HashSubCountTotal                    *prometheus.CounterVec
	MissSubTimeTotal                     *prometheus.CounterVec
	MissSubCountTotal                    *prometheus.CounterVec
	FetchSubTimeTotal                    *prometheus.CounterVec
	FetchSubCountTotal                   *prometheus.CounterVec
	DeliverSubTimeTotal                  *prometheus.CounterVec
	DeliverSubCountTotal                 *prometheus.CounterVec
	HitSubTimeTotal                      *prometheus.CounterVec
	HitSubCountTotal                     *prometheus.CounterVec
	PrehashSubTimeTotal                  *prometheus.CounterVec
	PrehashSubCountTotal                 *prometheus.CounterVec
	PredeliverSubTimeTotal               *prometheus.CounterVec
	PredeliverSubCountTotal              *prometheus.CounterVec
}

// NewMetrics returns a usable set of Prometheus metrics
// that have been registered to the provided registerer.
func NewMetrics(namespace, subsystem string, r prometheus.Registerer, excludes Stringmap) (*Metrics, error) {
	var m Metrics
	m.RealtimeAPIRequestsTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "realtime_api_requests_total",
		Help: "Total requests made to the real-time stats API.",
	}, []string{"service_id", "service_name", "result"}, r, excludes)
	m.ServiceInfo = registerGauge(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "service_info",
		Help: "Static gauge with service ID, name, and version information.",
	}, []string{"service_id", "service_name", "service_version"}, r, excludes)
	m.RequestsTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "requests_total",
		Help: "Total number of requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.TLSTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_total",
		Help: "Total number of TLS requests.",
	}, []string{"service_id", "service_name", "datacenter", "tls_version"}, r, excludes)
	m.ShieldTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_total",
		Help: "Total number of shield requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.IPv6Total = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "ipv6_total",
		Help: "Total number of IPv6 requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_total",
		Help: "Total number of image optimization requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoShieldTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_total",
		Help: "Total number of image optimization shield requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoTransformTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transforms_total",
		Help: "Total number of image optimization transforms.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_total",
		Help: "Total number of on-the-fly package requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPShieldTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_total",
		Help: "Total number of on-the-fly package shield requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPTransformTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transforms_total",
		Help: "Total number of on-the-fly package transforms.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPManifestTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_manifests_total",
		Help: "Total number of on-the-fly package manifests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.VideoTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "video_total",
		Help: "Total number of video requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PCITotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pci_total",
		Help: "Total number of PCI requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.LoggingTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "logging_total",
		Help: "Total number of logging requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HTTP2Total = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "http2_total",
		Help: "Total number of HTTP2 requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.RespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_header_bytes_total",
		Help: "Total size of response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HeaderSizeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "header_size_total",
		Help: "Total size of headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.RespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_body_bytes_total",
		Help: "Total size of response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.BodySizeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "body_size_total",
		Help: "Total size of bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "req_header_bytes_total",
		Help: "Total size of request headers, in bytes",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.BackendReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "bereq_header_bytes_total",
		Help: "Total size of backend headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.BilledHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_header_bytes_total",
		Help: "Total count of billed headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.BilledBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_body_bytes_total",
		Help: "Total count of billed bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.WAFBlockedTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_blocked_total",
		Help: "Total number of WAF blocked requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.WAFLoggedTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_logged_total",
		Help: "Total number of WAF logged requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.WAFPassedTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_passed_total",
		Help: "Total number of WAF passed requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackReqBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackRespSynthBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_resp_synth_bytes_total",
		Help: "Total count of 'attack' classified synth responses, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackLoggedReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers logged, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackLoggedReqBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies logged, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackBlockedReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers blocked, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackBlockedReqBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies blocked, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackPassedReqHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_header_bytes_total",
		Help: "Total size of 'attack' classified request headers passed, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.AttackPassedReqBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_body_bytes_total",
		Help: "Total size of 'attack' classified request bodies passed, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ShieldRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_header_bytes_total",
		Help: "Total size of shielded response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ShieldRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_body_bytes_total",
		Help: "Total size of shielded response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_header_bytes_total",
		Help: "Total size of on-the-fly package response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_body_bytes_total",
		Help: "Total size of on-the-fly package response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPShieldRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_header_bytes_total",
		Help: "Total size of on-the-fly package shield response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPShieldRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_body_bytes_total",
		Help: "Total size of on-the-fly package shield response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPTransformRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_header_bytes_total",
		Help: "Total size of on-the-fly package transform response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPTransformRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_body_bytes_total",
		Help: "Total size of on-the-fly package transform response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPShieldTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_time_total",
		Help: "Total time spent in on-the-fly package shield.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPTransformTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_time_total",
		Help: "Total time spent in on-the-fly package transforms.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.OTFPDeliverTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_deliver_time_total",
		Help: "Total time spent in on-the-fly package delivery.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_header_bytes_total",
		Help: "Total count of image optimization response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_body_bytes_total",
		Help: "Total count of image optimization response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoShieldRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_header_bytes_total",
		Help: "Total count of image optimization shield response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoShieldRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_body_bytes_total",
		Help: "Total count of image optimization shield response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoTransformRespHeaderBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_header_bytes_total",
		Help: "Total count of image optimization transform response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ImgOptoTransformRespBodyBytesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_body_bytes_total",
		Help: "Total count of image optimization transform response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.StatusGroupTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_group_total",
		Help: "Total count of requests, bucketed into status groups e.g. 1xx, 2xx.",
	}, []string{"service_id", "service_name", "datacenter", "status_group"}, r, excludes) // e.g. 1xx, 2xx
	m.StatusCodeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_code_total",
		Help: "Total count of requests, bucketed into individual status codes.",
	}, []string{"service_id", "service_name", "datacenter", "status_code"}, r, excludes) // e.g. 200, 404
	m.HitsTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_total",
		Help: "Total count of hits.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.MissesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_total",
		Help: "Total count of misses.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PassesTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_total",
		Help: "Total count of passes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.SynthsTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "synth_total",
		Help: "Total count of synths.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ErrorsTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "errors_total",
		Help: "Total count of errors.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.UncacheableTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "uncacheable_total",
		Help: "Total count of uncachable responses.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HitsTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_time_total",
		Help: "Total time spent serving hits.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.MissTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_time_total",
		Help: "Total time spent serving misses.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PassTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_time_total",
		Help: "Total time spent serving passes.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.MissDurationSeconds = registerHistogram(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "miss_duration_seconds",
		Help:    "Total time spent serving misses.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2, 4, 8, 16, 32, 60},
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.ObjectSizeBytes = registerHistogram(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "object_size_bytes",
		Help:    "Size of objects served in bytes.",
		Buckets: []float64{1 * 1024, 10 * 1024, 100 * 1024, 1 * 1000 * 1024, 10 * 1000 * 1024, 100 * 1000 * 1024, 1000 * 1000 * 1024},
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.RecvSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_time_total",
		Help: "Total receive sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.RecvSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_count_total",
		Help: "Total receive sub requests.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HashSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_time_total",
		Help: "Total hash sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HashSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_count_total",
		Help: "Total hash sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.MissSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_time_total",
		Help: "Total miss sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.MissSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_count_total",
		Help: "Total miss sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.FetchSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_time_total",
		Help: "Total fetch sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.FetchSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_count_total",
		Help: "Total fetch sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.DeliverSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_time_total",
		Help: "Total deliver sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.DeliverSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_count_total",
		Help: "Total deliver sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HitSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_time_total",
		Help: "Total hit sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.HitSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_count_total",
		Help: "Total hit sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PrehashSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_time_total",
		Help: "Total prehash sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PrehashSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_count_total",
		Help: "Total prehash sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PredeliverSubTimeTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_time_total",
		Help: "Total predeliver sub time.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)
	m.PredeliverSubCountTotal = registerCounter(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_count_total",
		Help: "Total predeliver sub count.",
	}, []string{"service_id", "service_name", "datacenter"}, r, excludes)

	if len(excludes) != 0 {
		return nil, fmt.Errorf("the following excludes don't seem to refer to existing metrics: %s", excludes.String())
	}

	return &m, nil
}

func registerCounter(opts prometheus.CounterOpts, labels []string, r prometheus.Registerer, excludes map[string]bool) *prometheus.CounterVec {
	vec := prometheus.NewCounterVec(opts, labels)
	if isExcluded(excludes, opts.Name) {
		// don't register if this metric is excluded
		return vec
	}
	if err := r.Register(vec); err != nil {
		panic(errors.Wrapf(err, "error registering counter metric %s", opts.Name))
	}
	return vec
}

func registerGauge(opts prometheus.GaugeOpts, labels []string, r prometheus.Registerer, excludes map[string]bool) *prometheus.GaugeVec {
	vec := prometheus.NewGaugeVec(opts, labels)
	if isExcluded(excludes, opts.Name) {
		// don't register if this metric is excluded
		return vec
	}
	if err := r.Register(vec); err != nil {
		panic(errors.Wrapf(err, "error registering counter metric %s", opts.Name))
	}
	return vec
}

func registerHistogram(opts prometheus.HistogramOpts, labels []string, r prometheus.Registerer, excludes map[string]bool) *prometheus.HistogramVec {
	vec := prometheus.NewHistogramVec(opts, labels)
	if isExcluded(excludes, opts.Name) {
		// don't register if this metric is excluded
		return vec
	}
	if err := r.Register(vec); err != nil {
		panic(errors.Wrapf(err, "error registering counter metric %s", opts.Name))
	}
	return vec
}

func isExcluded(excludes map[string]bool, name string) bool {
	if excludes[name] {
		// the map should end up empty after all excludes have been processed
		delete(excludes, name)
		return true
	}
	return false
}

type Stringmap map[string]bool

func (sm *Stringmap) Set(s string) error {
	stringmap := *sm
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			stringmap[v] = true
		}
	}
	return nil
}

func (sm *Stringmap) String() string {
	var sx []string
	for v, _ := range *sm {
		sx = append(sx, v)
	}
	return strings.Join(sx, ",")
}
