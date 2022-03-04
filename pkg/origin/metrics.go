package origin

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics collects all of the origin inspector Prometheus metrics.
type Metrics struct {
	RespBodyBytesTotal   *prometheus.CounterVec
	RespHeaderBytesTotal *prometheus.CounterVec
	ResponsesTotal       *prometheus.CounterVec
	StatusCodeTotal      *prometheus.CounterVec
	StatusGroupTotal     *prometheus.CounterVec
}

// NewMetrics returns a new set of metrics registered to the registerer.
// Only metrics whose names pass the name filter are registered.
func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {
	m := Metrics{
		RespBodyBytesTotal:   prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_body_bytes_total", Help: `Number of body bytes from origin.`}, []string{"service_id", "service_name", "datacenter", "origin"}),
		RespHeaderBytesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "resp_header_bytes_total", Help: `Number of header bytes from origin.`}, []string{"service_id", "service_name", "datacenter", "origin"}),
		ResponsesTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "responses_total", Help: `Number of responses from origin.`}, []string{"service_id", "service_name", "datacenter", "origin"}),
		StatusCodeTotal:      prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_code_total", Help: `Number of responses from origin, by status code e.g. 200, 419.`}, []string{"service_id", "service_name", "datacenter", "origin", "status_code"}),
		StatusGroupTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "status_group_total", Help: `Number of responses from origin, by status group e.g. 1xx, 2xx.`}, []string{"service_id", "service_name", "datacenter", "origin", "status_group"}),
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
