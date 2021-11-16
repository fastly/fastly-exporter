package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type DatacenterCache struct {
	client HTTPClient
	token  string

	mtx sync.Mutex
	dcs []datacenter
}

type datacenter struct {
	Code  string `json:"code"` // Prometheus label is "datacenter" to make grouping at query time less tedious
	Name  string `json:"name"`
	Group string `json:"group"`
}

func NewDatacenterCache(client HTTPClient, token string) *DatacenterCache {
	return &DatacenterCache{
		client: client,
		token:  token,
	}
}

func (c *DatacenterCache) Refresh() error {
	req, err := http.NewRequest("GET", "https://api.fastly.com/datacenters", nil)
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
		var response struct {
			Msg string `json:"msg"`
		}
		json.NewDecoder(resp.Body).Decode(&response)
		if response.Msg == "" {
			response.Msg = "unknown error"
		}
		return fmt.Errorf("api.fastly.com responded with %s (%s)", resp.Status, response.Msg)
	}

	var response []datacenter
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding API datacenters response: %w", err)
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.dcs = response

	return nil
}

func (c *DatacenterCache) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	var (
		fqName      = prometheus.BuildFQName(namespace, subsystem, "datacenter_info")
		help        = "Metadata about Fastly datacenters."
		labels      = []string{"datacenter", "name", "group"}
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

func (c *DatacenterCache) getDatacenters() []datacenter {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	dcs := make([]datacenter, len(c.dcs))
	copy(dcs, c.dcs)
	return dcs
}

//
//
//

type datacenterCollector struct {
	desc  *prometheus.Desc
	cache *DatacenterCache
}

func (c *datacenterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *datacenterCollector) Collect(ch chan<- prometheus.Metric) {
	for _, dc := range c.cache.getDatacenters() {
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, 1, dc.Code, dc.Name, dc.Group)
	}
}
