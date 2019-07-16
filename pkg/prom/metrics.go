package prom

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics are the concrete Prometheus metrics that get updated with data
// retrieved from the real-time stats API. The same set of metrics are updated
// for all service IDs, only the labels change.
type Metrics struct {
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
func NewMetrics(namespace, subsystem string, r prometheus.Registerer) (*Metrics, error) {
	var m Metrics
	m.ServiceInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "service_info",
		Help: "Static gauge with service ID, name, and version information.",
	}, []string{"service_id", "service_name", "service_version"})
	m.RequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "requests_total",
		Help: "Total number of requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.TLSTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_total",
		Help: "Total number of TLS requests.",
	}, []string{"service_id", "service_name", "datacenter", "tls_version"})
	m.ShieldTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_total",
		Help: "Total number of shield requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.IPv6Total = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "ipv6_total",
		Help: "Total number of IPv6 requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_total",
		Help: "Total number of image optimization requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoShieldTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_total",
		Help: "Total number of image optimization shield requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoTransformTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transforms_total",
		Help: "Total number of image optimization transforms.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_total",
		Help: "Total number of on-the-fly package requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPShieldTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_total",
		Help: "Total number of on-the-fly package shield requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPTransformTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transforms_total",
		Help: "Total number of on-the-fly package transforms.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPManifestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_manifests_total",
		Help: "Total number of on-the-fly package manifests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.VideoTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "video_total",
		Help: "Total number of video requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PCITotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pci_total",
		Help: "Total number of PCI requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.LoggingTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "logging_total",
		Help: "Total number of logging requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HTTP2Total = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "http2_total",
		Help: "Total number of HTTP2 requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.RespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_header_bytes_total",
		Help: "Total size of response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HeaderSizeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "header_size_total",
		Help: "Total size of headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.RespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_body_bytes_total",
		Help: "Total size of response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.BodySizeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "body_size_total",
		Help: "Total size of bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "req_header_bytes_total",
		Help: "Total size of request headers, in bytes",
	}, []string{"service_id", "service_name", "datacenter"})
	m.BackendReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "bereq_header_bytes_total",
		Help: "Total size of backend headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.BilledHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_header_bytes_total",
		Help: "Total count of billed headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.BilledBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_body_bytes_total",
		Help: "Total count of billed bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.WAFBlockedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_blocked_total",
		Help: "Total number of WAF blocked requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.WAFLoggedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_logged_total",
		Help: "Total number of WAF logged requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.WAFPassedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_passed_total",
		Help: "Total number of WAF passed requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackReqBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackRespSynthBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_resp_synth_bytes_total",
		Help: "Total count of 'attack' classified synth responses, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackLoggedReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers logged, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackLoggedReqBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies logged, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackBlockedReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers blocked, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackBlockedReqBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies blocked, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackPassedReqHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_header_bytes_total",
		Help: "Total size of 'attack' classified request headers passed, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.AttackPassedReqBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_body_bytes_total",
		Help: "Total size of 'attack' classified request bodies passed, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ShieldRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_header_bytes_total",
		Help: "Total size of shielded response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ShieldRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_body_bytes_total",
		Help: "Total size of shielded response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_header_bytes_total",
		Help: "Total size of on-the-fly package response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_body_bytes_total",
		Help: "Total size of on-the-fly package response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPShieldRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_header_bytes_total",
		Help: "Total size of on-the-fly package shield response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPShieldRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_body_bytes_total",
		Help: "Total size of on-the-fly package shield response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPTransformRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_header_bytes_total",
		Help: "Total size of on-the-fly package transform response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPTransformRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_body_bytes_total",
		Help: "Total size of on-the-fly package transform response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPShieldTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_time_total",
		Help: "Total time spent in on-the-fly package shield.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPTransformTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_time_total",
		Help: "Total time spent in on-the-fly package transforms.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.OTFPDeliverTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_deliver_time_total",
		Help: "Total time spent in on-the-fly package delivery.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_header_bytes_total",
		Help: "Total count of image optimization response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_body_bytes_total",
		Help: "Total count of image optimization response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoShieldRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_header_bytes_total",
		Help: "Total count of image optimization shield response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoShieldRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_body_bytes_total",
		Help: "Total count of image optimization shield response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoTransformRespHeaderBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_header_bytes_total",
		Help: "Total count of image optimization transform response headers, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ImgOptoTransformRespBodyBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_body_bytes_total",
		Help: "Total count of image optimization transform response bodies, in bytes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.StatusGroupTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_group_total",
		Help: "Total count of requests, bucketed into status groups e.g. 1xx, 2xx.",
	}, []string{"service_id", "service_name", "datacenter", "status_group"}) // e.g. 1xx, 2xx
	m.StatusCodeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_code_total",
		Help: "Total count of requests, bucketed into individual status codes.",
	}, []string{"service_id", "service_name", "datacenter", "status_code"}) // e.g. 200, 404
	m.HitsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_total",
		Help: "Total count of hits.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.MissesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_total",
		Help: "Total count of misses.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PassesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_total",
		Help: "Total count of passes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.SynthsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "synth_total",
		Help: "Total count of synths.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.ErrorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "errors_total",
		Help: "Total count of errors.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.UncacheableTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "uncacheable_total",
		Help: "Total count of uncachable responses.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HitsTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_time_total",
		Help: "Total time spent serving hits.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.MissTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_time_total",
		Help: "Total time spent serving misses.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PassTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_time_total",
		Help: "Total time spent serving passes.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.MissDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "miss_duration_seconds",
		Help:    "Total time spent serving misses.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2, 4, 8, 16, 32, 60},
	}, []string{"service_id", "service_name", "datacenter"})
	m.ObjectSizeBytes = prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "object_size_bytes",
		Help:    "Size of objects served in bytes.",
		Buckets: []float64{1 * 1024, 10 * 1024, 100 * 1024, 1 * 1000 * 1024, 10 * 1000 * 1024, 100 * 1000 * 1024, 1000 * 1000 * 1024},
	}, []string{"service_id", "service_name", "datacenter"})
	m.RecvSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_time_total",
		Help: "Total receive sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.RecvSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_count_total",
		Help: "Total receive sub requests.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HashSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_time_total",
		Help: "Total hash sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HashSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_count_total",
		Help: "Total hash sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.MissSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_time_total",
		Help: "Total miss sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.MissSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_count_total",
		Help: "Total miss sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.FetchSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_time_total",
		Help: "Total fetch sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.FetchSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_count_total",
		Help: "Total fetch sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.DeliverSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_time_total",
		Help: "Total deliver sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.DeliverSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_count_total",
		Help: "Total deliver sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HitSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_time_total",
		Help: "Total hit sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.HitSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_count_total",
		Help: "Total hit sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PrehashSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_time_total",
		Help: "Total prehash sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PrehashSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_count_total",
		Help: "Total prehash sub count.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PredeliverSubTimeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_time_total",
		Help: "Total predeliver sub time.",
	}, []string{"service_id", "service_name", "datacenter"})
	m.PredeliverSubCountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_count_total",
		Help: "Total predeliver sub count.",
	}, []string{"service_id", "service_name", "datacenter"})

	for i, v := 0, reflect.ValueOf(m); i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(prometheus.Collector)
		if !ok {
			panic(fmt.Sprintf("programmer error: field %d/%d in prom.Metrics isn't a prometheus.Collector", i+1, v.NumField()))
		}
		if err := r.Register(c); err != nil {
			return nil, errors.Wrapf(err, "error registering metric %d/%d", i+1, v.NumField())
		}
	}

	return &m, nil
}
