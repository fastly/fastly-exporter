package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func process(name string, src realtimeResponse, dst *prometheusMetrics) {
	for _, d := range src.Data {
		for datacenter, stats := range d.Datacenter {
			dst.requestsTotal.WithLabelValues(datacenter, name).Add(float64(stats.Requests))
			dst.tlsTotal.WithLabelValues(datacenter, name).Add(float64(stats.TLS))
			dst.shieldTotal.WithLabelValues(datacenter, name).Add(float64(stats.Shield))
			dst.iPv6Total.WithLabelValues(datacenter, name).Add(float64(stats.IPv6))
			dst.imgOptoTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOpto))
			dst.imgOptoShieldTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoShield))
			dst.imgOptoTransformTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoTransform))
			dst.otfpTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFP))
			dst.otfpShieldTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPShield))
			dst.otfpTransformTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPTransform))
			dst.otfpManifestTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPManifest))
			dst.videoTotal.WithLabelValues(datacenter, name).Add(float64(stats.Video))
			dst.pciTotal.WithLabelValues(datacenter, name).Add(float64(stats.PCI))
			dst.loggingTotal.WithLabelValues(datacenter, name).Add(float64(stats.Logging))
			dst.http2Total.WithLabelValues(datacenter, name).Add(float64(stats.HTTP2))
			dst.respHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.RespHeaderBytes))
			dst.headerSizeTotal.WithLabelValues(datacenter, name).Add(float64(stats.HeaderSize))
			dst.respBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.RespBodyBytes))
			dst.bodySizeTotal.WithLabelValues(datacenter, name).Add(float64(stats.BodySize))
			dst.reqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ReqHeaderBytes))
			dst.backendReqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.BackendReqHeaderBytes))
			dst.billedHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.BilledHeaderBytes))
			dst.billedBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.BilledBodyBytes))
			dst.wAFBlockedTotal.WithLabelValues(datacenter, name).Add(float64(stats.WAFBlocked))
			dst.wAFLoggedTotal.WithLabelValues(datacenter, name).Add(float64(stats.WAFLogged))
			dst.wAFPassedTotal.WithLabelValues(datacenter, name).Add(float64(stats.WAFPassed))
			dst.attackReqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackReqHeaderBytes))
			dst.attackReqBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackReqBodyBytes))
			dst.attackRespSynthBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackRespSynthBytes))
			dst.attackLoggedReqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackLoggedReqHeaderBytes))
			dst.attackLoggedReqBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackLoggedReqBodyBytes))
			dst.attackBlockedReqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackBlockedReqHeaderBytes))
			dst.attackBlockedReqBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackBlockedReqBodyBytes))
			dst.attackPassedReqHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackPassedReqHeaderBytes))
			dst.attackPassedReqBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.AttackPassedReqBodyBytes))
			dst.shieldRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ShieldRespHeaderBytes))
			dst.shieldRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ShieldRespBodyBytes))
			dst.otfpRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPRespHeaderBytes))
			dst.otfpRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPRespBodyBytes))
			dst.otfpShieldRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPShieldRespHeaderBytes))
			dst.otfpShieldRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPShieldRespBodyBytes))
			dst.otfpTransformRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPTransformRespHeaderBytes))
			dst.otfpTransformRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPTransformRespBodyBytes))
			dst.otfpShieldTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPShieldTime))
			dst.otfpTransformTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPTransformTime))
			dst.otfpDeliverTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.OTFPDeliverTime))
			dst.imgOptoRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoRespHeaderBytes))
			dst.imgOptoRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoRespBodyBytes))
			dst.imgOptoShieldRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoShieldRespHeaderBytes))
			dst.imgOptoShieldRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoShieldRespBodyBytes))
			dst.imgOptoTransformRespHeaderBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoTransformRespHeaderBytes))
			dst.imgOptoTransformRespBodyBytesTotal.WithLabelValues(datacenter, name).Add(float64(stats.ImgOptoTransformRespBodyBytes))
			dst.statusGroupTotal.WithLabelValues(datacenter, name, "1xx").Add(float64(stats.Status1xx))
			dst.statusGroupTotal.WithLabelValues(datacenter, name, "2xx").Add(float64(stats.Status2xx))
			dst.statusGroupTotal.WithLabelValues(datacenter, name, "3xx").Add(float64(stats.Status3xx))
			dst.statusGroupTotal.WithLabelValues(datacenter, name, "4xx").Add(float64(stats.Status4xx))
			dst.statusGroupTotal.WithLabelValues(datacenter, name, "5xx").Add(float64(stats.Status5xx))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "200").Add(float64(stats.Status200))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "204").Add(float64(stats.Status204))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "301").Add(float64(stats.Status301))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "302").Add(float64(stats.Status302))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "304").Add(float64(stats.Status304))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "400").Add(float64(stats.Status400))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "401").Add(float64(stats.Status401))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "403").Add(float64(stats.Status403))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "404").Add(float64(stats.Status404))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "416").Add(float64(stats.Status416))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "500").Add(float64(stats.Status500))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "501").Add(float64(stats.Status501))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "502").Add(float64(stats.Status502))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "503").Add(float64(stats.Status503))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "504").Add(float64(stats.Status504))
			dst.statusCodeTotal.WithLabelValues(datacenter, name, "505").Add(float64(stats.Status505))
			dst.hitsTotal.WithLabelValues(datacenter, name).Add(float64(stats.Hits))
			dst.missesTotal.WithLabelValues(datacenter, name).Add(float64(stats.Misses))
			dst.passesTotal.WithLabelValues(datacenter, name).Add(float64(stats.Passes))
			dst.synthsTotal.WithLabelValues(datacenter, name).Add(float64(stats.Synths))
			dst.errorsTotal.WithLabelValues(datacenter, name).Add(float64(stats.Errors))
			dst.uncacheableTotal.WithLabelValues(datacenter, name).Add(float64(stats.Uncacheable))
			dst.hitsTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.HitsTime))
			dst.missTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.MissTime))
			dst.passTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.PassTime))
			processHistogram(stats.MissHistogram, dst.missDurationSeconds.WithLabelValues(datacenter, name))
			dst.tlsv12Total.WithLabelValues(datacenter, name).Add(float64(stats.TLSv12))
			processObjectSizes(
				stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k,
				stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m,
				stats.ObjectSize1g, dst.objectSizeBytes.WithLabelValues(datacenter, name),
			)
			dst.recvSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.RecvSubTime))
			dst.recvSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.RecvSubCount))
			dst.hashSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.HashSubTime))
			dst.hashSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.HashSubCount))
			dst.missSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.MissSubTime))
			dst.missSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.MissSubCount))
			dst.fetchSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.FetchSubTime))
			dst.fetchSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.FetchSubCount))
			dst.deliverSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.DeliverSubTime))
			dst.deliverSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.DeliverSubCount))
			dst.hitSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.HitSubTime))
			dst.hitSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.HitSubCount))
			dst.prehashSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.PrehashSubTime))
			dst.prehashSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.PrehashSubCount))
			dst.predeliverSubTimeTotal.WithLabelValues(datacenter, name).Add(float64(stats.PredeliverSubTime))
			dst.predeliverSubCountTotal.WithLabelValues(datacenter, name).Add(float64(stats.PredeliverSubCount))
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
