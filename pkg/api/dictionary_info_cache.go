package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// dictionaryInfo is the Fastly API response for the /info endpoint.
type dictionaryInfo struct {
	Digest      string `json:"digest"`
	ItemCount   int64  `json:"item_count"`
	LastUpdated string `json:"last_updated"` // may be RFC3339 or "2006-01-02 15:04:05"; can be empty
}

// dictionaryResp represents an individual dictionary from the listDictionaries API call.
type dictionaryResp struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

// DictionaryInfoCache polls api.fastly.com/service/%s/version/%d/dictionary and
// api.fastly.com/service/%s/version/%d/dictionary/%s/info and maintains a local cache
// of the dictionary info. That information is exposed as Prometheus metrics.
type DictionaryInfoCache struct {
	client       HTTPClient
	token        string
	logger       log.Logger
	serviceCache *ServiceCache
	enabled      bool

	mtx sync.RWMutex

	// Cached snapshot used by the collector to avoid network I/O on scrape.
	dictionaries []Dictionary
}

// Dictionary holds information about a single dictionary,
// including its item count and unique identifier.
type Dictionary struct {
	ServiceID      string
	ServiceName    string
	Version        int
	DictionaryID   string
	DictionaryName string
	Digest         string
	ItemCount      float64
	LastUpdatedTS  float64
}

// NewDictionaryInfoCache returns an empty cache of dictionary metadata. Use the
// Refresh method to update the cache.
func NewDictionaryInfoCache(client HTTPClient, token string, logger log.Logger, serviceCache *ServiceCache, enabled bool) *DictionaryInfoCache {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &DictionaryInfoCache{
		client:       client,
		token:        token,
		logger:       log.With(logger, "component", "dictionary-info"),
		serviceCache: serviceCache,
		enabled:      enabled,
	}
}

// Enabled returns true if the DictionaryInfoCache is enabled
func (c *DictionaryInfoCache) Enabled() bool { return c.enabled }

// Refresh queries Fastly APIs and rebuilds the in-memory snapshot.
func (c *DictionaryInfoCache) Refresh(ctx context.Context) error {
	if !c.enabled {
		return nil
	}
	out := []Dictionary{}
	for _, s := range c.serviceCache.Services() {
		active := s.Version
		if active <= 0 {
			continue
		}
		dicts, err := c.listDictionaries(ctx, s.ID, active)
		if err != nil {
			level.Warn(c.logger).Log("during", "list dictionaries", "service", s.ID, "err", err)
			return err
		}
		for _, d := range dicts {
			info, err := c.getDictionaryInfo(ctx, s.ID, active, d.ID)
			if err != nil {
				level.Warn(c.logger).Log("during", "get dictionary info", "service", s.ID, "dictionary", d.ID, "err", err)
				continue
			}
			ts := parseFastlyTime(info.LastUpdated)
			if ts == 0 {
				ts = parseFastlyTime(d.UpdatedAt)
			}
			out = append(out, Dictionary{
				ServiceID:      s.ID,
				ServiceName:    s.Name,
				Version:        active,
				DictionaryID:   d.ID,
				DictionaryName: d.Name,
				Digest:         info.Digest,
				ItemCount:      float64(info.ItemCount),
				LastUpdatedTS:  ts,
			})
		}
	}

	c.mtx.Lock()
	c.dictionaries = out
	c.mtx.Unlock()

	level.Debug(c.logger).Log("msg", "refreshed dictionary info", "rows", len(out))
	return nil
}

// Dictionaries returns a copy of the currently cached dictionaries.
func (c *DictionaryInfoCache) Dictionaries() []Dictionary {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	dictionaries := c.dictionaries[:]
	return dictionaries
}

// dictionaryInfoCollector emits metrics from the DictionaryInfoCache snapshot.
type dictionaryInfoCollector struct {
	cache           *DictionaryInfoCache
	itemCountDesc   *prometheus.Desc
	lastUpdatedDesc *prometheus.Desc
}

func (dc *dictionaryInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- dc.itemCountDesc
	ch <- dc.lastUpdatedDesc
}

func (dc *dictionaryInfoCollector) Collect(ch chan<- prometheus.Metric) {
	for _, r := range dc.cache.Dictionaries() {
		labels := []string{
			r.ServiceID,
			r.ServiceName,
			strconv.Itoa(r.Version),
			r.DictionaryID,
			r.DictionaryName,
		}
		ch <- prometheus.MustNewConstMetric(dc.itemCountDesc, prometheus.GaugeValue, r.ItemCount, labels...)

		// Digest is intentionally not exposed as a label or exemplar to avoid cardinality explosion.
		ch <- prometheus.MustNewConstMetric(dc.lastUpdatedDesc, prometheus.GaugeValue, r.LastUpdatedTS, labels...)
	}
}

// Gatherer returns a prometheus.Gatherer backed by a custom collector.
func (c *DictionaryInfoCache) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	if !c.enabled {
		return prometheus.Gatherers{}, nil
	}

	labelKeys := []string{"service_id", "service_name", "version", "dictionary_id", "dictionary_name"}

	itemCountDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "dictionary_item_count"),
		"Number of items in a Fastly edge dictionary.",
		labelKeys, nil,
	)

	lastUpdatedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "dictionary_last_updated_timestamp_seconds"),
		"Unix timestamp (seconds) when the dictionary was last updated (UTC).",
		labelKeys, nil,
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(&dictionaryInfoCollector{
		cache:           c,
		itemCountDesc:   itemCountDesc,
		lastUpdatedDesc: lastUpdatedDesc,
	})

	return reg, nil
}

// --- Fastly API helpers ---

func (c *DictionaryInfoCache) req(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Fastly-Key", c.token)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *DictionaryInfoCache) listDictionaries(ctx context.Context, serviceID string, version int) ([]dictionaryResp, error) {
	u := fmt.Sprintf("https://api.fastly.com/service/%s/version/%d/dictionary", serviceID, version)
	req, err := c.req(ctx, http.MethodGet, u)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// disable if we get a 403 and c.certs is nil (only true on first run)
		if resp.StatusCode == http.StatusForbidden && c.dictionaries == nil {
			c.enabled = false
		}
		return nil, NewError(resp)
	}
	var dicts []dictionaryResp
	return dicts, json.NewDecoder(resp.Body).Decode(&dicts)
}

func (c *DictionaryInfoCache) getDictionaryInfo(ctx context.Context, serviceID string, version int, dictID string) (dictionaryInfo, error) {
	u := fmt.Sprintf("https://api.fastly.com/service/%s/version/%d/dictionary/%s/info", serviceID, version, dictID)
	req, err := c.req(ctx, http.MethodGet, u)
	if err != nil {
		return dictionaryInfo{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return dictionaryInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return dictionaryInfo{}, fmt.Errorf("get dictionary info: %s", resp.Status)
	}
	var info dictionaryInfo
	return info, json.NewDecoder(resp.Body).Decode(&info)
}

// parseFastlyTime parses the common Fastly timestamp formats into a Unix
// timestamp (seconds). Returns 0 if parsing fails.
func parseFastlyTime(s string) float64 {
	if s == "" {
		return 0
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return float64(t.Unix())
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return float64(t.Unix())
	}
	return 0
}
