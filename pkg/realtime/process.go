package realtime

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// Process updates the metrics with data from the API response.
func Process(response *Response, serviceID, serviceName, serviceVersion string, m *Metrics) {
	sl := prometheus.Labels{
		"service_id":   serviceID,
		"service_name": serviceName,
	}

	for _, d := range response.Data {
		for datacenter, stats := range d.Datacenter {
			dcl := prometheus.Labels{
				"datacenter": datacenter,
			}

			labels := mergeLabels(sl, dcl)

			process(labels, stats, m)
		}
	}
}

func process(labels prometheus.Labels, stats Datacenter, m *Metrics) {
	m.AttackBlockedReqBodyBytesTotal.With(labels).Add(float64(stats.AttackBlockedReqBodyBytes))
	m.AttackBlockedReqHeaderBytesTotal.With(labels).Add(float64(stats.AttackBlockedReqHeaderBytes))
	m.AttackLoggedReqBodyBytesTotal.With(labels).Add(float64(stats.AttackLoggedReqBodyBytes))
	m.AttackLoggedReqHeaderBytesTotal.With(labels).Add(float64(stats.AttackLoggedReqHeaderBytes))
	m.AttackPassedReqBodyBytesTotal.With(labels).Add(float64(stats.AttackPassedReqBodyBytes))
	m.AttackPassedReqHeaderBytesTotal.With(labels).Add(float64(stats.AttackPassedReqHeaderBytes))
	m.AttackReqBodyBytesTotal.With(labels).Add(float64(stats.AttackReqBodyBytes))
	m.AttackReqHeaderBytesTotal.With(labels).Add(float64(stats.AttackReqHeaderBytes))
	m.AttackRespSynthBytesTotal.With(labels).Add(float64(stats.AttackRespSynthBytes))
	m.BackendReqBodyBytesTotal.With(labels).Add(float64(stats.BackendReqBodyBytes))
	m.BackendReqHeaderBytesTotal.With(labels).Add(float64(stats.BackendReqHeaderBytes))
	m.BlacklistedTotal.With(labels).Add(float64(stats.Blacklisted))
	m.BodySizeTotal.With(labels).Add(float64(stats.BodySize))
	m.BotChallengeCompleteTokensCheckedTotal.With(labels).Add(float64(stats.BotChallengeCompleteTokensChecked))
	m.BotChallengeCompleteTokensDisabledTotal.With(labels).Add(float64(stats.BotChallengeCompleteTokensDisabled))
	m.BotChallengeCompleteTokensFailedTotal.With(labels).Add(float64(stats.BotChallengeCompleteTokensFailed))
	m.BotChallengeCompleteTokensIssuedTotal.With(labels).Add(float64(stats.BotChallengeCompleteTokensIssued))
	m.BotChallengeCompleteTokensPassedTotal.With(labels).Add(float64(stats.BotChallengeCompleteTokensPassed))
	m.BotChallengeStartsTotal.With(labels).Add(float64(stats.BotChallengeStarts))
	m.BotChallengesFailedTotal.With(labels).Add(float64(stats.BotChallengesFailed))
	m.BotChallengesIssuedTotal.With(labels).Add(float64(stats.BotChallengesIssued))
	m.BotChallengesSucceededTotal.With(labels).Add(float64(stats.BotChallengesSucceeded))
	m.ComputeBackendReqBodyBytesTotal.With(labels).Add(float64(stats.ComputeBackendReqBodyBytesTotal))
	m.ComputeBackendReqErrorsTotal.With(labels).Add(float64(stats.ComputeBackendReqErrorsTotal))
	m.ComputeBackendReqHeaderBytesTotal.With(labels).Add(float64(stats.ComputeBackendReqHeaderBytesTotal))
	m.ComputeBackendReqTotal.With(labels).Add(float64(stats.ComputeBackendReqTotal))
	m.ComputeBackendRespBodyBytesTotal.With(labels).Add(float64(stats.ComputeBackendRespBodyBytesTotal))
	m.ComputeBackendRespHeaderBytesTotal.With(labels).Add(float64(stats.ComputeBackendRespHeaderBytesTotal))
	m.ComputeExecutionTimeTotal.With(labels).Add(float64(stats.ComputeExecutionTimeMilliseconds) / 10000.0)
	m.ComputeGlobalsLimitExceededTotal.With(labels).Add(float64(stats.ComputeGlobalsLimitExceededTotal))
	m.ComputeGuestErrorsTotal.With(labels).Add(float64(stats.ComputeGuestErrorsTotal))
	m.ComputeHeapLimitExceededTotal.With(labels).Add(float64(stats.ComputeHeapLimitExceededTotal))
	m.ComputeRAMUsedBytesTotal.With(labels).Add(float64(stats.ComputeRAMUsed))
	m.ComputeReqBodyBytesTotal.With(labels).Add(float64(stats.ComputeReqBodyBytesTotal))
	m.ComputeReqHeaderBytesTotal.With(labels).Add(float64(stats.ComputeReqHeaderBytesTotal))
	m.ComputeRequestTimeBilledTotal.With(labels).Add(float64(stats.ComputeRequestTimeBilledMilliseconds) / 10000.0)
	m.ComputeRequestTimeTotal.With(labels).Add(float64(stats.ComputeRequestTimeMilliseconds) / 10000.0)
	m.ComputeRequestsTotal.With(labels).Add(float64(stats.ComputeRequests))
	m.ComputeResourceLimitExceedTotal.With(labels).Add(float64(stats.ComputeResourceLimitExceedTotal))
	m.ComputeRespBodyBytesTotal.With(labels).Add(float64(stats.ComputeRespBodyBytesTotal))
	m.ComputeRespHeaderBytesTotal.With(labels).Add(float64(stats.ComputeRespHeaderBytesTotal))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "1xx"})).Add(float64(stats.ComputeRespStatus1xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "2xx"})).Add(float64(stats.ComputeRespStatus2xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "3xx"})).Add(float64(stats.ComputeRespStatus3xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "4xx"})).Add(float64(stats.ComputeRespStatus4xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "5xx"})).Add(float64(stats.ComputeRespStatus5xx))
	m.ComputeRuntimeErrorsTotal.With(labels).Add(float64(stats.ComputeRuntimeErrorsTotal))
	m.ComputeStackLimitExceededTotal.With(labels).Add(float64(stats.ComputeStackLimitExceededTotal))
	m.DDOSActionBlackholeTotal.With(labels).Add(float64(stats.DDOSActionBlackhole))
	m.DDOSActionCloseTotal.With(labels).Add(float64(stats.DDOSActionClose))
	m.DDOSActionDowngradeTotal.With(labels).Add(float64(stats.DDOSActionDowngrade))
	m.DDOSActionDowngradedConnectionsTotal.With(labels).Add(float64(stats.DDOSActionDowngradedConnections))
	m.DDOSActionLimitStreamsConnectionsTotal.With(labels).Add(float64(stats.DDOSActionLimitStreamsConnections))
	m.DDOSActionLimitStreamsRequestsTotal.With(labels).Add(float64(stats.DDOSActionLimitStreamsRequests))
	m.DDOSActionTarpitAcceptTotal.With(labels).Add(float64(stats.DDOSActionTarpitAccept))
	m.DDOSActionTarpitTotal.With(labels).Add(float64(stats.DDOSActionTarpit))
	m.DeliverSubCountTotal.With(labels).Add(float64(stats.DeliverSubCount))
	m.DeliverSubTimeTotal.With(labels).Add(float64(stats.DeliverSubTime))
	m.EdgeHitRequestsTotal.With(labels).Add(float64(stats.EdgeHitRequests))
	m.EdgeHitRespBodyBytesTotal.With(labels).Add(float64(stats.EdgeHitRespBodyBytes))
	m.EdgeHitRespHeaderBytesTotal.With(labels).Add(float64(stats.EdgeHitRespHeaderBytes))
	m.EdgeMissRequestsTotal.With(labels).Add(float64(stats.EdgeMissRequests))
	m.EdgeMissRespBodyBytesTotal.With(labels).Add(float64(stats.EdgeMissRespBodyBytes))
	m.EdgeMissRespHeaderBytesTotal.With(labels).Add(float64(stats.EdgeMissRespHeaderBytes))
	m.EdgeRespBodyBytesTotal.With(labels).Add(float64(stats.EdgeRespBodyBytes))
	m.EdgeRespHeaderBytesTotal.With(labels).Add(float64(stats.EdgeRespHeaderBytes))
	m.EdgeTotal.With(labels).Add(float64(stats.Edge))
	m.ErrorSubCountTotal.With(labels).Add(float64(stats.ErrorSubCount))
	m.ErrorSubTimeTotal.With(labels).Add(float64(stats.ErrorSubTime))
	m.ErrorsTotal.With(labels).Add(float64(stats.Errors))
	m.FanoutBackendReqBodyBytesTotal.With(labels).Add(float64(stats.FanoutBackendReqBodyBytes))
	m.FanoutBackendReqHeaderBytesTotal.With(labels).Add(float64(stats.FanoutBackendReqHeaderBytes))
	m.FanoutBackendRespBodyBytesTotal.With(labels).Add(float64(stats.FanoutBackendRespBodyBytes))
	m.FanoutBackendRespHeaderBytesTotal.With(labels).Add(float64(stats.FanoutBackendRespHeaderBytes))
	m.FanoutConnTimeMsTotal.With(labels).Add(float64(stats.FanoutConnTimeMs))
	m.FanoutRecvPublishesTotal.With(labels).Add(float64(stats.FanoutRecvPublishes))
	m.FanoutReqBodyBytesTotal.With(labels).Add(float64(stats.FanoutReqBodyBytes))
	m.FanoutReqHeaderBytesTotal.With(labels).Add(float64(stats.FanoutReqHeaderBytes))
	m.FanoutRespBodyBytesTotal.With(labels).Add(float64(stats.FanoutRespBodyBytes))
	m.FanoutRespHeaderBytesTotal.With(labels).Add(float64(stats.FanoutRespHeaderBytes))
	m.FanoutSendPublishesTotal.With(labels).Add(float64(stats.FanoutSendPublishes))
	m.FetchSubCountTotal.With(labels).Add(float64(stats.FetchSubCount))
	m.FetchSubTimeTotal.With(labels).Add(float64(stats.FetchSubTime))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "1"})).Add(float64(stats.Requests - (stats.HTTP2 + stats.HTTP3)))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "2"})).Add(float64(stats.HTTP2))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "3"})).Add(float64(stats.HTTP3))
	m.HTTP2Total.With(labels).Add(float64(stats.HTTP2))
	m.HTTP3Total.With(labels).Add(float64(stats.HTTP3))
	m.HashSubCountTotal.With(labels).Add(float64(stats.HashSubCount))
	m.HashSubTimeTotal.With(labels).Add(float64(stats.HashSubTime))
	m.HeaderSizeTotal.With(labels).Add(float64(stats.HeaderSize))
	m.HitRespBodyBytesTotal.With(labels).Add(float64(stats.HitRespBodyBytes))
	m.HitSubCountTotal.With(labels).Add(float64(stats.HitSubCount))
	m.HitSubTimeTotal.With(labels).Add(float64(stats.HitSubTime))
	m.HitsTimeTotal.With(labels).Add(float64(stats.HitsTime))
	m.HitsTotal.With(labels).Add(float64(stats.Hits))
	m.IPv6Total.With(labels).Add(float64(stats.IPv6))
	m.ImgOptoRespBodyBytesTotal.With(labels).Add(float64(stats.ImgOptoRespBodyBytes))
	m.ImgOptoRespHeaderBytesTotal.With(labels).Add(float64(stats.ImgOptoRespHeaderBytes))
	m.ImgOptoShieldRespBodyBytesTotal.With(labels).Add(float64(stats.ImgOptoShieldRespBodyBytes))
	m.ImgOptoShieldRespHeaderBytesTotal.With(labels).Add(float64(stats.ImgOptoShieldRespHeaderBytes))
	m.ImgOptoShieldTotal.With(labels).Add(float64(stats.ImgOptoShield))
	m.ImgOptoTotal.With(labels).Add(float64(stats.ImgOpto))
	m.ImgOptoTransformRespBodyBytesTotal.With(labels).Add(float64(stats.ImgOptoTransformRespBodyBytes))
	m.ImgOptoTransformRespHeaderBytesTotal.With(labels).Add(float64(stats.ImgOptoTransformRespHeaderBytes))
	m.ImgOptoTransformTotal.With(labels).Add(float64(stats.ImgOptoTransform))
	m.ImgVideoFramesTotal.With(labels).Add(float64(stats.ImgVideoFrames))
	m.ImgVideoRespBodyBytesTotal.With(labels).Add(float64(stats.ImgVideoRespBodyBytes))
	m.ImgVideoRespHeaderBytesTotal.With(labels).Add(float64(stats.ImgVideoRespHeaderBytes))
	m.ImgVideoShieldFramesTotal.With(labels).Add(float64(stats.ImgVideoShieldFrames))
	m.ImgVideoShieldRespBodyBytesTotal.With(labels).Add(float64(stats.ImgVideoShieldRespBodyBytes))
	m.ImgVideoShieldRespHeaderBytesTotal.With(labels).Add(float64(stats.ImgVideoShieldRespHeaderBytes))
	m.ImgVideoShieldTotal.With(labels).Add(float64(stats.ImgVideoShield))
	m.ImgVideoTotal.With(labels).Add(float64(stats.ImgVideo))
	m.KVStoreClassAOperationsTotal.With(labels).Add(float64(stats.KVStoreClassAOperations))
	m.KVStoreClassBOperationsTotal.With(labels).Add(float64(stats.KVStoreClassBOperations))
	m.LogBytesTotal.With(labels).Add(float64(stats.LogBytes))
	m.LoggingTotal.With(labels).Add(float64(stats.Logging))
	m.MissRespBodyBytesTotal.With(labels).Add(float64(stats.MissRespBodyBytes))
	m.MissSubCountTotal.With(labels).Add(float64(stats.MissSubCount))
	m.MissSubTimeTotal.With(labels).Add(float64(stats.MissSubTime))
	m.MissTimeTotal.With(labels).Add(float64(stats.MissTime))
	m.MissesTotal.With(labels).Add(float64(stats.Misses))
	m.OTFPDeliverTimeTotal.With(labels).Add(float64(stats.OTFPDeliverTime))
	m.OTFPManifestTotal.With(labels).Add(float64(stats.OTFPManifest))
	m.OTFPRespBodyBytesTotal.With(labels).Add(float64(stats.OTFPRespBodyBytes))
	m.OTFPRespHeaderBytesTotal.With(labels).Add(float64(stats.OTFPRespHeaderBytes))
	m.OTFPShieldRespBodyBytesTotal.With(labels).Add(float64(stats.OTFPShieldRespBodyBytes))
	m.OTFPShieldRespHeaderBytesTotal.With(labels).Add(float64(stats.OTFPShieldRespHeaderBytes))
	m.OTFPShieldTimeTotal.With(labels).Add(float64(stats.OTFPShieldTime))
	m.OTFPShieldTotal.With(labels).Add(float64(stats.OTFPShield))
	m.OTFPTotal.With(labels).Add(float64(stats.OTFP))
	m.OTFPTransformRespBodyBytesTotal.With(labels).Add(float64(stats.OTFPTransformRespBodyBytes))
	m.OTFPTransformRespHeaderBytesTotal.With(labels).Add(float64(stats.OTFPTransformRespHeaderBytes))
	m.OTFPTransformTimeTotal.With(labels).Add(float64(stats.OTFPTransformTime))
	m.OTFPTransformTotal.With(labels).Add(float64(stats.OTFPTransform))
	m.OriginCacheFetchRespBodyBytesTotal.With(labels).Add(float64(stats.OriginCacheFetchRespBodyBytes))
	m.OriginCacheFetchRespHeaderBytesTotal.With(labels).Add(float64(stats.OriginCacheFetchRespHeaderBytes))
	m.OriginCacheFetchesTotal.With(labels).Add(float64(stats.OriginCacheFetches))
	m.OriginFetchBodyBytesTotal.With(labels).Add(float64(stats.OriginFetchBodyBytes))
	m.OriginFetchHeaderBytesTotal.With(labels).Add(float64(stats.OriginFetchHeaderBytes))
	m.OriginFetchRespBodyBytesTotal.With(labels).Add(float64(stats.OriginFetchRespBodyBytes))
	m.OriginFetchRespHeaderBytesTotal.With(labels).Add(float64(stats.OriginFetchRespHeaderBytes))
	m.OriginFetchesTotal.With(labels).Add(float64(stats.OriginFetches))
	m.OriginRevalidationsTotal.With(labels).Add(float64(stats.OriginRevalidations))
	m.PCITotal.With(labels).Add(float64(stats.PCI))
	m.PassRespBodyBytesTotal.With(labels).Add(float64(stats.PassRespBodyBytes))
	m.PassSubCountTotal.With(labels).Add(float64(stats.PassSubCount))
	m.PassSubTimeTotal.With(labels).Add(float64(stats.PassSubTime))
	m.PassTimeTotal.With(labels).Add(float64(stats.PassTime))
	m.PassesTotal.With(labels).Add(float64(stats.Passes))
	m.Pipe.With(labels).Add(float64(stats.Pipe))
	m.PipeSubCountTotal.With(labels).Add(float64(stats.PipeSubCount))
	m.PipeSubTimeTotal.With(labels).Add(float64(stats.PipeSubTime))
	m.PredeliverSubCountTotal.With(labels).Add(float64(stats.PredeliverSubCount))
	m.PredeliverSubTimeTotal.With(labels).Add(float64(stats.PredeliverSubTime))
	m.PrehashSubCountTotal.With(labels).Add(float64(stats.PrehashSubCount))
	m.PrehashSubTimeTotal.With(labels).Add(float64(stats.PrehashSubTime))
	m.RecvSubCountTotal.With(labels).Add(float64(stats.RecvSubCount))
	m.RecvSubTimeTotal.With(labels).Add(float64(stats.RecvSubTime))
	m.ReqBodyBytesTotal.With(labels).Add(float64(stats.ReqBodyBytes))
	m.ReqHeaderBytesTotal.With(labels).Add(float64(stats.ReqHeaderBytes))
	m.RequestsTotal.With(labels).Add(float64(stats.Requests))
	m.RespBodyBytesTotal.With(labels).Add(float64(stats.RespBodyBytes))
	m.RespHeaderBytesTotal.With(labels).Add(float64(stats.RespHeaderBytes))
	m.RestartTotal.With(labels).Add(float64(stats.Restart))
	m.SegBlockOriginFetchesTotal.With(labels).Add(float64(stats.SegBlockOriginFetches))
	m.SegBlockShieldFetchesTotal.With(labels).Add(float64(stats.SegBlockShieldFetches))
	m.ShieldCacheFetchesTotal.With(labels).Add(float64(stats.ShieldCacheFetches))
	m.ShieldFetchBodyBytesTotal.With(labels).Add(float64(stats.ShieldFetchBodyBytes))
	m.ShieldFetchHeaderBytesTotal.With(labels).Add(float64(stats.ShieldFetchHeaderBytes))
	m.ShieldFetchRespBodyBytesTotal.With(labels).Add(float64(stats.ShieldFetchRespBodyBytes))
	m.ShieldFetchRespHeaderBytesTotal.With(labels).Add(float64(stats.ShieldFetchRespHeaderBytes))
	m.ShieldFetchesTotal.With(labels).Add(float64(stats.ShieldFetches))
	m.ShieldHitRequestsTotal.With(labels).Add(float64(stats.ShieldHitRequests))
	m.ShieldHitRespBodyBytesTotal.With(labels).Add(float64(stats.ShieldHitRespBodyBytes))
	m.ShieldHitRespHeaderBytesTotal.With(labels).Add(float64(stats.ShieldHitRespHeaderBytes))
	m.ShieldMissRequestsTotal.With(labels).Add(float64(stats.ShieldMissRequests))
	m.ShieldMissRespBodyBytesTotal.With(labels).Add(float64(stats.ShieldMissRespBodyBytes))
	m.ShieldMissRespHeaderBytesTotal.With(labels).Add(float64(stats.ShieldMissRespHeaderBytes))
	m.ShieldRespBodyBytesTotal.With(labels).Add(float64(stats.ShieldRespBodyBytes))
	m.ShieldRespHeaderBytesTotal.With(labels).Add(float64(stats.ShieldRespHeaderBytes))
	m.ShieldRevalidationsTotal.With(labels).Add(float64(stats.ShieldRevalidations))
	m.ShieldTotal.With(labels).Add(float64(stats.Shield))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "200"})).Add(float64(stats.Status200))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "204"})).Add(float64(stats.Status204))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "206"})).Add(float64(stats.Status206))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "301"})).Add(float64(stats.Status301))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "302"})).Add(float64(stats.Status302))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "304"})).Add(float64(stats.Status304))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "400"})).Add(float64(stats.Status400))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "401"})).Add(float64(stats.Status401))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "403"})).Add(float64(stats.Status403))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "404"})).Add(float64(stats.Status404))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "406"})).Add(float64(stats.Status406))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "416"})).Add(float64(stats.Status416))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "429"})).Add(float64(stats.Status429))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "500"})).Add(float64(stats.Status500))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "501"})).Add(float64(stats.Status501))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "502"})).Add(float64(stats.Status502))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "503"})).Add(float64(stats.Status503))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "504"})).Add(float64(stats.Status504))
	m.StatusCodeTotal.With(mergeLabels(labels, prometheus.Labels{"status_code": "505"})).Add(float64(stats.Status505))
	m.StatusGroupTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "1xx"})).Add(float64(stats.Status1xx))
	m.StatusGroupTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "2xx"})).Add(float64(stats.Status2xx))
	m.StatusGroupTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "3xx"})).Add(float64(stats.Status3xx))
	m.StatusGroupTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "4xx"})).Add(float64(stats.Status4xx))
	m.StatusGroupTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "5xx"})).Add(float64(stats.Status5xx))
	m.SynthsTotal.With(labels).Add(float64(stats.Synths))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.0"})).Add(float64(stats.TLSv10))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.1"})).Add(float64(stats.TLSv11))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.2"})).Add(float64(stats.TLSv12))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.3"})).Add(float64(stats.TLSv13))
	m.UncacheableTotal.With(labels).Add(float64(stats.Uncacheable))
	m.VclOnComputeHitRequestsTotal.With(labels).Add(float64(stats.VclOnComputeHitRequests))
	m.VclOnComputeMissRequestsTotal.With(labels).Add(float64(stats.VclOnComputeMissRequests))
	m.VclOnComputePassRequestsTotal.With(labels).Add(float64(stats.VclOnComputePassRequests))
	m.VclOnComputeErrorRequestsTotal.With(labels).Add(float64(stats.VclOnComputeErrorRequests))
	m.VclOnComputeSynthRequestsTotal.With(labels).Add(float64(stats.VclOnComputeSynthRequests))
	m.VclOnComputeEdgeHitRequestsTotal.With(labels).Add(float64(stats.VclOnComputeEdgeHitRequests))
	m.VclOnComputeEdgeMissRequestsTotal.With(labels).Add(float64(stats.VclOnComputeEdgeMissRequests))
	m.VideoTotal.With(labels).Add(float64(stats.Video))
	m.WAFBlockedTotal.With(labels).Add(float64(stats.WAFBlocked))
	m.WAFLoggedTotal.With(labels).Add(float64(stats.WAFLogged))
	m.WAFPassedTotal.With(labels).Add(float64(stats.WAFPassed))
	m.WebsocketBackendReqBodyBytesTotal.With(labels).Add(float64(stats.WebsocketBackendReqBodyBytes))
	m.WebsocketBackendReqHeaderBytesTotal.With(labels).Add(float64(stats.WebsocketBackendReqHeaderBytes))
	m.WebsocketBackendRespBodyBytesTotal.With(labels).Add(float64(stats.WebsocketBackendRespBodyBytes))
	m.WebsocketBackendRespHeaderBytesTotal.With(labels).Add(float64(stats.WebsocketBackendRespHeaderBytes))
	m.WebsocketConnTimeMsTotal.With(labels).Add(float64(stats.WebsocketConnTimeMs))
	m.WebsocketReqBodyBytesTotal.With(labels).Add(float64(stats.WebsocketReqBodyBytes))
	m.WebsocketReqHeaderBytesTotal.With(labels).Add(float64(stats.WebsocketReqHeaderBytes))
	m.WebsocketRespBodyBytesTotal.With(labels).Add(float64(stats.WebsocketRespBodyBytes))
	m.WebsocketRespHeaderBytesTotal.With(labels).Add(float64(stats.WebsocketRespHeaderBytes))
	processHistogram(stats.MissHistogram, m.MissDurationSeconds.With(labels))
	processObjectSizes(stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k, stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m, stats.ObjectSize1g, m.ObjectSizeBytes.With(labels))
}

func processHistogram(src map[string]uint64, obs prometheus.Observer) {
	for str, count := range src {
		ms, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		s := float64(ms) / 1e3
		for i := 0; i < int(count); i++ {
			obs.Observe(s)
		}
	}
}

func processObjectSizes(n1k, n10k, n100k, n1m, n10m, n100m, n1g uint64, obs prometheus.Observer) {
	for v, n := range map[uint64]uint64{
		1 * 1024:           n1k,
		10 * 1024:          n10k,
		100 * 1024:         n100k,
		1 * 1000 * 1024:    n1m,
		10 * 1000 * 1024:   n10m,
		100 * 1000 * 1024:  n100m,
		1000 * 1000 * 1024: n1g,
	} {
		for i := uint64(0); i < n; i++ {
			obs.Observe(float64(v))
		}
	}
}

func mergeLabels(ls ...prometheus.Labels) prometheus.Labels {
	labels := make(map[string]string)

	for _, l := range ls {
		for k, v := range l {
			labels[k] = v
		}
	}

	return labels
}
