package rt

import (
	"strconv"

	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// process interprets the data in the realtime response, and feeds the
// interpreted results to the Prometheus metrics as observations.
func process(src realtimeResponse, serviceID, serviceName, serviceVersion string, m *prom.Metrics) {
	for _, d := range src.Data {
		for datacenter, stats := range d.Datacenter {
			m.ServiceInfo.WithLabelValues(serviceID, serviceName, serviceVersion).Set(1)
			m.RequestsTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Requests))
			m.TLSTotal.WithLabelValues(serviceID, serviceName, datacenter, "any").Add(float64(stats.TLS))
			m.TLSTotal.WithLabelValues(serviceID, serviceName, datacenter, "v10").Add(float64(stats.TLSv10))
			m.TLSTotal.WithLabelValues(serviceID, serviceName, datacenter, "v11").Add(float64(stats.TLSv11))
			m.TLSTotal.WithLabelValues(serviceID, serviceName, datacenter, "v12").Add(float64(stats.TLSv12))
			m.TLSTotal.WithLabelValues(serviceID, serviceName, datacenter, "v13").Add(float64(stats.TLSv13))
			m.ShieldTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Shield))
			m.IPv6Total.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.IPv6))
			m.ImgOptoTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOpto))
			m.ImgOptoShieldTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoShield))
			m.ImgOptoTransformTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoTransform))
			m.OTFPTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFP))
			m.OTFPShieldTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPShield))
			m.OTFPTransformTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPTransform))
			m.OTFPManifestTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPManifest))
			m.VideoTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Video))
			m.PCITotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PCI))
			m.LoggingTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Logging))
			m.HTTP2Total.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HTTP2))
			m.RespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.RespHeaderBytes))
			m.HeaderSizeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HeaderSize))
			m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.RespBodyBytes))
			m.BodySizeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.BodySize))
			m.ReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ReqHeaderBytes))
			m.BackendReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.BackendReqHeaderBytes))
			m.BilledHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.BilledHeaderBytes))
			m.BilledBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.BilledBodyBytes))
			m.WAFBlockedTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.WAFBlocked))
			m.WAFLoggedTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.WAFLogged))
			m.WAFPassedTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.WAFPassed))
			m.AttackReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackReqHeaderBytes))
			m.AttackReqBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackReqBodyBytes))
			m.AttackRespSynthBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackRespSynthBytes))
			m.AttackLoggedReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackLoggedReqHeaderBytes))
			m.AttackLoggedReqBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackLoggedReqBodyBytes))
			m.AttackBlockedReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackBlockedReqHeaderBytes))
			m.AttackBlockedReqBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackBlockedReqBodyBytes))
			m.AttackPassedReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackPassedReqHeaderBytes))
			m.AttackPassedReqBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.AttackPassedReqBodyBytes))
			m.ShieldRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ShieldRespHeaderBytes))
			m.ShieldRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ShieldRespBodyBytes))
			m.OTFPRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPRespHeaderBytes))
			m.OTFPRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPRespBodyBytes))
			m.OTFPShieldRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPShieldRespHeaderBytes))
			m.OTFPShieldRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPShieldRespBodyBytes))
			m.OTFPTransformRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPTransformRespHeaderBytes))
			m.OTFPTransformRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPTransformRespBodyBytes))
			m.OTFPShieldTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPShieldTime))
			m.OTFPTransformTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPTransformTime))
			m.OTFPDeliverTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.OTFPDeliverTime))
			m.ImgOptoRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoRespHeaderBytes))
			m.ImgOptoRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoRespBodyBytes))
			m.ImgOptoShieldRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoShieldRespHeaderBytes))
			m.ImgOptoShieldRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoShieldRespBodyBytes))
			m.ImgOptoTransformRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoTransformRespHeaderBytes))
			m.ImgOptoTransformRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.ImgOptoTransformRespBodyBytes))
			m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, "1xx").Add(float64(stats.Status1xx))
			m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, "2xx").Add(float64(stats.Status2xx))
			m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, "3xx").Add(float64(stats.Status3xx))
			m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, "4xx").Add(float64(stats.Status4xx))
			m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, "5xx").Add(float64(stats.Status5xx))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "200").Add(float64(stats.Status200))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "204").Add(float64(stats.Status204))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "301").Add(float64(stats.Status301))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "302").Add(float64(stats.Status302))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "304").Add(float64(stats.Status304))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "400").Add(float64(stats.Status400))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "401").Add(float64(stats.Status401))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "403").Add(float64(stats.Status403))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "404").Add(float64(stats.Status404))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "416").Add(float64(stats.Status416))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "500").Add(float64(stats.Status500))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "501").Add(float64(stats.Status501))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "502").Add(float64(stats.Status502))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "503").Add(float64(stats.Status503))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "504").Add(float64(stats.Status504))
			m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, "505").Add(float64(stats.Status505))
			m.HitsTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Hits))
			m.MissesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Misses))
			m.PassesTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Passes))
			m.SynthsTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Synths))
			m.ErrorsTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Errors))
			m.UncacheableTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Uncacheable))
			m.HitsTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HitsTime))
			m.MissTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.MissTime))
			m.PassTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PassTime))
			processHistogram(stats.MissHistogram, m.MissDurationSeconds.WithLabelValues(serviceID, serviceName, datacenter))
			processObjectSizes(
				stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k,
				stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m,
				stats.ObjectSize1g, m.ObjectSizeBytes.WithLabelValues(serviceID, serviceName, datacenter),
			)
			m.RecvSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.RecvSubTime))
			m.RecvSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.RecvSubCount))
			m.HashSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HashSubTime))
			m.HashSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HashSubCount))
			m.MissSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.MissSubTime))
			m.MissSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.MissSubCount))
			m.FetchSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.FetchSubTime))
			m.FetchSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.FetchSubCount))
			m.DeliverSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.DeliverSubTime))
			m.DeliverSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.DeliverSubCount))
			m.HitSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HitSubTime))
			m.HitSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.HitSubCount))
			m.PrehashSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PrehashSubTime))
			m.PrehashSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PrehashSubCount))
			m.PredeliverSubTimeTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PredeliverSubTime))
			m.PredeliverSubCountTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.PredeliverSubCount))
		}
	}
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
