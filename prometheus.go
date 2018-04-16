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
	objectSize1kTotal                    *prometheus.CounterVec
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
		Help: "TODO",
	}, []string{"datacenter"})
	m.tlsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.shieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.iPv6Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "ipv6_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoShieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoTransformTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transforms_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpShieldTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpTransformTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transforms_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpManifestTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_manifests_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.videoTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "video_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.pciTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pci_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.loggingTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "logging_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.http2Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "http2_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.respHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.headerSizeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "header_size_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.respBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.bodySizeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "body_size_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.reqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "req_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.backendReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "bereq_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.billedHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.billedBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "billed_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.wAFBlockedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_blocked_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.wAFLoggedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_logged_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.wAFPassedTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "waf_passed_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_req_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackRespSynthBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_resp_synth_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackLoggedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackLoggedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_logged_req_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackBlockedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackBlockedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_blocked_req_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackPassedReqHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.attackPassedReqBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "attack_passed_req_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.shieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.shieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "shield_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpShieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpShieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpTransformRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpTransformRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpShieldTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_shield_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpTransformTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_transform_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.otfpDeliverTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "otfp_deliver_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoShieldRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoShieldRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_shield_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoTransformRespHeaderBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_header_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.imgOptoTransformRespBodyBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "imgopto_transform_resp_body_bytes_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.statusGroupTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_group_total",
		Help: "TODO",
	}, []string{"datacenter", "status_group"}) // e.g. 1xx, 2xx
	m.statusCodeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "status_code_total",
		Help: "TODO",
	}, []string{"datacenter", "status_code"}) // e.g. 200, 404
	m.hitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.missesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.passesTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.synthsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "synth_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.errorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "errors_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.uncacheableTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "uncacheable_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.hitsTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hits_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.missTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.passTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "pass_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.missDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem,
		Name:    "miss_duration_seconds",
		Help:    "TODO",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2, 4, 8, 16, 32, 60},
	}, []string{"datacenter"})
	m.tlsv12Total = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "tls_v12_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.objectSize1kTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "object_size_1k_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.recvSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.recvSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "recv_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.hashSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.hashSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hash_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.missSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.missSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "miss_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.fetchSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.fetchSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "fetch_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.deliverSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.deliverSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "deliver_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.hitSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.hitSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "hit_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.prehashSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.prehashSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "prehash_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.predeliverSubTimeTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_time_total",
		Help: "TODO",
	}, []string{"datacenter"})
	m.predeliverSubCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem,
		Name: "predeliver_sub_count_total",
		Help: "TODO",
	}, []string{"datacenter"})
}
