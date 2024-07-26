package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// maxCertificatesPageSize is the maximum amount of results that can be requested
// from the api.fastly.com/tls/certificates endpoint.
const maxCertificatesPageSize = 1000

type Certificates struct {
	Certificate        []Certificate      `json:"data"`
}

type Certificate struct {
	Attributes        Attributes      `json:"attributes"`
	Id                string          `json:"id"`
}

type Attributes struct {
	CN          string      `json:"issued_to"`
	Name        string      `json:"name"`
	Issuer      string      `json:"issuer"`
	Not_after   string      `json:"not_after"`
	SN          string      `json:"serial_number"`
}

// CertificateCache polls api.fastly.com/tls/certificates and maintains a local cache
// of the returned metadata. That information is exposed as Prometheus metrics.
type CertificateCache struct {
	client  HTTPClient
	token   string
	enabled bool

	mtx sync.Mutex
	certs Certificates
}

// NewCertificateCache returns an empty cache of certificates metadata. Use the
// Refresh method to update the cache.
func NewCertificateCache(client HTTPClient, token string, enabled bool) *CertificateCache {
	return &CertificateCache{
		client:  client,
		token:   token,
		enabled: enabled,
	}
}

// Refresh the cache with metadata retreived from the Fastly API.
func (c *CertificateCache) Refresh(ctx context.Context) error {
	if !c.enabled {
		return nil
	}

	// TODO: Implement additional requests for next pages if there are more
	// TLS certificates than maxCertificatesPageSize
	var uri string = fmt.Sprintf("https://api.fastly.com/tls/certificates?page%%5Bnumber%%5D=1&page%%5Bsize%%5D=%d&sort=created_at", maxCertificatesPageSize)

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return fmt.Errorf("error constructing API certificates request: %w", err)
	}

	req.Header.Set("Fastly-Key", c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing API certificates request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewError(resp)
	}

	var response Certificates
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding API certificates response: %w", err)
	}

	sort.Slice(response.Certificate, func(i, j int) bool {
		return response.Certificate[i].Attributes.CN < response.Certificate[j].Attributes.CN
	})

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.certs = response

	return nil
}

// Certificates returns a copy of the currently cached certificates.
func (c *CertificateCache) Certificates() Certificates {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	certs := c.certs
	return certs
}

// Gatherer returns a Prometheus gatherer which will yield current metadata
// about Fastly certificates as labels on a gauge metric.
func (c *CertificateCache) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	var (
		fqName      = prometheus.BuildFQName(namespace, subsystem, "cert_expiry_timestamp_seconds")
		help        = "Metadata about Fastly certificates."
		labels      = []string{"cn", "name", "id", "issuer", "sn"}
		constLabels = prometheus.Labels{}
		desc        = prometheus.NewDesc(fqName, help, labels, constLabels)
		collector   = &certificateCollector{desc: desc, cache: c}
	)

	registry := prometheus.NewRegistry()
	if err := registry.Register(collector); err != nil {
		return nil, fmt.Errorf("registering certificate collector: %w", err)
	}

	return registry, nil
}

// Enabled returns true if the CertificateCache is enabled
func (c *CertificateCache) Enabled() bool {
	return c.enabled
}

type certificateCollector struct {
	desc  *prometheus.Desc
	cache *CertificateCache
}

func (c *certificateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *certificateCollector) Collect(ch chan<- prometheus.Metric) {
	for _, cert := range c.cache.Certificates().Certificate {
		format := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(format, cert.Attributes.Not_after)
		var (
			desc        = c.desc
			valueType   = prometheus.GaugeValue
			value       = float64(t.Unix())
			labelValues = []string{cert.Attributes.CN, cert.Attributes.Name, cert.Id, cert.Attributes.Issuer, cert.Attributes.SN}
		)
		ch <- prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
	}
}
