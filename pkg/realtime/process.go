package realtime

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// Process updates the metrics with data from the API response.
func Process(response *Response, serviceID, serviceName, serviceVersion string, m *Metrics) {
	for _, d := range response.Data {
		for datacenter, stats := range d.Datacenter {
			processDatacenter(&datacenter, serviceID, serviceName, serviceVersion, stats, m)
		}
	}
}

func ProcessAggregated(response *Response, serviceID, serviceName, serviceVersion string, m *Metrics) {
	for _, d := range response.Data {
		processDatacenter(nil, serviceID, serviceName, serviceVersion, d.Aggregated, m)
	}
}

func processDatacenter(datacenter *string, serviceID, serviceName, serviceVersion string, stats Datacenter, m *Metrics) {
	labels := prometheus.Labels{
		"service_id":   serviceID,
		"service_name": serviceName,
	}

	if datacenter != nil {
		labels["datacenter"] = *datacenter
	}

	m.AttackBlockedReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackBlockedReqBodyBytes))
	m.AttackBlockedReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackBlockedReqHeaderBytes))
	m.AttackLoggedReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackLoggedReqBodyBytes))
	m.AttackLoggedReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackLoggedReqHeaderBytes))
	m.AttackPassedReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackPassedReqBodyBytes))
	m.AttackPassedReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackPassedReqHeaderBytes))
	m.AttackReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackReqBodyBytes))
	m.AttackReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackReqHeaderBytes))
	m.AttackRespSynthBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.AttackRespSynthBytes))
	m.BackendReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.BackendReqBodyBytes))
	m.BackendReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.BackendReqHeaderBytes))
	m.BlacklistedTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Blacklisted))
	m.BodySizeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.BodySize))
	m.ComputeBackendReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendReqBodyBytesTotal))
	m.ComputeBackendReqErrorsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendReqErrorsTotal))
	m.ComputeBackendReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendReqHeaderBytesTotal))
	m.ComputeBackendReqTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendReqTotal))
	m.ComputeBackendRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendRespBodyBytesTotal))
	m.ComputeBackendRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeBackendRespHeaderBytesTotal))
	m.ComputeExecutionTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeExecutionTimeMilliseconds) / 10000.0)
	m.ComputeGlobalsLimitExceededTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeGlobalsLimitExceededTotal))
	m.ComputeGuestErrorsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeGuestErrorsTotal))
	m.ComputeHeapLimitExceededTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeHeapLimitExceededTotal))
	m.ComputeRAMUsedBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRAMUsed))
	m.ComputeReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeReqBodyBytesTotal))
	m.ComputeReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeReqHeaderBytesTotal))
	m.ComputeRequestTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRequestTimeMilliseconds) / 10000.0)
	m.ComputeRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRequests))
	m.ComputeResourceLimitExceedTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeResourceLimitExceedTotal))
	m.ComputeRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRespBodyBytesTotal))
	m.ComputeRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRespHeaderBytesTotal))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "1xx"})).Add(float64(stats.ComputeRespStatus1xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "2xx"})).Add(float64(stats.ComputeRespStatus2xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "3xx"})).Add(float64(stats.ComputeRespStatus3xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "4xx"})).Add(float64(stats.ComputeRespStatus4xx))
	m.ComputeRespStatusTotal.With(mergeLabels(labels, prometheus.Labels{"status_group": "5xx"})).Add(float64(stats.ComputeRespStatus5xx))
	m.ComputeRuntimeErrorsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeRuntimeErrorsTotal))
	m.ComputeStackLimitExceededTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ComputeStackLimitExceededTotal))
	m.DDOSActionBlackholeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionBlackhole))
	m.DDOSActionCloseTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionClose))
	m.DDOSActionLimitStreamsConnectionsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionLimitStreamsConnections))
	m.DDOSActionLimitStreamsRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionLimitStreamsRequests))
	m.DDOSActionTarpitAcceptTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionTarpitAccept))
	m.DDOSActionTarpitTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DDOSActionTarpit))
	m.DeliverSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DeliverSubCount))
	m.DeliverSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.DeliverSubTime))
	m.EdgeHitRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeHitRequests))
	m.EdgeHitRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeHitRespBodyBytes))
	m.EdgeHitRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeHitRespHeaderBytes))
	m.EdgeMissRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeMissRequests))
	m.EdgeMissRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeMissRespBodyBytes))
	m.EdgeMissRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeMissRespHeaderBytes))
	m.EdgeRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeRespBodyBytes))
	m.EdgeRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.EdgeRespHeaderBytes))
	m.EdgeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Edge))
	m.ErrorSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ErrorSubCount))
	m.ErrorSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ErrorSubTime))
	m.ErrorsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Errors))
	m.FanoutBackendReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutBackendReqBodyBytes))
	m.FanoutBackendReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutBackendReqHeaderBytes))
	m.FanoutBackendRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutBackendRespBodyBytes))
	m.FanoutBackendRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutBackendRespHeaderBytes))
	m.FanoutConnTimeMsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutConnTimeMs))
	m.FanoutRecvPublishesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutRecvPublishes))
	m.FanoutReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutReqBodyBytes))
	m.FanoutReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutReqHeaderBytes))
	m.FanoutRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutRespBodyBytes))
	m.FanoutRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutRespHeaderBytes))
	m.FanoutSendPublishesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FanoutSendPublishes))
	m.FetchSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FetchSubCount))
	m.FetchSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.FetchSubTime))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "1"})).Add(float64(stats.Requests - (stats.HTTP2 + stats.HTTP3)))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "2"})).Add(float64(stats.HTTP2))
	m.HTTPTotal.With(mergeLabels(labels, prometheus.Labels{"http_version": "3"})).Add(float64(stats.HTTP3))
	m.HTTP2Total.With(mergeLabels(labels, nil)).Add(float64(stats.HTTP2))
	m.HTTP3Total.With(mergeLabels(labels, nil)).Add(float64(stats.HTTP3))
	m.HashSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HashSubCount))
	m.HashSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HashSubTime))
	m.HeaderSizeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HeaderSize))
	m.HitRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HitRespBodyBytes))
	m.HitSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HitSubCount))
	m.HitSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HitSubTime))
	m.HitsTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.HitsTime))
	m.HitsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Hits))
	m.IPv6Total.With(mergeLabels(labels, nil)).Add(float64(stats.IPv6))
	m.ImgOptoRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoRespBodyBytes))
	m.ImgOptoRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoRespHeaderBytes))
	m.ImgOptoShieldRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoShieldRespBodyBytes))
	m.ImgOptoShieldRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoShieldRespHeaderBytes))
	m.ImgOptoShieldTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoShield))
	m.ImgOptoTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOpto))
	m.ImgOptoTransformRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoTransformRespBodyBytes))
	m.ImgOptoTransformRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoTransformRespHeaderBytes))
	m.ImgOptoTransformTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgOptoTransform))
	m.ImgVideoFramesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoFrames))
	m.ImgVideoRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoRespBodyBytes))
	m.ImgVideoRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoRespHeaderBytes))
	m.ImgVideoShieldFramesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoShieldFrames))
	m.ImgVideoShieldRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoShieldRespBodyBytes))
	m.ImgVideoShieldRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoShieldRespHeaderBytes))
	m.ImgVideoShieldTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideoShield))
	m.ImgVideoTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ImgVideo))
	m.KVStoreClassAOperationsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.KVStoreClassAOperations))
	m.KVStoreClassBOperationsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.KVStoreClassBOperations))
	m.LogBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.LogBytes))
	m.LoggingTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Logging))
	m.MissRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.MissRespBodyBytes))
	m.MissSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.MissSubCount))
	m.MissSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.MissSubTime))
	m.MissTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.MissTime))
	m.MissesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Misses))
	m.OTFPDeliverTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPDeliverTime))
	m.OTFPManifestTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPManifest))
	m.OTFPRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPRespBodyBytes))
	m.OTFPRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPRespHeaderBytes))
	m.OTFPShieldRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPShieldRespBodyBytes))
	m.OTFPShieldRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPShieldRespHeaderBytes))
	m.OTFPShieldTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPShieldTime))
	m.OTFPShieldTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPShield))
	m.OTFPTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFP))
	m.OTFPTransformRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPTransformRespBodyBytes))
	m.OTFPTransformRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPTransformRespHeaderBytes))
	m.OTFPTransformTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPTransformTime))
	m.OTFPTransformTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OTFPTransform))
	m.OriginCacheFetchRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginCacheFetchRespBodyBytes))
	m.OriginCacheFetchRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginCacheFetchRespHeaderBytes))
	m.OriginCacheFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginCacheFetches))
	m.OriginFetchBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginFetchBodyBytes))
	m.OriginFetchHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginFetchHeaderBytes))
	m.OriginFetchRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginFetchRespBodyBytes))
	m.OriginFetchRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginFetchRespHeaderBytes))
	m.OriginFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginFetches))
	m.OriginRevalidationsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.OriginRevalidations))
	m.PCITotal.With(mergeLabels(labels, nil)).Add(float64(stats.PCI))
	m.PassRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PassRespBodyBytes))
	m.PassSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PassSubCount))
	m.PassSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PassSubTime))
	m.PassTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PassTime))
	m.PassesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Passes))
	m.Pipe.With(mergeLabels(labels, nil)).Add(float64(stats.Pipe))
	m.PipeSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PipeSubCount))
	m.PipeSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PipeSubTime))
	m.PredeliverSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PredeliverSubCount))
	m.PredeliverSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PredeliverSubTime))
	m.PrehashSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PrehashSubCount))
	m.PrehashSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.PrehashSubTime))
	m.RecvSubCountTotal.With(mergeLabels(labels, nil)).Add(float64(stats.RecvSubCount))
	m.RecvSubTimeTotal.With(mergeLabels(labels, nil)).Add(float64(stats.RecvSubTime))
	m.ReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ReqBodyBytes))
	m.ReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ReqHeaderBytes))
	m.RequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Requests))
	m.RespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.RespBodyBytes))
	m.RespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.RespHeaderBytes))
	m.RestartTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Restart))
	m.SegBlockOriginFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.SegBlockOriginFetches))
	m.SegBlockShieldFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.SegBlockShieldFetches))
	m.ShieldCacheFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldCacheFetches))
	m.ShieldFetchBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldFetchBodyBytes))
	m.ShieldFetchHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldFetchHeaderBytes))
	m.ShieldFetchRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldFetchRespBodyBytes))
	m.ShieldFetchRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldFetchRespHeaderBytes))
	m.ShieldFetchesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldFetches))
	m.ShieldHitRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldHitRequests))
	m.ShieldHitRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldHitRespBodyBytes))
	m.ShieldHitRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldHitRespHeaderBytes))
	m.ShieldMissRequestsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldMissRequests))
	m.ShieldMissRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldMissRespBodyBytes))
	m.ShieldMissRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldMissRespHeaderBytes))
	m.ShieldRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldRespBodyBytes))
	m.ShieldRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldRespHeaderBytes))
	m.ShieldRevalidationsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.ShieldRevalidations))
	m.ShieldTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Shield))
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
	m.SynthsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Synths))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.0"})).Add(float64(stats.TLSv10))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.1"})).Add(float64(stats.TLSv11))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.2"})).Add(float64(stats.TLSv12))
	m.TLSTotal.With(mergeLabels(labels, prometheus.Labels{"tls_version": "1.3"})).Add(float64(stats.TLSv13))
	m.UncacheableTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Uncacheable))
	m.VideoTotal.With(mergeLabels(labels, nil)).Add(float64(stats.Video))
	m.WAFBlockedTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WAFBlocked))
	m.WAFLoggedTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WAFLogged))
	m.WAFPassedTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WAFPassed))
	m.WebsocketBackendReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketBackendReqBodyBytes))
	m.WebsocketBackendReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketBackendReqHeaderBytes))
	m.WebsocketBackendRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketBackendRespBodyBytes))
	m.WebsocketBackendRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketBackendRespHeaderBytes))
	m.WebsocketConnTimeMsTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketConnTimeMs))
	m.WebsocketReqBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketReqBodyBytes))
	m.WebsocketReqHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketReqHeaderBytes))
	m.WebsocketRespBodyBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketRespBodyBytes))
	m.WebsocketRespHeaderBytesTotal.With(mergeLabels(labels, nil)).Add(float64(stats.WebsocketRespHeaderBytes))
	processHistogram(stats.MissHistogram, m.MissDurationSeconds.With(mergeLabels(labels, nil)))
	processObjectSizes(stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k, stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m, stats.ObjectSize1g, m.ObjectSizeBytes.With(mergeLabels(labels, nil)))
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

func mergeLabels(labels prometheus.Labels, extraLabels prometheus.Labels) prometheus.Labels {
	merged := make(map[string]string)

	for k, v := range labels {
		merged[k] = v
	}

	for k, v := range extraLabels {
		merged[k] = v
	}

	return merged
}
