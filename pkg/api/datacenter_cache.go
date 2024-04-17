package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Datacenter models a single datacenter as returned by the Fastly API.
type Datacenter struct {
	Code        string      `json:"code"` // Prometheus label is "datacenter" to make grouping at query time less tedious
	Name        string      `json:"name"`
	Group       string      `json:"group"`
	Coördinates Coördinates `json:"coordinates"`
}

// Coördinates of a specific datacenter.
type Coördinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DatacenterCache polls api.fastly.com/datacenters and maintains a local cache
// of the returned metadata. That information is exposed as Prometheus metrics.
type DatacenterCache struct {
	client  HTTPClient
	token   string
	enabled bool

	mtx sync.Mutex
	dcs []Datacenter
}

// NewDatacenterCache returns an empty cache of datacenter metadata. Use the
// Refresh method to update the cache.
func NewDatacenterCache(client HTTPClient, token string, enabled bool) *DatacenterCache {
	return &DatacenterCache{
		client:  client,
		token:   token,
		enabled: enabled,
	}
}

// Refresh the cache with metadata retreived from the Fastly API.
func (c *DatacenterCache) Refresh(ctx context.Context) error {
	if !c.enabled {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.fastly.com/datacenters", nil)
	if err != nil {
		return fmt.Errorf("error constructing API datacenters request: %w", err)
	}

	req.Header.Set("Fastly-Key", c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing API datacenters request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewError(resp)
	}

	var response []Datacenter
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding API datacenters response: %w", err)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Code < response[j].Code
	})

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.dcs = response

	return nil
}

// Datacenters returns a copy of the currently cached datacenters.
func (c *DatacenterCache) Datacenters() []Datacenter {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	dcs := make([]Datacenter, len(c.dcs))
	copy(dcs, c.dcs)
	return dcs
}

// Gatherer returns a Prometheus gatherer which will yield current metadata
// about Fastly datacenters as labels on a gauge metric.
func (c *DatacenterCache) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	var (
		fqName      = prometheus.BuildFQName(namespace, subsystem, "datacenter_info")
		help        = "Metadata about Fastly datacenters."
		labels      = []string{"datacenter", "name", "group", "latitude", "longitude"}
		constLabels = prometheus.Labels{}
		desc        = prometheus.NewDesc(fqName, help, labels, constLabels)
		collector   = &datacenterCollector{desc: desc, cache: c}
	)

	registry := prometheus.NewRegistry()
	if err := registry.Register(collector); err != nil {
		return nil, fmt.Errorf("registering datacenter collector: %w", err)
	}

	return registry, nil
}

// Enabled returns true if the DatacenterCache is enabled
func (c *DatacenterCache) Enabled() bool {
	return c.enabled
}

type datacenterCollector struct {
	desc  *prometheus.Desc
	cache *DatacenterCache
}

func (c *datacenterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *datacenterCollector) Collect(ch chan<- prometheus.Metric) {
	for _, dc := range c.cache.Datacenters() {
		var (
			desc        = c.desc
			valueType   = prometheus.GaugeValue
			value       = 1.0
			latitude    = strconv.FormatFloat(dc.Coördinates.Latitude, 'f', -1, 64)
			longitude   = strconv.FormatFloat(dc.Coördinates.Longitude, 'f', -1, 64)
			labelValues = []string{dc.Code, dc.Name, dc.Group, latitude, longitude}
		)
		ch <- prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
	}
}
