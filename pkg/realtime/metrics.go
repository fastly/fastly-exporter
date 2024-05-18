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
	AttackBlockedReqBodyBytesTotal          *prometheus.CounterVec
	AttackBlockedReqHeaderBytesTotal        *prometheus.CounterVec
	AttackLoggedReqBodyBytesTotal           *prometheus.CounterVec
	AttackLoggedReqHeaderBytesTotal         *prometheus.CounterVec
	AttackPassedReqBodyBytesTotal           *prometheus.CounterVec
	AttackPassedReqHeaderBytesTotal         *prometheus.CounterVec
	AttackReqBodyBytesTotal                 *prometheus.CounterVec
	AttackReqHeaderBytesTotal               *prometheus.CounterVec
	AttackRespSynthBytesTotal               *prometheus.CounterVec
	BackendReqBodyBytesTotal                *prometheus.CounterVec
	BackendReqHeaderBytesTotal              *prometheus.CounterVec
	BlacklistedTotal                        *prometheus.CounterVec
	BodySizeTotal                           *prometheus.CounterVec
	BotChallengeCompleteTokensCheckedTotal  *prometheus.CounterVec
	BotChallengeCompleteTokensDisabledTotal *prometheus.CounterVec
	BotChallengeCompleteTokensFailedTotal   *prometheus.CounterVec
	BotChallengeCompleteTokensIssuedTotal   *prometheus.CounterVec
	BotChallengeCompleteTokensPassedTotal   *prometheus.CounterVec
	BotChallengeStartsTotal                 *prometheus.CounterVec
	BotChallengesFailedTotal                *prometheus.CounterVec
	BotChallengesIssuedTotal                *prometheus.CounterVec
	BotChallengesSucceededTotal             *prometheus.CounterVec
	ComputeBackendReqBodyBytesTotal         *prometheus.CounterVec
	ComputeBackendReqErrorsTotal            *prometheus.CounterVec
	ComputeBackendReqHeaderBytesTotal       *prometheus.CounterVec
	ComputeBackendReqTotal                  *prometheus.CounterVec
	ComputeBackendRespBodyBytesTotal        *prometheus.CounterVec
	ComputeBackendRespHeaderBytesTotal      *prometheus.CounterVec
	ComputeExecutionTimeTotal               *prometheus.CounterVec
	ComputeGlobalsLimitExceededTotal        *prometheus.CounterVec
	ComputeGuestErrorsTotal                 *prometheus.CounterVec
	ComputeHeapLimitExceededTotal           *prometheus.CounterVec
	ComputeRAMUsedBytesTotal                *prometheus.CounterVec
	ComputeReqBodyBytesTotal                *prometheus.CounterVec
	ComputeReqHeaderBytesTotal              *prometheus.CounterVec
	ComputeRequestTimeBilledTotal           *prometheus.CounterVec
	ComputeRequestTimeTotal                 *prometheus.CounterVec
	ComputeRequestsTotal                    *prometheus.CounterVec
	ComputeResourceLimitExceedTotal         *prometheus.CounterVec
	ComputeRespBodyBytesTotal               *prometheus.CounterVec
	ComputeRespHeaderBytesTotal             *prometheus.CounterVec
	ComputeRespStatusTotal                  *prometheus.CounterVec
	ComputeRuntimeErrorsTotal               *prometheus.CounterVec
	ComputeStackLimitExceededTotal          *prometheus.CounterVec
	DDOSActionBlackholeTotal                *prometheus.CounterVec
	DDOSActionCloseTotal                    *prometheus.CounterVec
	DDOSActionDowngradeTotal                *prometheus.CounterVec
	DDOSActionDowngradedConnectionsTotal    *prometheus.CounterVec
	DDOSActionLimitStreamsConnectionsTotal  *prometheus.CounterVec
	DDOSActionLimitStreamsRequestsTotal     *prometheus.CounterVec
	DDOSActionTarpitAcceptTotal             *prometheus.CounterVec
	DDOSActionTarpitTotal                   *prometheus.CounterVec
	DeliverSubCountTotal                    *prometheus.CounterVec
	DeliverSubTimeTotal                     *prometheus.CounterVec
	EdgeHitRequestsTotal                    *prometheus.CounterVec
	EdgeHitRespBodyBytesTotal               *prometheus.CounterVec
	EdgeHitRespHeaderBytesTotal             *prometheus.CounterVec
	EdgeMissRequestsTotal                   *prometheus.CounterVec
	EdgeMissRespBodyBytesTotal              *prometheus.CounterVec
	EdgeMissRespHeaderBytesTotal            *prometheus.CounterVec
	EdgeRespBodyBytesTotal                  *prometheus.CounterVec
	EdgeRespHeaderBytesTotal                *prometheus.CounterVec
	EdgeTotal                               *prometheus.CounterVec
	ErrorSubCountTotal                      *prometheus.CounterVec
	ErrorSubTimeTotal                       *prometheus.CounterVec
	ErrorsTotal                             *prometheus.CounterVec
	FanoutBackendReqBodyBytesTotal          *prometheus.CounterVec
	FanoutBackendReqHeaderBytesTotal        *prometheus.CounterVec
	FanoutBackendRespBodyBytesTotal         *prometheus.CounterVec
	FanoutBackendRespHeaderBytesTotal       *prometheus.CounterVec
	FanoutConnTimeMsTotal                   *prometheus.CounterVec
	FanoutRecvPublishesTotal                *prometheus.CounterVec
	FanoutReqBodyBytesTotal                 *prometheus.CounterVec
	FanoutReqHeaderBytesTotal               *prometheus.CounterVec
	FanoutRespBodyBytesTotal                *prometheus.CounterVec
	FanoutRespHeaderBytesTotal              *prometheus.CounterVec
	FanoutSendPublishesTotal                *prometheus.CounterVec
	FetchSubCountTotal                      *prometheus.CounterVec
	FetchSubTimeTotal                       *prometheus.CounterVec
	HTTPTotal                               *prometheus.CounterVec
	HTTP2Total                              *prometheus.CounterVec
	HTTP3Total                              *prometheus.CounterVec
	HashSubCountTotal                       *prometheus.CounterVec
	HashSubTimeTotal                        *prometheus.CounterVec
	HeaderSizeTotal                         *prometheus.CounterVec
	HitRespBodyBytesTotal                   *prometheus.CounterVec
	HitSubCountTotal                        *prometheus.CounterVec
	HitSubTimeTotal                         *prometheus.CounterVec
	HitsTimeTotal                           *prometheus.CounterVec
	HitsTotal                               *prometheus.CounterVec
	IPv6Total                               *prometheus.CounterVec
	ImgOptoRespBodyBytesTotal               *prometheus.CounterVec
	ImgOptoRespHeaderBytesTotal             *prometheus.CounterVec
	ImgOptoShieldRespBodyBytesTotal         *prometheus.CounterVec
	ImgOptoShieldRespHeaderBytesTotal       *prometheus.CounterVec
	ImgOptoShieldTotal                      *prometheus.CounterVec
	ImgOptoTotal                            *prometheus.CounterVec
	ImgOptoTransformRespBodyBytesTotal      *prometheus.CounterVec
	ImgOptoTransformRespHeaderBytesTotal    *prometheus.CounterVec
	ImgOptoTransformTotal                   *prometheus.CounterVec
	ImgVideoFramesTotal                     *prometheus.CounterVec
	ImgVideoRespBodyBytesTotal              *prometheus.CounterVec
	ImgVideoRespHeaderBytesTotal            *prometheus.CounterVec
	ImgVideoShieldFramesTotal               *prometheus.CounterVec
	ImgVideoShieldRespBodyBytesTotal        *prometheus.CounterVec
	ImgVideoShieldRespHeaderBytesTotal      *prometheus.CounterVec
	ImgVideoShieldTotal                     *prometheus.CounterVec
	ImgVideoTotal                           *prometheus.CounterVec
	KVStoreClassAOperationsTotal            *prometheus.CounterVec
	KVStoreClassBOperationsTotal            *prometheus.CounterVec
	LogBytesTotal                           *prometheus.CounterVec
	LoggingTotal                            *prometheus.CounterVec
	MissDurationSeconds                     *prometheus.HistogramVec
	MissRespBodyBytesTotal                  *prometheus.CounterVec
	MissSubCountTotal                       *prometheus.CounterVec
	MissSubTimeTotal                        *prometheus.CounterVec
	MissTimeTotal                           *prometheus.CounterVec
	MissesTotal                             *prometheus.CounterVec
	OTFPDeliverTimeTotal                    *prometheus.CounterVec
	OTFPManifestTotal                       *prometheus.CounterVec
	OTFPRespBodyBytesTotal                  *prometheus.CounterVec
	OTFPRespHeaderBytesTotal                *prometheus.CounterVec
	OTFPShieldRespBodyBytesTotal            *prometheus.CounterVec
	OTFPShieldRespHeaderBytesTotal          *prometheus.CounterVec
	OTFPShieldTimeTotal                     *prometheus.CounterVec
	OTFPShieldTotal                         *prometheus.CounterVec
	OTFPTotal                               *prometheus.CounterVec
	OTFPTransformRespBodyBytesTotal         *prometheus.CounterVec
	OTFPTransformRespHeaderBytesTotal       *prometheus.CounterVec
	OTFPTransformTimeTotal                  *prometheus.CounterVec
	OTFPTransformTotal                      *prometheus.CounterVec
	ObjectSizeBytes                         *prometheus.HistogramVec
	OriginCacheFetchRespBodyBytesTotal      *prometheus.CounterVec
	OriginCacheFetchRespHeaderBytesTotal    *prometheus.CounterVec
	OriginCacheFetchesTotal                 *prometheus.CounterVec
	OriginFetchBodyBytesTotal               *prometheus.CounterVec
	OriginFetchHeaderBytesTotal             *prometheus.CounterVec
	OriginFetchRespBodyBytesTotal           *prometheus.CounterVec
	OriginFetchRespHeaderBytesTotal         *prometheus.CounterVec
	OriginFetchesTotal                      *prometheus.CounterVec
	OriginRevalidationsTotal                *prometheus.CounterVec
	PCITotal                                *prometheus.CounterVec
	PassRespBodyBytesTotal                  *prometheus.CounterVec
	PassSubCountTotal                       *prometheus.CounterVec
	PassSubTimeTotal                        *prometheus.CounterVec
	PassTimeTotal                           *prometheus.CounterVec
	PassesTotal                             *prometheus.CounterVec
	Pipe                                    *prometheus.CounterVec
	PipeSubCountTotal                       *prometheus.CounterVec
	PipeSubTimeTotal                        *prometheus.CounterVec
	PredeliverSubCountTotal                 *prometheus.CounterVec
	PredeliverSubTimeTotal                  *prometheus.CounterVec
	PrehashSubCountTotal                    *prometheus.CounterVec
	PrehashSubTimeTotal                     *prometheus.CounterVec
	RealtimeAPIRequestsTotal                *prometheus.CounterVec
	RecvSubCountTotal                       *prometheus.CounterVec
	RecvSubTimeTotal                        *prometheus.CounterVec
	ReqBodyBytesTotal                       *prometheus.CounterVec
	ReqHeaderBytesTotal                     *prometheus.CounterVec
	RequestsTotal                           *prometheus.CounterVec
	RespBodyBytesTotal                      *prometheus.CounterVec
	RespHeaderBytesTotal                    *prometheus.CounterVec
	RestartTotal                            *prometheus.CounterVec
	SegBlockOriginFetchesTotal              *prometheus.CounterVec
	SegBlockShieldFetchesTotal              *prometheus.CounterVec
	ShieldCacheFetchesTotal                 *prometheus.CounterVec
	ShieldFetchBodyBytesTotal               *prometheus.CounterVec
	ShieldFetchHeaderBytesTotal             *prometheus.CounterVec
	ShieldFetchRespBodyBytesTotal           *prometheus.CounterVec
	ShieldFetchRespHeaderBytesTotal         *prometheus.CounterVec
	ShieldFetchesTotal                      *prometheus.CounterVec
	ShieldHitRequestsTotal                  *prometheus.CounterVec
	ShieldHitRespBodyBytesTotal             *prometheus.CounterVec
	ShieldHitRespHeaderBytesTotal           *prometheus.CounterVec
	ShieldMissRequestsTotal                 *prometheus.CounterVec
	ShieldMissRespBodyBytesTotal            *prometheus.CounterVec
	ShieldMissRespHeaderBytesTotal          *prometheus.CounterVec
	ShieldRespBodyBytesTotal                *prometheus.CounterVec
	ShieldRespHeaderBytesTotal              *prometheus.CounterVec
	ShieldRevalidationsTotal                *prometheus.CounterVec
	ShieldTotal                             *prometheus.CounterVec
	StatusCodeTotal                         *prometheus.CounterVec
	StatusGroupTotal                        *prometheus.CounterVec
	SynthsTotal                             *prometheus.CounterVec
	TLSTotal                                *prometheus.CounterVec
	UncacheableTotal                        *prometheus.CounterVec
	VclOnComputeEdgeHitRequestsTotal        *prometheus.CounterVec
	VclOnComputeEdgeMissRequestsTotal       *prometheus.CounterVec
	VclOnComputeErrorRequestsTotal          *prometheus.CounterVec
	VclOnComputeHitRequestsTotal            *prometheus.CounterVec
	VclOnComputeMissRequestsTotal           *prometheus.CounterVec
	VclOnComputePassRequestsTotal           *prometheus.CounterVec
	VclOnComputeSynthRequestsTotal          *prometheus.CounterVec
	VideoTotal                              *prometheus.CounterVec
	WAFBlockedTotal                         *prometheus.CounterVec
	WAFLoggedTotal                          *prometheus.CounterVec
	WAFPassedTotal                          *prometheus.CounterVec
	WebsocketBackendReqBodyBytesTotal       *prometheus.CounterVec
	WebsocketBackendReqHeaderBytesTotal     *prometheus.CounterVec
	WebsocketBackendRespBodyBytesTotal      *prometheus.CounterVec
	WebsocketBackendRespHeaderBytesTotal    *prometheus.CounterVec
	WebsocketConnTimeMsTotal                *prometheus.CounterVec
	WebsocketReqBodyBytesTotal              *prometheus.CounterVec
	WebsocketReqHeaderBytesTotal            *prometheus.CounterVec
	WebsocketRespBodyBytesTotal             *prometheus.CounterVec
	WebsocketRespHeaderBytesTotal           *prometheus.CounterVec
}

// NewMetrics returns a new set of metrics registered to the registerer.
// Only metrics whose names pass the name filter are registered.
func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {
	labels := []string{
		"service_id",
		"service_name",
		"datacenter",
	}

	m := Metrics{
		AttackBlockedReqBodyBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_blocked_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was blocked."}, labels),
		AttackBlockedReqHeaderBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_blocked_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was blocked."}, labels),
		AttackLoggedReqBodyBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_logged_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was logged."}, labels),
		AttackLoggedReqHeaderBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_logged_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was logged."}, labels),
		AttackPassedReqBodyBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_passed_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule that was passed."}, labels),
		AttackPassedReqHeaderBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_passed_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule that was passed."}, labels),
		AttackReqBodyBytesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_req_body_bytes_total", Help: "Total body bytes received from requests that triggered a WAF rule."}, labels),
		AttackReqHeaderBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_req_header_bytes_total", Help: "Total header bytes received from requests that triggered a WAF rule."}, labels),
		AttackRespSynthBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "attack_resp_synth_bytes_total", Help: "Total bytes delivered for requests that triggered a WAF rule and returned a synthetic response."}, labels),
		BackendReqBodyBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_body_bytes_total", Help: "Total body bytes sent to origin."}, labels),
		BackendReqHeaderBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_header_bytes_total", Help: "Total header bytes sent to origin."}, labels),
		BlacklistedTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "blacklist_total", Help: "TODO"}, labels),
		BodySizeTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "body_size_total", Help: "Total body bytes delivered (alias for resp_body_bytes)."}, labels),
		BotChallengeCompleteTokensCheckedTotal:  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_complete_tokens_checked_total", Help: "The number of challenge-complete tokens checked."}, labels),
		BotChallengeCompleteTokensDisabledTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_complete_tokens_disabled_total", Help: "TThe number of challenge-complete tokens not checked because the feature was disabled."}, labels),
		BotChallengeCompleteTokensFailedTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_complete_tokens_failed_total", Help: "TThe number of challenge-complete tokens that failed validation."}, labels),
		BotChallengeCompleteTokensIssuedTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_complete_tokens_issued_total", Help: "The number of challenge-complete tokens issued. For example, issuing a challenge-complete token after a series of CAPTCHA challenges ending in success."}, labels),
		BotChallengeCompleteTokensPassedTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_complete_tokens_passed_total", Help: "The number of challenge-complete tokens that passed validation."}, labels),
		BotChallengeStartsTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenge_starts_total", Help: "The number of challenge-start tokens created."}, labels),
		BotChallengesFailedTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenges_failed_total", Help: "The number of failed challenge solutions processed. For example, an incorrect CAPTCHA solution."}, labels),
		BotChallengesIssuedTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenges_issued_total", Help: "The number of challenges issued. For example, the issuance of a CAPTCHA challenge."}, labels),
		BotChallengesSucceededTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bot_challenges_succeeded_total", Help: "The number of successful challenge solutions processed. For example, a correct CAPTCHA solution."}, labels),
		ComputeBackendReqBodyBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_body_bytes_total", Help: "Total body bytes sent to backends (origins) by Compute@Edge."}, labels),
		ComputeBackendReqErrorsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_errors_total", Help: "Number of backend request errors, including timeouts."}, labels),
		ComputeBackendReqHeaderBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_header_bytes_total", Help: "Total header bytes sent to backends (origins) by Compute@Edge."}, labels),
		ComputeBackendReqTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_bereq_total", Help: "Number of backend requests started."}, labels),
		ComputeBackendRespBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_beresp_body_bytes_total", Help: "Total body bytes received from backends (origins) by Compute@Edge."}, labels),
		ComputeBackendRespHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_beresp_header_bytes_total", Help: "Total header bytes received from backends (origins) by Compute@Edge."}, labels),
		ComputeExecutionTimeTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_execution_time_total", Help: "The amount of active CPU time used to process your requests (in seconds)."}, labels),
		ComputeGlobalsLimitExceededTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_globals_limit_exceeded_total", Help: "Number of times a guest exceeded its globals limit."}, labels),
		ComputeGuestErrorsTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_guest_errors_total", Help: "Number of times a service experienced a guest code error."}, labels),
		ComputeHeapLimitExceededTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_heap_limit_exceeded_total", Help: "Number of times a guest exceeded its heap limit."}, labels),
		ComputeRAMUsedBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_ram_used_bytes_total", Help: "The amount of RAM used for your site by Fastly."}, labels),
		ComputeReqBodyBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_req_body_bytes_total", Help: "Total body bytes received by Compute@Edge."}, labels),
		ComputeReqHeaderBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_req_header_bytes_total", Help: "Total header bytes received by Compute@Edge."}, labels),
		ComputeRequestTimeBilledTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_request_time_billed_total", Help: "The total amount of request processing time you will be billed for, measured in 50 millisecond increments. (in seconds)"}, labels),
		ComputeRequestTimeTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_request_time_total", Help: "The total amount of time used to process your requests, including active CPU time (in seconds)."}, labels),
		ComputeRequestsTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_requests_total", Help: "The total number of requests that were received for your site by Fastly."}, labels),
		ComputeResourceLimitExceedTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resource_limit_exceeded_total", Help: "Number of times a guest exceeded its resource limit, includes heap, stack, globals, and code execution timeout."}, labels),
		ComputeRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_body_bytes_total", Help: "Total body bytes sent from Compute@Edge to end user."}, labels),
		ComputeRespHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_header_bytes_total", Help: "Total header bytes sent from Compute@Edge to end user."}, labels),
		ComputeRespStatusTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_resp_status_total", Help: "Number of responses delivered delivered by Compute@Edge, by status code group."}, addLabels(labels, "status_group")),
		ComputeRuntimeErrorsTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_runtime_errors_total", Help: "Number of times a service experienced a guest runtime error."}, labels),
		ComputeStackLimitExceededTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "compute_stack_limit_exceeded_total", Help: "Number of times a guest exceeded its stack limit."}, labels),
		DDOSActionBlackholeTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_blackhole_total", Help: "The number of times the blackhole action was taken. The blackhole action quietly closes a TCP connection without sending a reset. The blackhole action quietly closes a TCP connection without notifying its peer (all TCP state is dropped)."}, labels),
		DDOSActionCloseTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_close_total", Help: "The number of times the close action was taken. The close action aborts the connection as soon as possible. The close action takes effect either right after accept, right after the client hello, or right after the response was sent."}, labels),
		DDOSActionDowngradeTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_downgrade_total", Help: "The number of times the downgrade action was taken. The downgrade action restricts the client to http1."}, labels),
		DDOSActionDowngradedConnectionsTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_downgraded_connections_total", Help: "The number of connections the downgrade action was applied to. The downgrade action restricts the connection to http1."}, labels),
		DDOSActionLimitStreamsConnectionsTotal:  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_limit_streams_connections_total", Help: "For HTTP/2, the number of connections the limit-streams action was applied to. The limit-streams action caps the allowed number of concurrent streams in a connection."}, labels),
		DDOSActionLimitStreamsRequestsTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_limit_streams_requests_total", Help: "For HTTP/2, the number of requests made on a connection for which the limit-streams action was taken. The limit-streams action caps the allowed number of concurrent streams in a connection."}, labels),
		DDOSActionTarpitAcceptTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_tarpit_accept_total", Help: "The number of times the tarpit-accept action was taken. The tarpit-accept action adds a delay when accepting future connections."}, labels),
		DDOSActionTarpitTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ddos_action_tarpit_total", Help: "The number of times the tarpit action was taken. The tarpit action delays writing the response to the client."}, labels),
		DeliverSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "deliver_sub_count_total", Help: "Number of executions of the 'deliver' Varnish subroutine."}, labels),
		DeliverSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "deliver_sub_time_total", Help: "Time spent inside the 'deliver' Varnish subroutine (in seconds)."}, labels),
		EdgeHitRequestsTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_hit_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a hit at the edge."}, labels),
		EdgeHitRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_hit_resp_body_bytes_total", Help: "Body bytes delivered for edge hits."}, labels),
		EdgeHitRespHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_hit_resp_header_bytes_total", Help: "Header bytes delivered for edge hits."}, labels),
		EdgeMissRequestsTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_miss_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a miss at the edge."}, labels),
		EdgeMissRespBodyBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_miss_resp_body_bytes_total", Help: "Body bytes delivered for edge misses."}, labels),
		EdgeMissRespHeaderBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_miss_resp_header_bytes_total", Help: "Header bytes delivered for edge misses."}, labels),
		EdgeRespBodyBytesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_body_bytes_total", Help: "Total body bytes delivered from Fastly to the end user."}, labels),
		EdgeRespHeaderBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_header_bytes_total", Help: "Total header bytes delivered from Fastly to the end user."}, labels),
		EdgeTotal:                               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_total", Help: "Number of requests sent by end users to Fastly."}, labels),
		ErrorSubCountTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "error_sub_count_total", Help: "Number of executions of the 'error' Varnish subroutine."}, labels),
		ErrorSubTimeTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "error_sub_time_total", Help: "Time spent inside the 'error' Varnish subroutine (in seconds)."}, labels),
		ErrorsTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "errors_total", Help: "Number of cache errors."}, labels),
		FanoutBackendReqBodyBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_bereq_body_bytes_total", Help: "Total body or message content bytes sent to backends over Fanout connections."}, labels),
		FanoutBackendReqHeaderBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_bereq_header_bytes_total", Help: "Total header bytes sent to backends over Fanout connections."}, labels),
		FanoutBackendRespBodyBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_beresp_body_bytes_total", Help: "Total body or message content bytes received from backends over Fanout connections."}, labels),
		FanoutBackendRespHeaderBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_beresp_header_bytes_total", Help: "Total header bytes received from backends over Fanout connections."}, labels),
		FanoutConnTimeMsTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_conn_time_ms_total", Help: "Total duration of Fanout connections with end users."}, labels),
		FanoutRecvPublishesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_recv_publishes_total", Help: "Total published messages received from the publish API endpoint."}, labels),
		FanoutReqBodyBytesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_req_body_bytes_total", Help: "Total body or message content bytes received from end users over Fanout connections."}, labels),
		FanoutReqHeaderBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_req_header_bytes_total", Help: "Total header bytes received from end users over Fanout connections."}, labels),
		FanoutRespBodyBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_resp_body_bytes_total", Help: "Total body or message content bytes sent to end users over Fanout connections, excluding published message content."}, labels),
		FanoutRespHeaderBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_resp_header_bytes_total", Help: "Total header bytes sent to end users over Fanout connections."}, labels),
		FanoutSendPublishesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fanout_send_publishes_total", Help: "Total published messages sent to end users."}, labels),
		FetchSubCountTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fetch_sub_count_total", Help: "Number of executions of the 'fetch' Varnish subroutine."}, labels),
		FetchSubTimeTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "fetch_sub_time_total", Help: "Time spent inside the 'fetch' Varnish subroutine (in seconds)."}, labels),
		HTTPTotal:                               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "http_total", Help: "Number of requests received, by HTTP version."}, addLabels(labels, "http_version")),
		HTTP2Total:                              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "http2_total", Help: "Number of requests received over HTTP2."}, labels),
		HTTP3Total:                              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "http3_total", Help: "Number of requests received over HTTP3."}, labels),
		HashSubCountTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hash_sub_count_total", Help: "Number of executions of the 'hash' Varnish subroutine."}, labels),
		HashSubTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hash_sub_time_total", Help: "Time spent inside the 'hash' Varnish subroutine (in seconds)."}, labels),
		HeaderSizeTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "header_size_total", Help: "Total header bytes delivered (alias for resp_header_bytes)."}, labels),
		HitRespBodyBytesTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_resp_body_bytes_total", Help: "Total body bytes delivered for cache hits."}, labels),
		HitSubCountTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_sub_count_total", Help: "Number of executions of the 'hit' Varnish subroutine."}, labels),
		HitSubTimeTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hit_sub_time_total", Help: "Time spent inside the 'hit' Varnish subroutine (in seconds)."}, labels),
		HitsTimeTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hits_time_total", Help: "Total amount of time spent processing cache hits (in seconds)."}, labels),
		HitsTotal:                               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "hits_total", Help: "Number of cache hits."}, labels),
		IPv6Total:                               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "ipv6_total", Help: "Number of requests that were received over IPv6."}, labels),
		ImgOptoRespBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_resp_body_bytes_total", Help: "Total body bytes delivered from the Fastly Image Optimizer service."}, labels),
		ImgOptoRespHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_resp_header_bytes_total", Help: "Total header bytes delivered from the Fastly Image Optimizer service."}, labels),
		ImgOptoShieldRespBodyBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgOptoShieldRespHeaderBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgOptoShieldTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_shield_total", Help: "Number of responses delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgOptoTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_total", Help: "Number of responses that came from the Fastly Image Optimizer service."}, labels),
		ImgOptoTransformRespBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transform_resp_body_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly Image Optimizer service."}, labels),
		ImgOptoTransformRespHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transform_resp_header_bytes_total", Help: "Total header bytes of transforms delivered from the Fastly Image Optimizer service."}, labels),
		ImgOptoTransformTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgopto_transforms_total", Help: "Total transforms performed by the Fastly Image Optimizer service."}, labels),
		ImgVideoFramesTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_frames_total", Help: "Number of video frames that came from the Fastly Image Optimizer service."}, labels),
		ImgVideoRespBodyBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_resp_body_bytes_total", Help: "Total body bytes of video delivered from the Fastly Image Optimizer service."}, labels),
		ImgVideoRespHeaderBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_resp_header_bytes_total", Help: "Total header bytes of video delivered from the Fastly Image Optimizer service."}, labels),
		ImgVideoShieldFramesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_frames_total", Help: "Number of video frames delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgVideoShieldRespBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_resp_body_bytes_total", Help: "Total body bytes of video delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgVideoShieldRespHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_resp_header_bytes_total", Help: "Total header bytes of video delivered via a shield from the Fastly Image Optimizer service."}, labels),
		ImgVideoShieldTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_shield_total", Help: "Number of video responses that came via a shield from the Fastly Image Optimizer service."}, labels),
		ImgVideoTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "imgvideo_total", Help: "Number of video responses that came via a shield from the Fastly Image Optimizer service."}, labels),
		KVStoreClassAOperationsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "kv_store_class_a_operations_total", Help: "The total number of class a operations for the KV store."}, labels),
		KVStoreClassBOperationsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "kv_store_class_b_operations_total", Help: "The total number of class b operations for the KV store."}, labels),
		LogBytesTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "log_bytes_total", Help: "Total log bytes sent."}, labels),
		LoggingTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "logging_total", Help: "Number of log lines sent."}, labels),
		MissDurationSeconds:                     prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_duration_seconds", Help: "Histogram of time spent processing cache misses (in seconds).", Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 4, 8, 16, 32, 60}}, labels),
		MissRespBodyBytesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_resp_body_bytes_total", Help: "Total body bytes delivered for cache misses."}, labels),
		MissSubCountTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_sub_count_total", Help: "Number of executions of the 'miss' Varnish subroutine."}, labels),
		MissSubTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_sub_time_total", Help: "Time spent inside the 'miss' Varnish subroutine (in seconds)."}, labels),
		MissTimeTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_time_total", Help: "Total amount of time spent processing cache misses (in seconds)."}, labels),
		MissesTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "miss_total", Help: "Number of cache misses."}, labels),
		OTFPDeliverTimeTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_total", Help: "Number of responses that came from the Fastly On-the-Fly Packager."}, labels),
		OTFPManifestTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_deliver_time_total", Help: "Total amount of time spent delivering a response from the Fastly On-the-Fly Packager (in seconds)."}, labels),
		OTFPRespBodyBytesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_manifests_total", Help: "Number of responses that were manifest files from the Fastly On-the-Fly Packager."}, labels),
		OTFPRespHeaderBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_resp_body_bytes_total", Help: "Total body bytes delivered from the Fastly On-the-Fly Packager."}, labels),
		OTFPShieldRespBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_resp_header_bytes_total", Help: "Total header bytes delivered from the Fastly On-the-Fly Packager."}, labels),
		OTFPShieldRespHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_total", Help: "Number of responses delivered from the Fastly On-the-Fly Packager"}, labels),
		OTFPShieldTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield for the Fastly On-the-Fly Packager."}, labels),
		OTFPShieldTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield for the Fastly On-the-Fly Packager."}, labels),
		OTFPTotal:                               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_shield_time_total", Help: "Total amount of time spent delivering a response via a shield from the Fastly On-the-Fly Packager (in seconds)."}, labels),
		OTFPTransformRespBodyBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transforms_total", Help: "Number of transforms performed by the Fastly On-the-Fly Packager."}, labels),
		OTFPTransformRespHeaderBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_resp_body_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly On-the-Fly Packager."}, labels),
		OTFPTransformTimeTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_resp_header_bytes_total", Help: "Total body bytes of transforms delivered from the Fastly On-the-Fly Packager."}, labels),
		OTFPTransformTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "otfp_transform_time_total", Help: "Total amount of time spent performing transforms from the Fastly On-the-Fly Packager."}, labels),
		ObjectSizeBytes:                         prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: namespace, Subsystem: subsystem, Name: "object_size_bytes", Help: "Histogram of count of objects served, bucketed by object size range.", Buckets: []float64{1024, 10240, 102400, 1.024e+06, 1.024e+07, 1.024e+08, 1.024e+09}}, labels),
		OriginCacheFetchRespBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_cache_fetch_resp_body_bytes_total", Help: "Body bytes received from origin for cacheable content."}, labels),
		OriginCacheFetchRespHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_cache_fetch_resp_header_bytes_total", Help: "Header bytes received from an origin for cacheable content."}, labels),
		OriginCacheFetchesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_cache_fetches_total", Help: "The total number of completed requests made to backends (origins) that returned cacheable content."}, labels),
		OriginFetchBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_body_bytes_total", Help: "Total request body bytes sent to origin."}, labels),
		OriginFetchHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_header_bytes_total", Help: "Total request header bytes sent to origin."}, labels),
		OriginFetchRespBodyBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_body_bytes_total", Help: "Total body bytes received from origin."}, labels),
		OriginFetchRespHeaderBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_header_bytes_total", Help: "Total header bytes received from origin."}, labels),
		OriginFetchesTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetches_total", Help: "Number of requests sent to origin."}, labels),
		OriginRevalidationsTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_revalidations_total", Help: "Number of responses received from origin with a 304 status code in response to an If-Modified-Since or If-None-Match request. Under regular scenarios, a revalidation will imply a cache hit. However, if using Fastly Image Optimizer or segmented caching this may result in a cache miss."}, labels),
		PCITotal:                                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pci_total", Help: "Number of responses with the PCI flag turned on."}, labels),
		PassRespBodyBytesTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_resp_body_bytes_total", Help: "Total body bytes delivered for cache passes."}, labels),
		PassSubCountTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_sub_count_total", Help: "Number of executions of the 'pass' Varnish subroutine."}, labels),
		PassSubTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_sub_time_total", Help: "Time spent inside the 'pass' Varnish subroutine (in seconds)."}, labels),
		PassTimeTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_time_total", Help: "Total amount of time spent processing cache passes (in seconds)."}, labels),
		PassesTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pass_total", Help: "Number of requests that passed through the CDN without being cached."}, labels),
		Pipe:                                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe", Help: "Pipe operations performed."}, labels),
		PipeSubCountTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe_sub_count_total", Help: "Number of executions of the 'pipe' Varnish subroutine."}, labels),
		PipeSubTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "pipe_sub_time_total", Help: "Time spent inside the 'pipe' Varnish subroutine (in seconds)."}, labels),
		PredeliverSubCountTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "predeliver_sub_count_total", Help: "Number of executions of the 'predeliver' Varnish subroutine."}, labels),
		PredeliverSubTimeTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "predeliver_sub_time_total", Help: "Time spent inside the 'predeliver' Varnish subroutine (in seconds)."}, labels),
		PrehashSubCountTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "prehash_sub_count_total", Help: "Number of executions of the 'prehash' Varnish subroutine."}, labels),
		PrehashSubTimeTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "prehash_sub_time_total", Help: "Time spent inside the 'prehash' Varnish subroutine (in seconds)."}, labels),
		RealtimeAPIRequestsTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "realtime_api_requests_total", Help: "Total requests made to the real-time stats API."}, []string{"service_id", "service_name", "result"}),
		RecvSubCountTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "recv_sub_count_total", Help: "Number of executions of the 'recv' Varnish subroutine."}, labels),
		RecvSubTimeTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "recv_sub_time_total", Help: "Time spent inside the 'recv' Varnish subroutine (in seconds)."}, labels),
		ReqBodyBytesTotal:                       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "req_body_bytes_total", Help: "Total body bytes received."}, labels),
		ReqHeaderBytesTotal:                     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "req_header_bytes_total", Help: "Total header bytes received."}, labels),
		RequestsTotal:                           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "requests_total", Help: "Number of requests processed."}, labels),
		RespBodyBytesTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_body_bytes_total", Help: "Total body bytes delivered."}, labels),
		RespHeaderBytesTotal:                    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_header_bytes_total", Help: "Total header bytes delivered."}, labels),
		RestartTotal:                            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "restarts_total", Help: "Number of restarts performed."}, labels),
		SegBlockOriginFetchesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "segblock_origin_fetches_total", Help: "Number of Range requests to origin for segments of resources when using segmented caching."}, labels),
		SegBlockShieldFetchesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "segblock_shield_fetches_total", Help: "Number of Range requests to a shield for segments of resources when using segmented caching."}, labels),
		ShieldCacheFetchesTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_cache_fetches_total", Help: "The total number of completed requests made to shields that returned cacheable content."}, labels),
		ShieldFetchBodyBytesTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_body_bytes_total", Help: "Total request body bytes sent to a shield."}, labels),
		ShieldFetchHeaderBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_header_bytes_total", Help: "Total request header bytes sent to a shield."}, labels),
		ShieldFetchRespBodyBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_resp_body_bytes_total", Help: "Total response body bytes sent from a shield to the edge."}, labels),
		ShieldFetchRespHeaderBytesTotal:         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetch_resp_header_bytes_total", Help: "Total response header bytes sent from a shield to the edge."}, labels),
		ShieldFetchesTotal:                      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_fetches_total", Help: "Number of requests made from one Fastly data center to another, as part of shielding."}, labels),
		ShieldHitRequestsTotal:                  prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_hit_requests_total", Help: "Number of requests that resulted in a hit at a shield."}, labels),
		ShieldHitRespBodyBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_hit_resp_body_bytes_total", Help: "Body bytes delivered for shield hits."}, labels),
		ShieldHitRespHeaderBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_hit_resp_header_bytes_total", Help: "Header bytes delivered for shield hits."}, labels),
		ShieldMissRequestsTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_miss_requests_total", Help: "Number of requests that resulted in a miss at a shield."}, labels),
		ShieldMissRespBodyBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_miss_resp_body_bytes_total", Help: "Body bytes delivered for shield misses."}, labels),
		ShieldMissRespHeaderBytesTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_miss_resp_header_bytes_total", Help: "Header bytes delivered for shield misses."}, labels),
		ShieldRespBodyBytesTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_resp_body_bytes_total", Help: "Total body bytes delivered via a shield."}, labels),
		ShieldRespHeaderBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_resp_header_bytes_total", Help: "Total header bytes delivered via a shield."}, labels),
		ShieldRevalidationsTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_revalidations_total", Help: "Number of responses received from origin with a 304 status code, in response to an If-Modified-Since or If-None-Match request to a shield. Under regular scenarios, a revalidation will imply a cache hit. However, if using segmented caching this may result in a cache miss."}, labels),
		ShieldTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "shield_total", Help: "Number of requests from edge to the shield POP."}, labels),
		StatusCodeTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_code_total", Help: "Number of responses sent with status code 500 (Internal Server Error)."}, addLabels(labels, "status_code")),
		StatusGroupTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_group_total", Help: "Number of 'Client Error' category status codes delivered."}, addLabels(labels, "status_group")),
		SynthsTotal:                             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "synth_total", Help: "TODO"}, labels),
		TLSTotal:                                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "tls_total", Help: "Number of requests that were received over TLS."}, addLabels(labels, "tls_version")),
		UncacheableTotal:                        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "uncacheable_total", Help: "Number of requests that were designated uncachable."}, labels),
		VclOnComputeHitRequestsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_hit_requests_total", Help: "Number of cache hits for a VCL service running on Compute."}, labels),
		VclOnComputeMissRequestsTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_miss_requests_total", Help: "Number of cache misses for a VCL service running on Compute."}, labels),
		VclOnComputePassRequestsTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_pass_requests_total", Help: "Number of requests that passed through the CDN without being cached for a VCL service running on Compute."}, labels),
		VclOnComputeErrorRequestsTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_error_requests_total", Help: "Number of cache errors for a VCL service running on Compute."}, labels),
		VclOnComputeSynthRequestsTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_synth_requests_total", Help: "Number of requests that returned a synthetic response (i.e., response objects created with the synthetic VCL statement) for a VCL service running on Compute."}, labels),
		VclOnComputeEdgeHitRequestsTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_edge_hit_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a hit at the edge for a VCL service running on Compute."}, labels),
		VclOnComputeEdgeMissRequestsTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "vcl_on_compute_edge_miss_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a miss at the edge for a VCL service running on Compute."}, labels),
		VideoTotal:                              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "video_total", Help: "Number of responses with the video segment or video manifest MIME type (i.e., application/x-mpegurl, application/vnd.apple.mpegurl, application/f4m, application/dash+xml, application/vnd.ms-sstr+xml, ideo/mp2t, audio/aac, video/f4f, video/x-flv, video/mp4, audio/mp4)."}, labels),
		WAFBlockedTotal:                         prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_blocked_total", Help: "Number of requests that triggered a WAF rule and were blocked."}, labels),
		WAFLoggedTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_logged_total", Help: "Number of requests that triggered a WAF rule and were logged."}, labels),
		WAFPassedTotal:                          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "waf_passed_total", Help: "Number of requests that triggered a WAF rule and were passed."}, labels),
		WebsocketBackendReqBodyBytesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_bereq_body_bytes_total", Help: "Total message content bytes sent to backends over passthrough WebSocket connections."}, labels),
		WebsocketBackendReqHeaderBytesTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_bereq_header_bytes_total", Help: "Total header bytes sent to backends over passthrough WebSocket connections."}, labels),
		WebsocketBackendRespBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_beresp_body_bytes_total", Help: "Total message content bytes received from backends over passthrough WebSocket connections."}, labels),
		WebsocketBackendRespHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_beresp_header_bytes_total", Help: "Total header bytes received from backends over passthrough WebSocket connections."}, labels),
		WebsocketConnTimeMsTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_conn_time_ms_total", Help: "Total duration of passthrough WebSocket connections with end users."}, labels),
		WebsocketReqBodyBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_req_body_bytes_total", Help: "Total message content bytes received from end users over passthrough WebSocket connections."}, labels),
		WebsocketReqHeaderBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_req_header_bytes_total", Help: "Total header bytes received from end users over passthrough WebSocket connections."}, labels),
		WebsocketRespBodyBytesTotal:             prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_resp_body_bytes_total", Help: "Total message content bytes sent to end users over passthrough WebSocket connections."}, labels),
		WebsocketRespHeaderBytesTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "websocket_resp_header_bytes_total", Help: "Total header bytes sent to end users over passthrough WebSocket connections."}, labels),
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

func addLabels(ls []string, l ...string) []string {
	labels := make([]string, len(ls), len(ls)+len(l))

	copy(labels, ls)

	return append(labels, l...)
}
