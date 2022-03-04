package realtime

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics collects all of the Prometheus metrics that map to real-time stats.
type Metrics struct {
	RealtimeAPIRequestsTotal             *prometheus.CounterVec
	AttackBlockedReqBodyBytesTotal       *prometheus.CounterVec
	AttackBlockedReqHeaderBytesTotal     *prometheus.CounterVec
	AttackLoggedReqBodyBytesTotal        *prometheus.CounterVec
	AttackLoggedReqHeaderBytesTotal      *prometheus.CounterVec
	AttackPassedReqBodyBytesTotal        *prometheus.CounterVec
	AttackPassedReqHeaderBytesTotal      *prometheus.CounterVec
	AttackReqBodyBytesTotal              *prometheus.CounterVec
	AttackReqHeaderBytesTotal            *prometheus.CounterVec
	AttackRespSynthBytesTotal            *prometheus.CounterVec
	BackendReqBodyBytesTotal             *prometheus.CounterVec
	BackendReqHeaderBytesTotal           *prometheus.CounterVec
	BilledBodyBytesTotal                 *prometheus.CounterVec
	BilledHeaderBytesTotal               *prometheus.CounterVec
	BilledTotal                          *prometheus.CounterVec
	BlacklistedTotal                     *prometheus.CounterVec
	BodySizeTotal                        *prometheus.CounterVec
	ComputeBackendReqBodyBytesTotal      *prometheus.CounterVec
	ComputeBackendReqErrorsTotal         *prometheus.CounterVec
	ComputeBackendReqHeaderBytesTotal    *prometheus.CounterVec
	ComputeBackendReqTotal               *prometheus.CounterVec
	ComputeBackendRespBodyBytesTotal     *prometheus.CounterVec
	ComputeBackendRespHeaderBytesTotal   *prometheus.CounterVec
	ComputeExecutionTimeTotal            *prometheus.CounterVec
	ComputeGlobalsLimitExceededTotal     *prometheus.CounterVec
	ComputeGuestErrorsTotal              *prometheus.CounterVec
	ComputeHeapLimitExceededTotal        *prometheus.CounterVec
	ComputeRAMUsedBytesTotal             *prometheus.CounterVec
	ComputeReqBodyBytesTotal             *prometheus.CounterVec
	ComputeReqHeaderBytesTotal           *prometheus.CounterVec
	ComputeRequestsTotal                 *prometheus.CounterVec
	ComputeRequestTimeTotal              *prometheus.CounterVec
	ComputeResourceLimitExceedTotal      *prometheus.CounterVec
	ComputeRespBodyBytesTotal            *prometheus.CounterVec
	ComputeRespHeaderBytesTotal          *prometheus.CounterVec
	ComputeRespStatusTotal               *prometheus.CounterVec
	ComputeRuntimeErrorsTotal            *prometheus.CounterVec
	ComputeStackLimitExceededTotal       *prometheus.CounterVec
	DeliverSubCountTotal                 *prometheus.CounterVec
	DeliverSubTimeTotal                  *prometheus.CounterVec
	EdgeRespBodyBytesTotal               *prometheus.CounterVec
	EdgeRespHeaderBytesTotal             *prometheus.CounterVec
	EdgeTotal                            *prometheus.CounterVec
	ErrorsTotal                          *prometheus.CounterVec
	ErrorSubCountTotal                   *prometheus.CounterVec
	ErrorSubTimeTotal                    *prometheus.CounterVec
	FetchSubCountTotal                   *prometheus.CounterVec
	FetchSubTimeTotal                    *prometheus.CounterVec
	HashSubCountTotal                    *prometheus.CounterVec
	HashSubTimeTotal                     *prometheus.CounterVec
	HeaderSizeTotal                      *prometheus.CounterVec
	HitRespBodyBytesTotal                *prometheus.CounterVec
	HitsTimeTotal                        *prometheus.CounterVec
	HitsTotal                            *prometheus.CounterVec
	HitSubCountTotal                     *prometheus.CounterVec
	HitSubTimeTotal                      *prometheus.CounterVec
	HTTP2Total                           *prometheus.CounterVec
	ImgOptoRespBodyBytesTotal            *prometheus.CounterVec
	ImgOptoRespHeaderBytesTotal          *prometheus.CounterVec
	ImgOptoShieldRespBodyBytesTotal      *prometheus.CounterVec
	ImgOptoShieldRespHeaderBytesTotal    *prometheus.CounterVec
	ImgOptoShieldTotal                   *prometheus.CounterVec
	ImgOptoTotal                         *prometheus.CounterVec
	ImgOptoTransformRespBodyBytesTotal   *prometheus.CounterVec
	ImgOptoTransformRespHeaderBytesTotal *prometheus.CounterVec
	ImgOptoTransformTotal                *prometheus.CounterVec
	ImgVideoFramesTotal                  *prometheus.CounterVec
	ImgVideoRespBodyBytesTotal           *prometheus.CounterVec
	ImgVideoRespHeaderBytesTotal         *prometheus.CounterVec
	ImgVideoShieldFramesTotal            *prometheus.CounterVec
	ImgVideoShieldRespBodyBytesTotal     *prometheus.CounterVec
	ImgVideoShieldRespHeaderBytesTotal   *prometheus.CounterVec
	ImgVideoShieldTotal                  *prometheus.CounterVec
	ImgVideoTotal                        *prometheus.CounterVec
	IPv6Total                            *prometheus.CounterVec
	LogBytesTotal                        *prometheus.CounterVec
	LoggingTotal                         *prometheus.CounterVec
	MissDurationSeconds                  *prometheus.HistogramVec
	MissesTotal                          *prometheus.CounterVec
	MissRespBodyBytesTotal               *prometheus.CounterVec
	MissSubCountTotal                    *prometheus.CounterVec
	MissSubTimeTotal                     *prometheus.CounterVec
	MissTimeTotal                        *prometheus.CounterVec
	ObjectSizeBytes                      *prometheus.HistogramVec
	OriginFetchBodyBytesTotal            *prometheus.CounterVec
	OriginFetchesTotal                   *prometheus.CounterVec
	OriginFetchHeaderBytesTotal          *prometheus.CounterVec
	OriginFetchRespBodyBytesTotal        *prometheus.CounterVec
	OriginFetchRespHeaderBytesTotal      *prometheus.CounterVec
	OriginRevalidationsTotal             *prometheus.CounterVec
	OTFPDeliverTimeTotal                 *prometheus.CounterVec
	OTFPManifestTotal                    *prometheus.CounterVec
	OTFPRespBodyBytesTotal               *prometheus.CounterVec
	OTFPRespHeaderBytesTotal             *prometheus.CounterVec
	OTFPShieldRespBodyBytesTotal         *prometheus.CounterVec
	OTFPShieldRespHeaderBytesTotal       *prometheus.CounterVec
	OTFPShieldTimeTotal                  *prometheus.CounterVec
	OTFPShieldTotal                      *prometheus.CounterVec
	OTFPTotal                            *prometheus.CounterVec
	OTFPTransformRespBodyBytesTotal      *prometheus.CounterVec
	OTFPTransformRespHeaderBytesTotal    *prometheus.CounterVec
	OTFPTransformTimeTotal               *prometheus.CounterVec
	OTFPTransformTotal                   *prometheus.CounterVec
	PassesTotal                          *prometheus.CounterVec
	PassRespBodyBytesTotal               *prometheus.CounterVec
	PassSubCountTotal                    *prometheus.CounterVec
	PassSubTimeTotal                     *prometheus.CounterVec
	PassTimeTotal                        *prometheus.CounterVec
	PCITotal                             *prometheus.CounterVec
	Pipe                                 *prometheus.CounterVec
	PipeSubCountTotal                    *prometheus.CounterVec
	PipeSubTimeTotal                     *prometheus.CounterVec
	PredeliverSubCountTotal              *prometheus.CounterVec
	PredeliverSubTimeTotal               *prometheus.CounterVec
	PrehashSubCountTotal                 *prometheus.CounterVec
	PrehashSubTimeTotal                  *prometheus.CounterVec
	RecvSubCountTotal                    *prometheus.CounterVec
	RecvSubTimeTotal                     *prometheus.CounterVec
	ReqBodyBytesTotal                    *prometheus.CounterVec
	ReqHeaderBytesTotal                  *prometheus.CounterVec
	RequestsTotal                        *prometheus.CounterVec
	RespBodyBytesTotal                   *prometheus.CounterVec
	RespHeaderBytesTotal                 *prometheus.CounterVec
	RestartTotal                         *prometheus.CounterVec
	SegBlockOriginFetchesTotal           *prometheus.CounterVec
	SegBlockShieldFetchesTotal           *prometheus.CounterVec
	ShieldFetchBodyBytesTotal            *prometheus.CounterVec
	ShieldFetchesTotal                   *prometheus.CounterVec
	ShieldFetchHeaderBytesTotal          *prometheus.CounterVec
	ShieldFetchRespBodyBytesTotal        *prometheus.CounterVec
	ShieldFetchRespHeaderBytesTotal      *prometheus.CounterVec
	ShieldRespBodyBytesTotal             *prometheus.CounterVec
	ShieldRespHeaderBytesTotal           *prometheus.CounterVec
	ShieldRevalidationsTotal             *prometheus.CounterVec
	ShieldTotal                          *prometheus.CounterVec
	StatusCodeTotal                      *prometheus.CounterVec
	StatusGroupTotal                     *prometheus.CounterVec
	SynthsTotal                          *prometheus.CounterVec
	TLSTotal                             *prometheus.CounterVec
	UncacheableTotal                     *prometheus.CounterVec
	VideoTotal                           *prometheus.CounterVec
	WAFBlockedTotal                      *prometheus.CounterVec
	WAFLoggedTotal                       *prometheus.CounterVec
	WAFPassedTotal                       *prometheus.CounterVec
}

// NewMetrics returns a new set of metrics registered to the registerer.
// Only metrics whose names pass the name filter are registered.
func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {
	m := Metrics{
		RealtimeAPIRequestsTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "realtime_api_requests_total", Help: "Total requests made to the real-time stats API."}, []string{"service_id", "service_name", "result"}),
		AttackBlockedReqBodyBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_blocked_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was blocked."}, []string{"service_id", "service_name", "datacenter"}),
		AttackBlockedReqHeaderBytesTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_blocked_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was blocked."}, []string{"service_id", "service_name", "datacenter"}),
		AttackLoggedReqBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_logged_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was logged."}, []string{"service_id", "service_name", "datacenter"}),
		AttackLoggedReqHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_logged_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was logged."}, []string{"service_id", "service_name", "datacenter"}),
		AttackPassedReqBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_passed_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was passed."}, []string{"service_id", "service_name", "datacenter"}),
		AttackPassedReqHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_passed_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was passed."}, []string{"service_id", "service_name", "datacenter"}),
		AttackReqBodyBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule."}, []string{"service_id", "service_name", "datacenter"}),
		AttackReqHeaderBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule."}, []string{"service_id", "service_name", "datacenter"}),
		AttackRespSynthBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_resp_synth_bytes_total", Help: "Total bytes delivered for requests that triggered a WAF rule and returned a synthetic response."}, []string{"service_id", "service_name", "datacenter"}),
		BackendReqBodyBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_body_bytes_total", Help: "Total body bytes sent to origin."}, []string{"service_id", "service_name", "datacenter"}),
		BackendReqHeaderBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_header_bytes_total", Help: "Total header bytes sent to origin."}, []string{"service_id", "service_name", "datacenter"}),
		BilledBodyBytesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "billed_body_bytes_total", Help: "TODO"}, []string{"service_id", "service_name", "datacenter"}),
		BilledHeaderBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "billed_header_bytes_total", Help: "TODO"}, []string{"service_id", "service_name", "datacenter"}),
		BilledTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "billed_total", Help: "TODO"}, []string{"service_id", "service_name", "datacenter"}),
		BlacklistedTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "blacklist_total", Help: "TODO"}, []string{"service_id", "service_name", "datacenter"}),
		BodySizeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "body_size_total", Help: "Total body bytes delivered (alias for resp_body_bytes)."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendReqBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_body_bytes_total", Help: "Total body bytes sent to backends (origins) by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendReqErrorsTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_errors_total", Help: "Number of backend request errors, including timeouts."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendReqHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_header_bytes_total", Help: "Total header bytes sent to backends (origins) by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendReqTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_total", Help: "Number of backend requests started."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendRespBodyBytesTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_beresp_body_bytes_total", Help: "Total body bytes received from backends (origins) by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeBackendRespHeaderBytesTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_beresp_header_bytes_total", Help: "Total header bytes received from backends (origins) by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeExecutionTimeTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_execution_time_total", Help: "The amount of active CPU time used to process your requests (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeGlobalsLimitExceededTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_globals_limit_exceeded_total", Help: "Number of times a guest exceeded its globals limit."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeGuestErrorsTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_guest_errors_total", Help: "Number of times a service experienced a guest code error."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeHeapLimitExceededTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_heap_limit_exceeded_total", Help: "Number of times a guest exceeded its heap limit."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRAMUsedBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_ram_used_bytes_total", Help: "The amount of RAM used for your site by Fastly."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeReqBodyBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_req_body_bytes_total", Help: "Total body bytes received by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeReqHeaderBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_req_header_bytes_total", Help: "Total header bytes received by Compute@Edge."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRequestsTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_requests_total", Help: "The total number of requests that were received for your site by Fastly."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRequestTimeTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_request_time_total", Help: "The total amount of time used to process your requests, including active CPU time (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeResourceLimitExceedTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resource_limit_exceeded_total", Help: "Number of times a guest exceeded its resource limit, includes heap, stack, globals, and code execution timeout."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRespBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_body_bytes_total", Help: "Total body bytes sent from Compute@Edge to end user."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRespHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_header_bytes_total", Help: "Total header bytes sent from Compute@Edge to end user."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeRespStatusTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_status_total", Help: "Number of responses delivered delivered by Compute@Edge, by status code group."}, []string{"service_id", "service_name", "datacenter", "status_group"}),
		ComputeRuntimeErrorsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_runtime_errors_total", Help: "Number of times a service experienced a guest runtime error."}, []string{"service_id", "service_name", "datacenter"}),
		ComputeStackLimitExceededTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_stack_limit_exceeded_total", Help: "Number of times a guest exceeded its stack limit."}, []string{"service_id", "service_name", "datacenter"}),
		DeliverSubCountTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "deliver_sub_count_total", Help: "Number of executions of the 'deliver' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		DeliverSubTimeTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "deliver_sub_time_total", Help: "Time spent inside the 'deliver' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		EdgeRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_body_bytes_total", Help: "Total body bytes delivered from Fastly to the end user."}, []string{"service_id", "service_name", "datacenter"}),
		EdgeRespHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_header_bytes_total", Help: "Total header bytes delivered from Fastly to the end user."}, []string{"service_id", "service_name", "datacenter"}),
		EdgeTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_total", Help: "Number of requests sent by end users to Fastly."}, []string{"service_id", "service_name", "datacenter"}),
		ErrorsTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "errors_total", Help: "Number of cache errors."}, []string{"service_id", "service_name", "datacenter"}),
		ErrorSubCountTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "error_sub_count_total", Help: "Number of executions of the 'error' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		ErrorSubTimeTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "error_sub_time_total", Help: "Time spent inside the 'error' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		FetchSubCountTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fetch_sub_count_total", Help: "Number of executions of the 'fetch' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		FetchSubTimeTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fetch_sub_time_total", Help: "Time spent inside the 'fetch' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		HashSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hash_sub_count_total", Help: "Number of executions of the 'hash' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		HashSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hash_sub_time_total", Help: "Time spent inside the 'hash' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		HeaderSizeTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "header_size_total", Help: "Total header bytes delivered (alias for resp_header_bytes)."}, []string{"service_id", "service_name", "datacenter"}),
		HitRespBodyBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_resp_body_bytes_total", Help: "Total body bytes delivered for cache hits."}, []string{"service_id", "service_name", "datacenter"}),
		HitsTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hits_time_total", Help: "Total amount of time spent processing cache hits (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		HitsTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hits_total", Help: "Number of cache hits."}, []string{"service_id", "service_name", "datacenter"}),
		HitSubCountTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_sub_count_total", Help: "Number of executions of the 'hit' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		HitSubTimeTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_sub_time_total", Help: "Time spent inside the 'hit' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		HTTP2Total:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "http2_total", Help: "Number of requests received over HTTP2."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoRespBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_resp_body_bytes_total", Help: "Total body bytes delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoRespHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_resp_header_bytes_total", Help: "Total header bytes delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoShieldRespBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoShieldRespHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoShieldTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_total", Help: "Number of responses delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_total", Help: "Number of responses that came from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoTransformRespBodyBytesTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transform_resp_body_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoTransformRespHeaderBytesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transform_resp_header_bytes_total", Help: "Total header bytes of transforms delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgOptoTransformTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transforms_total", Help: "Total transforms performed by the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoFramesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_frames_total", Help: "Number of video frames that came from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoRespBodyBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_resp_body_bytes_total", Help: "Total body bytes of video delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoRespHeaderBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_resp_header_bytes_total", Help: "Total header bytes of video delivered from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoShieldFramesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_frames_total", Help: "Number of video frames delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoShieldRespBodyBytesTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_resp_body_bytes_total", Help: "Total body bytes of video delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoShieldRespHeaderBytesTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_resp_header_bytes_total", Help: "Total header bytes of video delivered via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoShieldTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_total", Help: "Number of video responses that came via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		ImgVideoTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_total", Help: "Number of video responses that came via a shield from the Fastly Image Optimizer service."}, []string{"service_id", "service_name", "datacenter"}),
		IPv6Total:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ipv6_total", Help: "Number of requests that were received over IPv6."}, []string{"service_id", "service_name", "datacenter"}),
		LogBytesTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "log_bytes_total", Help: "Total log bytes sent."}, []string{"service_id", "service_name", "datacenter"}),
		LoggingTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "logging_total", Help: "Number of log lines sent."}, []string{"service_id", "service_name", "datacenter"}),
		MissDurationSeconds:                  prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_duration_seconds", Help: "Histogram of time spent processing cache misses (in seconds).", Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 4, 8, 16, 32, 60}}, []string{"service_id", "service_name", "datacenter"}),
		MissesTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_total", Help: "Number of cache misses."}, []string{"service_id", "service_name", "datacenter"}),
		MissRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_resp_body_bytes_total", Help: "Total body bytes delivered for cache misses."}, []string{"service_id", "service_name", "datacenter"}),
		MissSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_sub_count_total", Help: "Number of executions of the 'miss' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		MissSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_sub_time_total", Help: "Time spent inside the 'miss' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		MissTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_time_total", Help: "Total amount of time spent processing cache misses (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		ObjectSizeBytes:                      prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "object_size_bytes", Help: "Histogram of count of objects served, bucketed by object size range.", Buckets: []float64{1024, 10240, 102400, 1.024e+06, 1.024e+07, 1.024e+08, 1.024e+09}}, []string{"service_id", "service_name", "datacenter"}),
		OriginFetchBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_body_bytes_total", Help: "Total request body bytes sent to origin."}, []string{"service_id", "service_name", "datacenter"}),
		OriginFetchesTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetches_total", Help: "Number of requests sent to origin."}, []string{"service_id", "service_name", "datacenter"}),
		OriginFetchHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_header_bytes_total", Help: "Total request header bytes sent to origin."}, []string{"service_id", "service_name", "datacenter"}),
		OriginFetchRespBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_body_bytes_total", Help: "Total body bytes received from origin."}, []string{"service_id", "service_name", "datacenter"}),
		OriginFetchRespHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_header_bytes_total", Help: "Total header bytes received from origin."}, []string{"service_id", "service_name", "datacenter"}),
		OriginRevalidationsTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_revalidations_total", Help: "Number of responses received from origin with a 304 status code in response to an If-Modified-Since or If-None-Match request. Under regular scenarios, a revalidation will imply a cache hit. However, if using Fastly Image Optimizer or segmented caching this may result in a cache miss."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPDeliverTimeTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_total", Help: "Number of responses that came from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPManifestTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_deliver_time_total", Help: "Total amount of time spent delivering a response from the Fastly On-the-Fly Packager (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_manifests_total", Help: "Number of responses that were manifest files from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPRespHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_resp_body_bytes_total", Help: "Total body bytes delivered from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPShieldRespBodyBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_resp_header_bytes_total", Help: "Total header bytes delivered from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPShieldRespHeaderBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_total", Help: "Number of responses delivered from the Fastly On-the-Fly Packager"}, []string{"service_id", "service_name", "datacenter"}),
		OTFPShieldTimeTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield for the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPShieldTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield for the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_time_total", Help: "Total amount of time spent delivering a response via a shield from the Fastly On-the-Fly Packager (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPTransformRespBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transforms_total", Help: "Number of transforms performed by the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPTransformRespHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_resp_body_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPTransformTimeTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_resp_header_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		OTFPTransformTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_time_total", Help: "Total amount of time spent performing transforms from the Fastly On-the-Fly Packager."}, []string{"service_id", "service_name", "datacenter"}),
		PassesTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_total", Help: "Number of requests that passed through the CDN without being cached."}, []string{"service_id", "service_name", "datacenter"}),
		PassRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_resp_body_bytes_total", Help: "Total body bytes delivered for cache passes."}, []string{"service_id", "service_name", "datacenter"}),
		PassSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_sub_count_total", Help: "Number of executions of the 'pass' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		PassSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_sub_time_total", Help: "Time spent inside the 'pass' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		PassTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_time_total", Help: "Total amount of time spent processing cache passes (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		PCITotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pci_total", Help: "Number of responses with the PCI flag turned on."}, []string{"service_id", "service_name", "datacenter"}),
		Pipe:                                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe", Help: "Pipe operations performed."}, []string{"service_id", "service_name", "datacenter"}),
		PipeSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe_sub_count_total", Help: "Number of executions of the 'pipe' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		PipeSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe_sub_time_total", Help: "Time spent inside the 'pipe' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		PredeliverSubCountTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "predeliver_sub_count_total", Help: "Number of executions of the 'predeliver' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		PredeliverSubTimeTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "predeliver_sub_time_total", Help: "Time spent inside the 'predeliver' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		PrehashSubCountTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "prehash_sub_count_total", Help: "Number of executions of the 'prehash' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		PrehashSubTimeTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "prehash_sub_time_total", Help: "Time spent inside the 'prehash' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		RecvSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "recv_sub_count_total", Help: "Number of executions of the 'recv' Varnish subroutine."}, []string{"service_id", "service_name", "datacenter"}),
		RecvSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "recv_sub_time_total", Help: "Time spent inside the 'recv' Varnish subroutine (in seconds)."}, []string{"service_id", "service_name", "datacenter"}),
		ReqBodyBytesTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "req_body_bytes_total", Help: "Total body bytes received."}, []string{"service_id", "service_name", "datacenter"}),
		ReqHeaderBytesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "req_header_bytes_total", Help: "Total header bytes received."}, []string{"service_id", "service_name", "datacenter"}),
		RequestsTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "requests_total", Help: "Number of requests processed."}, []string{"service_id", "service_name", "datacenter"}),
		RespBodyBytesTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_body_bytes_total", Help: "Total body bytes delivered."}, []string{"service_id", "service_name", "datacenter"}),
		RespHeaderBytesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_header_bytes_total", Help: "Total header bytes delivered."}, []string{"service_id", "service_name", "datacenter"}),
		RestartTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "restarts_total", Help: "Number of restarts performed."}, []string{"service_id", "service_name", "datacenter"}),
		SegBlockOriginFetchesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "segblock_origin_fetches_total", Help: "Number of Range requests to origin for segments of resources when using segmented caching."}, []string{"service_id", "service_name", "datacenter"}),
		SegBlockShieldFetchesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "segblock_shield_fetches_total", Help: "Number of Range requests to a shield for segments of resources when using segmented caching."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldFetchBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_body_bytes_total", Help: "Total request body bytes sent to a shield."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldFetchesTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetches_total", Help: "Number of requests made from one Fastly data center to another, as part of shielding."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldFetchHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_header_bytes_total", Help: "Total request header bytes sent to a shield."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldFetchRespBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_resp_body_bytes_total", Help: "Total response body bytes sent from a shield to the edge."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldFetchRespHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_resp_header_bytes_total", Help: "Total response header bytes sent from a shield to the edge."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldRespBodyBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldRespHeaderBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldRevalidationsTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_revalidations_total", Help: "Number of responses received from origin with a 304 status code, in response to an If-Modified-Since or If-None-Match request to a shield. Under regular scenarios, a revalidation will imply a cache hit. However, if using segmented caching this may result in a cache miss."}, []string{"service_id", "service_name", "datacenter"}),
		ShieldTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_total", Help: "Number of requests from edge to the shield POP."}, []string{"service_id", "service_name", "datacenter"}),
		StatusCodeTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_code_total", Help: "Number of responses sent with status code 500 (Internal Server Error)."}, []string{"service_id", "service_name", "datacenter", "status_code"}),
		StatusGroupTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_group_total", Help: "Number of 'Client Error' category status codes delivered."}, []string{"service_id", "service_name", "datacenter", "status_group"}),
		SynthsTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "synth_total", Help: "TODO"}, []string{"service_id", "service_name", "datacenter"}),
		TLSTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "tls_total", Help: "Number of requests that were received over TLS."}, []string{"service_id", "service_name", "datacenter", "tls_version"}),
		UncacheableTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "uncacheable_total", Help: "Number of requests that were designated uncachable."}, []string{"service_id", "service_name", "datacenter"}),
		VideoTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "video_total", Help: "Number of responses with the video segment or video manifest MIME type (i.e., application/x-mpegurl, application/vnd.apple.mpegurl, application/f4m, application/dash+xml, application/vnd.ms-sstr+xml, ideo/mp2t, audio/aac, video/f4f, video/x-flv, video/mp4, audio/mp4)."}, []string{"service_id", "service_name", "datacenter"}),
		WAFBlockedTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_blocked_total", Help: "Number of requests that triggered a WAF rule and were blocked."}, []string{"service_id", "service_name", "datacenter"}),
		WAFLoggedTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_logged_total", Help: "Number of requests that triggered a WAF rule and were logged."}, []string{"service_id", "service_name", "datacenter"}),
		WAFPassedTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_passed_total", Help: "Number of requests that triggered a WAF rule and were passed."}, []string{"service_id", "service_name", "datacenter"}),
	}

	for i, v := 0, reflect.ValueOf(m); i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(prometheus.Collector)
		if !ok {
			panic(fmt.Errorf("field %d/%d in Metrics type isn't a prometheus.Collector", i+1, v.NumField()))
		}
		if name := getName(c); !nameFilter.Permit(name) {
			continue
		}
		if err := r.Register(c); err != nil {
			panic(fmt.Errorf("error registering metric %d/%d: %w", i+1, v.NumField(), err))
		}
	}

	return &m
}

var descNameRegex = regexp.MustCompile("fqName: \"([^\"]+)\"")

func getName(c prometheus.Collector) string {
	d := make(chan *prometheus.Desc, 1)
	c.Describe(d)
	desc := (<-d).String()
	matches := descNameRegex.FindAllStringSubmatch(desc, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1]
	}
	return ""
}
