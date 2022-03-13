package prom

import (
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/origin"
	"github.com/peterbourgon/fastly-exporter/pkg/realtime"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ServiceInfo            *prometheus.GaugeVec
	LastSuccessfulResponse *prometheus.GaugeVec
	Realtime               *realtime.Metrics
	Origin                 *origin.Metrics
}

func NewMetrics(namespace, deprecatedSubsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {
	var (
		serviceInfo            = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: deprecatedSubsystem, Name: "service_info", Help: "Static gauge with service ID, name, and version information."}, []string{"service_id", "service_name", "service_version"})
		lastSuccessfulResponse = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: deprecatedSubsystem, Name: "last_successful_response", Help: "Unix timestamp of the last successful response received from the real-time stats API."}, []string{"service_id", "service_name"})
	)
	r.MustRegister(serviceInfo, lastSuccessfulResponse)

	return &Metrics{
		ServiceInfo:            serviceInfo,
		LastSuccessfulResponse: lastSuccessfulResponse,
		Realtime:               realtime.NewMetrics(namespace, deprecatedSubsystem, nameFilter, r), // TODO(pb): change this to "rt" or "realtime"
		Origin:                 origin.NewMetrics(namespace, "origin", nameFilter, r),
	}
}
