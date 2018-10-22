package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type prometheusMetrics struct {
	requestsTotal                        *prometheus.CounterVec
	tlsTotal                             *prometheus.CounterVec
	shieldTotal                          *prometheus.CounterVec
	iPv6Total                            *prometheus.CounterVec
	imgOptoTotal                         *prometheus.CounterVec
	imgOptoShieldTotal                   *prometheus.CounterVec
	imgOptoTransformTotal                *prometheus.CounterVec
	otfpTotal                            *prometheus.CounterVec
	otfpShieldTotal                      *prometheus.CounterVec
	otfpTransformTotal                   *prometheus.CounterVec
	otfpManifestTotal                    *prometheus.CounterVec
	videoTotal                           *prometheus.CounterVec
	pciTotal                             *prometheus.CounterVec
	loggingTotal                         *prometheus.CounterVec
	http2Total                           *prometheus.CounterVec
	respHeaderBytesTotal                 *prometheus.CounterVec
	headerSizeTotal                      *prometheus.CounterVec
	respBodyBytesTotal                   *prometheus.CounterVec
	bodySizeTotal                        *prometheus.CounterVec
	reqHeaderBytesTotal                  *prometheus.CounterVec
	backendReqHeaderBytesTotal           *prometheus.CounterVec
	billedHeaderBytesTotal               *prometheus.CounterVec
	billedBodyBytesTotal                 *prometheus.CounterVec
	wAFBlockedTotal                      *prometheus.CounterVec
	wAFLoggedTotal                       *prometheus.CounterVec
	wAFPassedTotal                       *prometheus.CounterVec
	attackReqHeaderBytesTotal            *prometheus.CounterVec
	attackReqBodyBytesTotal              *prometheus.CounterVec
	attackRespSynthBytesTotal            *prometheus.CounterVec
	attackLoggedReqHeaderBytesTotal      *prometheus.CounterVec
	attackLoggedReqBodyBytesTotal        *prometheus.CounterVec
	attackBlockedReqHeaderBytesTotal     *prometheus.CounterVec
	attackBlockedReqBodyBytesTotal       *prometheus.CounterVec
	attackPassedReqHeaderBytesTotal      *prometheus.CounterVec
	attackPassedReqBodyBytesTotal        *prometheus.CounterVec
	shieldRespHeaderBytesTotal           *prometheus.CounterVec
	shieldRespBodyBytesTotal             *prometheus.CounterVec
	otfpRespHeaderBytesTotal             *prometheus.CounterVec
	otfpRespBodyBytesTotal               *prometheus.CounterVec
	otfpShieldRespHeaderBytesTotal       *prometheus.CounterVec
	otfpShieldRespBodyBytesTotal         *prometheus.CounterVec
	otfpTransformRespHeaderBytesTotal    *prometheus.CounterVec
	otfpTransformRespBodyBytesTotal      *prometheus.CounterVec
	otfpShieldTimeTotal                  *prometheus.CounterVec
	otfpTransformTimeTotal               *prometheus.CounterVec
	otfpDeliverTimeTotal                 *prometheus.CounterVec
	imgOptoRespHeaderBytesTotal          *prometheus.CounterVec
	imgOptoRespBodyBytesTotal            *prometheus.CounterVec
	imgOptoShieldRespHeaderBytesTotal    *prometheus.CounterVec
	imgOptoShieldRespBodyBytesTotal      *prometheus.CounterVec
	imgOptoTransformRespHeaderBytesTotal *prometheus.CounterVec
	imgOptoTransformRespBodyBytesTotal   *prometheus.CounterVec
	statusGroupTotal                     *prometheus.CounterVec
	statusCodeTotal                      *prometheus.CounterVec
	hitsTotal                            *prometheus.CounterVec
	missesTotal                          *prometheus.CounterVec
	passesTotal                          *prometheus.CounterVec
	synthsTotal                          *prometheus.CounterVec
	errorsTotal                          *prometheus.CounterVec
	uncacheableTotal                     *prometheus.CounterVec
	hitsTimeTotal                        *prometheus.CounterVec
	missTimeTotal                        *prometheus.CounterVec
	passTimeTotal                        *prometheus.CounterVec
	missDurationSeconds                  *prometheus.HistogramVec
	tlsv12Total                          *prometheus.CounterVec
	objectSizeBytes                      *prometheus.HistogramVec
	recvSubTimeTotal                     *prometheus.CounterVec
	recvSubCountTotal                    *prometheus.CounterVec
	hashSubTimeTotal                     *prometheus.CounterVec
	hashSubCountTotal                    *prometheus.CounterVec
	missSubTimeTotal                     *prometheus.CounterVec
	missSubCountTotal                    *prometheus.CounterVec
	fetchSubTimeTotal                    *prometheus.CounterVec
	fetchSubCountTotal                   *prometheus.CounterVec
	deliverSubTimeTotal                  *prometheus.CounterVec
	deliverSubCountTotal                 *prometheus.CounterVec
	hitSubTimeTotal                      *prometheus.CounterVec
	hitSubCountTotal                     *prometheus.CounterVec
	prehashSubTimeTotal                  *prometheus.CounterVec
	prehashSubCountTotal                 *prometheus.CounterVec
	predeliverSubTimeTotal               *prometheus.CounterVec
	predeliverSubCountTotal              *prometheus.CounterVec
}

func (m *prometheusMetrics) register(namespace, subsystem string) {
	m.requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "requests_total",
		Help: "Total number of requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.tlsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_total",
		Help: "Total number of TLS requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.shieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_total",
		Help: "Total number of shield requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.iPv6Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "ipv6_total",
		Help: "Total number of IPv6 requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_total",
		Help: "Total number of image optimization requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoShieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_total",
		Help: "Total number of image optimization shield requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoTransformTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transforms_total",
		Help: "Total number of image optimization transforms.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_total",
		Help: "Total number of on-the-fly package requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpShieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_total",
		Help: "Total number of on-the-fly package shield requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpTransformTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transforms_total",
		Help: "Total number of on-the-fly package transforms.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpManifestTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_manifests_total",
		Help: "Total number of on-the-fly package manifests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.videoTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "video_total",
		Help: "Total number of video requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.pciTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pci_total",
		Help: "Total number of PCI requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.loggingTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "logging_total",
		Help: "Total number of logging requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.http2Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "http2_total",
		Help: "Total number of HTTP2 requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.respHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_header_bytes_total",
		Help: "Total size of response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.headerSizeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "header_size_total",
		Help: "Total size of headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.respBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_body_bytes_total",
		Help: "Total size of response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.bodySizeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "body_size_total",
		Help: "Total size of bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.reqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "req_header_bytes_total",
		Help: "Total size of request headers, in bytes",
	}, []string{"service_name", "service_id", "datacenter"})
	m.backendReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "bereq_header_bytes_total",
		Help: "Total size of backend headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.billedHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_header_bytes_total",
		Help: "Total count of billed headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.billedBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_body_bytes_total",
		Help: "Total count of billed bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.wAFBlockedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_blocked_total",
		Help: "Total number of WAF blocked requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.wAFLoggedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_logged_total",
		Help: "Total number of WAF logged requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.wAFPassedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_passed_total",
		Help: "Total number of WAF passed requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackRespSynthBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_resp_synth_bytes_total",
		Help: "Total count of 'attack' classified synth responses, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackLoggedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers logged, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackLoggedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies logged, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackBlockedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_header_bytes_total",
		Help: "Total count of 'attack' classified request headers blocked, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackBlockedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_body_bytes_total",
		Help: "Total count of 'attack' classified request bodies blocked, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackPassedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_header_bytes_total",
		Help: "Total size of 'attack' classified request headers passed, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.attackPassedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_body_bytes_total",
		Help: "Total size of 'attack' classified request bodies passed, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.shieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_header_bytes_total",
		Help: "Total size of shielded response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.shieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_body_bytes_total",
		Help: "Total size of shielded response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_header_bytes_total",
		Help: "Total size of on-the-fly package response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_body_bytes_total",
		Help: "Total size of on-the-fly package response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpShieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_header_bytes_total",
		Help: "Total size of on-the-fly package shield response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpShieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_body_bytes_total",
		Help: "Total size of on-the-fly package shield response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpTransformRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_header_bytes_total",
		Help: "Total size of on-the-fly package transform response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpTransformRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_body_bytes_total",
		Help: "Total size of on-the-fly package transform response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpShieldTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_time_total",
		Help: "Total time spent in on-the-fly package shield.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpTransformTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_time_total",
		Help: "Total time spent in on-the-fly package transforms.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.otfpDeliverTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_deliver_time_total",
		Help: "Total time spent in on-the-fly package delivery.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_header_bytes_total",
		Help: "Total count of image optimization response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_body_bytes_total",
		Help: "Total count of image optimization response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoShieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_header_bytes_total",
		Help: "Total count of image optimization shield response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoShieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_body_bytes_total",
		Help: "Total count of image optimization shield response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoTransformRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_header_bytes_total",
		Help: "Total count of image optimization transform response headers, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.imgOptoTransformRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_body_bytes_total",
		Help: "Total count of image optimization transform response bodies, in bytes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.statusGroupTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_group_total",
		Help: "Total count of requests, bucketed into status groups e.g. 1xx, 2xx.",
	}, []string{"service_name", "service_id", "datacenter", "status_group"}) // e.g. 1xx, 2xx
	m.statusCodeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_code_total",
		Help: "Total count of requests, bucketed into individual status codes.",
	}, []string{"service_name", "service_id", "datacenter", "status_code"}) // e.g. 200, 404
	m.hitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_total",
		Help: "Total count of hits.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.missesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_total",
		Help: "Total count of misses.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.passesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_total",
		Help: "Total count of passes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.synthsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "synth_total",
		Help: "Total count of synths.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.errorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "errors_total",
		Help: "Total count of errors.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.uncacheableTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "uncacheable_total",
		Help: "Total count of uncachable responses.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.hitsTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_time_total",
		Help: "Total time spent serving hits.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.missTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_time_total",
		Help: "Total time spent serving misses.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.passTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_time_total",
		Help: "Total time spent serving passes.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.missDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "miss_duration_seconds",
		Help:    "Total time spent serving misses.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2, 4, 8, 16, 32, 60},
	}, []string{"service_name", "service_id", "datacenter"})
	m.tlsv12Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_v12_total",
		Help: "Total number of TLS v1.2 requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.objectSizeBytes = promauto.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "object_size_bytes",
		Help:    "Size of objects served in bytes.",
		Buckets: []float64{1 * 1024, 10 * 1024, 100 * 1024, 1 * 1000 * 1024, 10 * 1000 * 1024, 100 * 1000 * 1024, 1000 * 1000 * 1024},
	}, []string{"service_name", "service_id", "datacenter"})
	m.recvSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_time_total",
		Help: "Total receive sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.recvSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_count_total",
		Help: "Total receive sub requests.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.hashSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_time_total",
		Help: "Total hash sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.hashSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_count_total",
		Help: "Tothash al sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.missSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_time_total",
		Help: "Total miss sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.missSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_count_total",
		Help: "Totmiss al sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.fetchSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_time_total",
		Help: "Total fetch sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.fetchSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_count_total",
		Help: "Totafetch l sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.deliverSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_time_total",
		Help: "Total deliver sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.deliverSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_count_total",
		Help: "Total deliver sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.hitSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_time_total",
		Help: "Total hit sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.hitSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_count_total",
		Help: "Tohit tal sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.prehashSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_time_total",
		Help: "Total prehash sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.prehashSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_count_total",
		Help: "Total prehash sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.predeliverSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_time_total",
		Help: "Total predeliver sub time.",
	}, []string{"service_name", "service_id", "datacenter"})
	m.predeliverSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_count_total",
		Help: "Total predeliver sub count.",
	}, []string{"service_name", "service_id", "datacenter"})
}
