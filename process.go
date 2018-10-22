package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func process(src realtimeResponse, serviceName string, serviceID string, dst *prometheusMetrics) {
	for _, d := range src.Data {
		for datacenter, stats := range d.Datacenter {
			dst.requestsTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Requests))
			dst.tlsTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.TLS))
			dst.shieldTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Shield))
			dst.iPv6Total.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.IPv6))
			dst.imgOptoTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOpto))
			dst.imgOptoShieldTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoShield))
			dst.imgOptoTransformTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoTransform))
			dst.otfpTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFP))
			dst.otfpShieldTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPShield))
			dst.otfpTransformTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPTransform))
			dst.otfpManifestTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPManifest))
			dst.videoTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Video))
			dst.pciTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PCI))
			dst.loggingTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Logging))
			dst.http2Total.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HTTP2))
			dst.respHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.RespHeaderBytes))
			dst.headerSizeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HeaderSize))
			dst.respBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.RespBodyBytes))
			dst.bodySizeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.BodySize))
			dst.reqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ReqHeaderBytes))
			dst.backendReqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.BackendReqHeaderBytes))
			dst.billedHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.BilledHeaderBytes))
			dst.billedBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.BilledBodyBytes))
			dst.wAFBlockedTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.WAFBlocked))
			dst.wAFLoggedTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.WAFLogged))
			dst.wAFPassedTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.WAFPassed))
			dst.attackReqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackReqHeaderBytes))
			dst.attackReqBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackReqBodyBytes))
			dst.attackRespSynthBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackRespSynthBytes))
			dst.attackLoggedReqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackLoggedReqHeaderBytes))
			dst.attackLoggedReqBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackLoggedReqBodyBytes))
			dst.attackBlockedReqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackBlockedReqHeaderBytes))
			dst.attackBlockedReqBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackBlockedReqBodyBytes))
			dst.attackPassedReqHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackPassedReqHeaderBytes))
			dst.attackPassedReqBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.AttackPassedReqBodyBytes))
			dst.shieldRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ShieldRespHeaderBytes))
			dst.shieldRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ShieldRespBodyBytes))
			dst.otfpRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPRespHeaderBytes))
			dst.otfpRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPRespBodyBytes))
			dst.otfpShieldRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPShieldRespHeaderBytes))
			dst.otfpShieldRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPShieldRespBodyBytes))
			dst.otfpTransformRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPTransformRespHeaderBytes))
			dst.otfpTransformRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPTransformRespBodyBytes))
			dst.otfpShieldTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPShieldTime))
			dst.otfpTransformTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPTransformTime))
			dst.otfpDeliverTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.OTFPDeliverTime))
			dst.imgOptoRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoRespHeaderBytes))
			dst.imgOptoRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoRespBodyBytes))
			dst.imgOptoShieldRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoShieldRespHeaderBytes))
			dst.imgOptoShieldRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoShieldRespBodyBytes))
			dst.imgOptoTransformRespHeaderBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoTransformRespHeaderBytes))
			dst.imgOptoTransformRespBodyBytesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.ImgOptoTransformRespBodyBytes))
			dst.statusGroupTotal.WithLabelValues(serviceName, serviceID, datacenter, "1xx").Add(float64(stats.Status1xx))
			dst.statusGroupTotal.WithLabelValues(serviceName, serviceID, datacenter, "2xx").Add(float64(stats.Status2xx))
			dst.statusGroupTotal.WithLabelValues(serviceName, serviceID, datacenter, "3xx").Add(float64(stats.Status3xx))
			dst.statusGroupTotal.WithLabelValues(serviceName, serviceID, datacenter, "4xx").Add(float64(stats.Status4xx))
			dst.statusGroupTotal.WithLabelValues(serviceName, serviceID, datacenter, "5xx").Add(float64(stats.Status5xx))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "200").Add(float64(stats.Status200))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "204").Add(float64(stats.Status204))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "301").Add(float64(stats.Status301))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "302").Add(float64(stats.Status302))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "304").Add(float64(stats.Status304))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "400").Add(float64(stats.Status400))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "401").Add(float64(stats.Status401))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "403").Add(float64(stats.Status403))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "404").Add(float64(stats.Status404))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "416").Add(float64(stats.Status416))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "500").Add(float64(stats.Status500))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "501").Add(float64(stats.Status501))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "502").Add(float64(stats.Status502))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "503").Add(float64(stats.Status503))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "504").Add(float64(stats.Status504))
			dst.statusCodeTotal.WithLabelValues(serviceName, serviceID, datacenter, "505").Add(float64(stats.Status505))
			dst.hitsTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Hits))
			dst.missesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Misses))
			dst.passesTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Passes))
			dst.synthsTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Synths))
			dst.errorsTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Errors))
			dst.uncacheableTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.Uncacheable))
			dst.hitsTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HitsTime))
			dst.missTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.MissTime))
			dst.passTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PassTime))
			processHistogram(stats.MissHistogram, dst.missDurationSeconds.WithLabelValues(serviceName, serviceID, datacenter))
			dst.tlsv12Total.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.TLSv12))
			processObjectSizes(
				stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k,
				stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m,
				stats.ObjectSize1g, dst.objectSizeBytes.WithLabelValues(serviceName, serviceID, datacenter),
			)
			dst.recvSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.RecvSubTime))
			dst.recvSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.RecvSubCount))
			dst.hashSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HashSubTime))
			dst.hashSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HashSubCount))
			dst.missSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.MissSubTime))
			dst.missSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.MissSubCount))
			dst.fetchSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.FetchSubTime))
			dst.fetchSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.FetchSubCount))
			dst.deliverSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.DeliverSubTime))
			dst.deliverSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.DeliverSubCount))
			dst.hitSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HitSubTime))
			dst.hitSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.HitSubCount))
			dst.prehashSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PrehashSubTime))
			dst.prehashSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PrehashSubCount))
			dst.predeliverSubTimeTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PredeliverSubTime))
			dst.predeliverSubCountTotal.WithLabelValues(serviceName, serviceID, datacenter).Add(float64(stats.PredeliverSubCount))
		}
	}
}

func processHistogram(src map[string]uint64, dst prometheus.Observer) {
	for str, count := range src {
		ms, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		s := float64(ms) / 1e3
		for i := 0; i < int(count); i++ {
			dst.Observe(s)
		}
	}
}

func processObjectSizes(n1k, n10k, n100k, n1m, n10m, n100m, n1g uint64, dst prometheus.Observer) {
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
			dst.Observe(float64(v))
		}
	}
}
