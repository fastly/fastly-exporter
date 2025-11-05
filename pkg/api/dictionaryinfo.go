// pkg/api/dictionaryinfo.go
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// DictionaryInfo is the Fastly API response for the /info endpoint.
type DictionaryInfo struct {
	Digest      string `json:"digest"`
	ItemCount   int64  `json:"item_count"`
	LastUpdated string `json:"last_updated"` // may be RFC3339 or other formats; can be empty
}

type svc struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Versions []struct {
		Number int  `json:"number"`
		Active bool `json:"active"`
	} `json:"versions"`
}

type dict struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// dictionary item model (subset)
type dictItem struct {
	ItemKey   string `json:"item_key"`
	ItemValue string `json:"item_value"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

type DictionaryInfoCache struct {
	client       *http.Client
	token        string
	logger       log.Logger
	serviceCache *ServiceCache
	enabled      bool

	// Cached snapshot used by the collector to avoid network I/O on scrape.
	snapshot []metricRow
}

type metricRow struct {
	ServiceID      string
	ServiceName    string
	Version        int
	DictionaryID   string
	DictionaryName string
	Digest         string
	ItemCount      float64
	LastUpdatedTS  float64
}

func NewDictionaryInfoCache(client *http.Client, token string, logger log.Logger, serviceCache *ServiceCache, enabled bool) *DictionaryInfoCache {
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

func (c *DictionaryInfoCache) Enabled() bool { return c.enabled }

// parse multiple time layouts commonly seen in Fastly APIs
func parseFastlyTime(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",       // "2025-10-21 18:48:00"
		"2006-01-02T15:04:05",       // without zone
		"2006-01-02 15:04:05Z07:00", // space + zone
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return float64(t.Unix()), true
		}
	}
	return 0, false
}

// GET /service/{sid}/version/{v}/dictionary/{did}/item
func (c *DictionaryInfoCache) listDictionaryItems(ctx context.Context, serviceID string, version int, dictID string) ([]dictItem, error) {
	u := fmt.Sprintf("https://api.fastly.com/service/%s/version/%d/dictionary/%s/item", serviceID, version, dictID)
	req, err := c.req(ctx, http.MethodGet, u)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("list dictionary items: %s", resp.Status)
	}
	var items []dictItem
	return items, json.NewDecoder(resp.Body).Decode(&items)
}

// Refresh queries Fastly APIs and rebuilds the in-memory snapshot.
func (c *DictionaryInfoCache) Refresh(ctx context.Context) error {
	if !c.enabled {
		return nil
	}
	var services []svc
	if c.serviceCache != nil {
		for _, s := range c.serviceCache.Services() {
			active := s.Version
			services = append(services, svc{
				ID:   s.ID,
				Name: s.Name,
				Versions: []struct {
					Number int  `json:"number"`
					Active bool `json:"active"`
				}{
					{
						Number: active,
						Active: active > 0,
					},
				},
			})
		}
	} else {
		var err error
		services, err = c.listServices(ctx)
		if err != nil {
			return err
		}
	}
	var out []metricRow
	for _, s := range services {
		active := -1
		for _, v := range s.Versions {
			if v.Active {
				active = v.Number
				break
			}
		}
		if active < 0 {
			continue
		}
		dicts, err := c.listDictionaries(ctx, s.ID, active)
		if err != nil {
			level.Warn(c.logger).Log("during", "list dictionaries", "service", s.ID, "err", err)
			continue
		}
		for _, d := range dicts {
			info, err := c.getDictionaryInfo(ctx, s.ID, active, d.ID)
			if err != nil {
				level.Warn(c.logger).Log("during", "get dictionary info", "service", s.ID, "dictionary", d.ID, "err", err)
				continue
			}
			var ts float64
			if v, ok := parseFastlyTime(info.LastUpdated); ok {
				ts = v
			} else {
				// fallback: derive from latest item.updated_at (if any)
				if items, err := c.listDictionaryItems(ctx, s.ID, active, d.ID); err == nil && len(items) > 0 {
					var maxTS float64
					for _, it := range items {
						if v2, ok2 := parseFastlyTime(it.UpdatedAt); ok2 && v2 > maxTS {
							maxTS = v2
						}
					}
					ts = maxTS // stays 0 if all missing/unparseable
				}
			}
			out = append(out, metricRow{
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
	c.snapshot = out
	level.Debug(c.logger).Log("msg", "refreshed dictionary info", "rows", len(out))
	return nil
}

// dictionaryInfoCollector emits metrics from the DictionaryInfoCache snapshot.
type dictionaryInfoCollector struct {
	cache           *DictionaryInfoCache
	itemCountDesc   *prometheus.Desc
	lastUpdatedDesc *prometheus.Desc
	digestInfoDesc  *prometheus.Desc
}

func (dc *dictionaryInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- dc.itemCountDesc
	ch <- dc.lastUpdatedDesc
	ch <- dc.digestInfoDesc
}

func (dc *dictionaryInfoCollector) Collect(ch chan<- prometheus.Metric) {
	for _, r := range dc.cache.snapshot {
		labels := []string{
			r.ServiceID,
			r.ServiceName,
			strconv.Itoa(r.Version),
			r.DictionaryID,
			r.DictionaryName,
		}
		ch <- prometheus.MustNewConstMetric(dc.itemCountDesc, prometheus.GaugeValue, r.ItemCount, labels...)
		ch <- prometheus.MustNewConstMetric(dc.lastUpdatedDesc, prometheus.GaugeValue, r.LastUpdatedTS, labels...)
		digestLabels := append(labels, r.Digest)
		ch <- prometheus.MustNewConstMetric(dc.digestInfoDesc, prometheus.GaugeValue, 1, digestLabels...)
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

	digestInfoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "dictionary_digest_info"),
		"Digest of the dictionary content (info metric; value is always 1).",
		append(labelKeys, "digest"), nil,
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(&dictionaryInfoCollector{
		cache:           c,
		itemCountDesc:   itemCountDesc,
		lastUpdatedDesc: lastUpdatedDesc,
		digestInfoDesc:  digestInfoDesc,
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

func (c *DictionaryInfoCache) listServices(ctx context.Context) ([]svc, error) {
	u := "https://api.fastly.com/service"
	req, err := c.req(ctx, http.MethodGet, u)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("list services: %s", resp.Status)
	}
	var svcs []svc
	return svcs, json.NewDecoder(resp.Body).Decode(&svcs)
}

func (c *DictionaryInfoCache) listDictionaries(ctx context.Context, serviceID string, version int) ([]dict, error) {
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
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("list dictionaries: %s", resp.Status)
	}
	var dicts []dict
	return dicts, json.NewDecoder(resp.Body).Decode(&dicts)
}

func (c *DictionaryInfoCache) getDictionaryInfo(ctx context.Context, serviceID string, version int, dictID string) (DictionaryInfo, error) {
	u := fmt.Sprintf("https://api.fastly.com/service/%s/version/%d/dictionary/%s/info", serviceID, version, dictID)
	req, err := c.req(ctx, http.MethodGet, u)
	if err != nil {
		return DictionaryInfo{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return DictionaryInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return DictionaryInfo{}, fmt.Errorf("get dictionary info: %s", resp.Status)
	}
	var info DictionaryInfo
	return info, json.NewDecoder(resp.Body).Decode(&info)
}
