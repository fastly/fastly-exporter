package prom_test

import (
	"testing"

	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

func TestRegistration(t *testing.T) {
	var (
		namespace = "namespace"
		subsystem = "subsystem"
		registry  = prometheus.NewRegistry()
	)

	{
		_, err := prom.NewMetrics(namespace, subsystem, registry, prom.Stringmap{})
		if err != nil {
			t.Errorf("unexpected error on first construction: %v", err)
		}
	}
	{
		_, err := prom.NewMetrics("alt"+namespace, subsystem, registry, prom.Stringmap{})
		if err != nil {
			t.Errorf("unexpected error on third, alt-namespace construction: %v", err)
		}
	}
}
