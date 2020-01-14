package prom_test

import (
	"testing"

	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

func TestRegistration(t *testing.T) {
	var (
		namespace  = "namespace"
		subsystem  = "subsystem"
		nameFilter = filter.Filter{} // allow all
		registry   = prometheus.NewRegistry()
	)

	{
		_, err := prom.NewMetrics(namespace, subsystem, nameFilter, registry)
		if err != nil {
			t.Errorf("unexpected error on first construction: %v", err)
		}
	}
	{
		_, err := prom.NewMetrics(namespace, subsystem, nameFilter, registry)
		if err == nil {
			t.Error("unexpected success on second construction")
		}
	}
	{
		_, err := prom.NewMetrics("alt"+namespace, subsystem, nameFilter, registry)
		if err != nil {
			t.Errorf("unexpected error on third, alt-namespace construction: %v", err)
		}
	}
}
