package prom

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
)

// BuildInfoGatherer returns a Prometheus Gatherer that includes build information.
func BuildInfoGatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	registry := prometheus.NewRegistry()
	buildinfoCollector := version.NewCollector(fmt.Sprintf("%s_%s", namespace, subsystem))
	err := registry.Register(buildinfoCollector)
	if err != nil {
		return nil, fmt.Errorf("registering build info collector: %w", err)
	}
	return registry, nil
}
