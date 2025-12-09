package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// maxCertificatesPageSize is the maximum amount of results that can be requested
// from the api.fastly.com/tls/certificates endpoint.
const maxCertificatesPageSize = 200

// format is the expected time format for the NotAfter timestamp from the Fastly API.
const format = "2006-01-02T15:04:05.000Z"

// CertificateResponse represents the top-level structure of the response from
// the /tls/certificates endpoint.
type CertificateResponse struct {
	Certificates []Certificate     `json:"data"`
	Links        map[string]string `json:"links"`
}

// Certificate holds information about a single TLS certificate,
// including its attributes and unique identifier.
type Certificate struct {
	Attributes Attributes `json:"attributes"`
	ID         string     `json:"id"`
}

// Attributes contains the specific metadata for a TLS certificate,
// such as its Common Name (CN), user-defined name, issuer,
// expiration date (NotAfter), and serial number (SN).
type Attributes struct {
	CN       string `json:"issued_to"`
	Name     string `json:"name"`
	Issuer   string `json:"issuer"`
	NotAfter string `json:"not_after"`
	SN       string `json:"serial_number"`
}

// CertificateCache polls api.fastly.com/tls/certificates and maintains a local cache
// of the returned metadata. That information is exposed as Prometheus metrics.
type CertificateCache struct {
	client  HTTPClient
	token   string
	enabled bool
	logger  log.Logger

	mtx   sync.Mutex
	certs []Certificate
}

// NewCertificateCache returns an empty cache of certificates metadata. Use the
// Refresh method to update the cache.
func NewCertificateCache(client HTTPClient, token string, enabled bool, logger log.Logger) *CertificateCache {
	return &CertificateCache{
		client:  client,
		token:   token,
		enabled: enabled,
		logger:  logger,
	}
}

// Refresh the cache with metadata retreived from the Fastly API.
func (c *CertificateCache) Refresh(ctx context.Context) error {
	if !c.enabled {
		return nil
	}
	begin := time.Now()

	var (
		uri       = fmt.Sprintf("https://api.fastly.com/tls/certificates?page%%5Bnumber%%5D=1&page%%5Bsize%%5D=%d&sort=created_at", maxCertificatesPageSize)
		nextCerts = []Certificate{}
		total     = 0
	)
	for {
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
			// disable if we get a 403 and c.certs is nil (only true on first run)
			if resp.StatusCode == http.StatusForbidden && c.certs == nil {
				c.enabled = false
			}
			return NewError(resp)
		}

		var response CertificateResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("error decoding API certificates response: %w", err)
		}

		nextCerts = append(nextCerts, response.Certificates...)
		total += len(response.Certificates)

		next, err := getNextCertificateLink(response)
		if err != nil {
			break
		}

		uri = next.String()
	}

	level.Debug(c.logger).Log(
		"certificate_refresh_took", time.Since(begin),
		"total_certificate_count", total,
	)

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.certs = nextCerts

	return nil
}

// Certificates returns a copy of the currently cached certificates.
func (c *CertificateCache) Certificates() []Certificate {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	certs := c.certs[:]
	return certs
}

// Gatherer returns a Prometheus gatherer which will yield current metadata
// about Fastly certificates as labels on a gauge metric.
func (c *CertificateCache) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	var (
		fqName      = prometheus.BuildFQName(namespace, subsystem, "cert_expiry_timestamp_seconds")
		help        = "Metadata about Fastly custom TLS certificates."
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

// Describe implements prometheus.Collector
func (c *certificateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect implements prometheus.Collector
func (c *certificateCollector) Collect(ch chan<- prometheus.Metric) {
	for _, cert := range c.cache.Certificates() {
		t, _ := time.Parse(format, cert.Attributes.NotAfter)

		bigISN := new(big.Int)
		bigISN.SetString(cert.Attributes.SN, 10)

		// Convert certificate SN to hex.
		hexSN := fmt.Sprintf("%X", bigISN)

		var (
			desc        = c.desc
			valueType   = prometheus.GaugeValue
			value       = float64(t.Unix())
			labelValues = []string{cert.Attributes.CN, cert.Attributes.Name, cert.ID, cert.Attributes.Issuer, hexSN}
		)
		ch <- prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
	}
}

func getNextCertificateLink(certResp CertificateResponse) (*url.URL, error) {
	next, ok := certResp.Links["next"]
	if !ok {
		return nil, errors.New("no next link found")
	}
	if next == "" {
		return nil, errors.New("next link is empty")
	}
	return url.Parse(next)
}
