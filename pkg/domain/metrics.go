package domain

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/fastly/fastly-exporter/pkg/filter"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics collects all of the domain inspector Prometheus metrics.
type Metrics struct {
	BackendReqBodyBytesTotal        *prometheus.CounterVec
	BackendReqHeaderBytesTotal      *prometheus.CounterVec
	EdgeHitRequestsTotal            *prometheus.CounterVec
	EdgeMissRequestsTotal           *prometheus.CounterVec
	EdgeRequestsTotal               *prometheus.CounterVec
	EdgeResponseBodyBytesTotal      *prometheus.CounterVec
	EdgeResponseHeaderBytesTotal    *prometheus.CounterVec
	OriginFetchRespBodyBytesTotal   *prometheus.CounterVec
	OriginFetchRespHeaderBytesTotal *prometheus.CounterVec
	OriginFetches                   *prometheus.CounterVec
	OriginStatusCodeTotal           *prometheus.CounterVec
	OriginStatusGroupTotal          *prometheus.CounterVec
	RequestsTotal                   *prometheus.CounterVec
	RespBodyBytesTotal              *prometheus.CounterVec
	RespHeaderBytesTotal            *prometheus.CounterVec
	StatusCodeTotal                 *prometheus.CounterVec
	StatusGroupTotal                *prometheus.CounterVec
}

// NewMetrics returns a new set of metrics registered to the Registerer.
// Only metrics whose names pass the name filter are registered.
func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {
	m := Metrics{
		BackendReqBodyBytesTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_body_bytes_total", Help: "Total body bytes sent to origin."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		BackendReqHeaderBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "bereq_header_bytes_total", Help: "Total header bytes sent to origin."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		EdgeHitRequestsTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_hit_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a hit at the edge."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		EdgeMissRequestsTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_miss_requests_total", Help: "Number of requests sent by end users to Fastly that resulted in a miss at the edge."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		EdgeRequestsTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_requests_total", Help: "Number of requests sent by end users to Fastly."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		EdgeResponseBodyBytesTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_body_bytes_total", Help: "Total body bytes delivered from Fastly to the end user."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		EdgeResponseHeaderBytesTotal:    prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "edge_resp_header_bytes_total", Help: "Total header bytes delivered from Fastly to the end user."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		OriginFetchRespBodyBytesTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_body_bytes", Help: "Total body bytes received from origin."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		OriginFetchRespHeaderBytesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetch_resp_header_bytes", Help: "Total header bytes received from origin."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		OriginFetches:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_fetches", Help: "Number of requests sent to origin."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		OriginStatusCodeTotal:           prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_status_code_total", Help: `Number of responses from origin, by status code e.g. 200, 419.`}, []string{"service_id", "service_name", "datacenter", "domain", "status_code"}),
		OriginStatusGroupTotal:          prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "origin_status_group_total", Help: `Number of responses from origin, by status group e.g. 1xx, 2xx.`}, []string{"service_id", "service_name", "datacenter", "domain", "status_group"}),
		RequestsTotal:                   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "requests_total", Help: "Number of requests processed."}, []string{"service_id", "service_name", "datacenter", "domain"}),
		RespBodyBytesTotal:              prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_body_bytes_total", Help: `Total body bytes delivered.`}, []string{"service_id", "service_name", "datacenter", "domain"}),
		RespHeaderBytesTotal:            prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_header_bytes_total", Help: `Total header bytes delivered.`}, []string{"service_id", "service_name", "datacenter", "domain"}),
		StatusCodeTotal:                 prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_code_total", Help: `Number of responses, by status code e.g. 200, 419.`}, []string{"service_id", "service_name", "datacenter", "domain", "status_code"}),
		StatusGroupTotal:                prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_group_total", Help: `Number of responses, by status group e.g. 1xx, 2xx.`}, []string{"service_id", "service_name", "datacenter", "domain", "status_group"}),
	}

	for i, v := 0, reflect.ValueOf(m); i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(prometheus.Collector)
		if !ok {
			panic(fmt.Errorf("field %d/%d isn't a prometheus.Collector", i+1, v.NumField()))
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
