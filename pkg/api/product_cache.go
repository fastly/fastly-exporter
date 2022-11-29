package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	Default         = "default"
	OriginInspector = "origin_inspector"
)

var Products = []string{Default, OriginInspector}

// Product models the response from the Fastly Product Entitlement API
type Product struct {
	HasAccess bool `json:"has_access"`
	Meta      struct {
		Name string `json:"id"`
	} `json:"product"`
}

// ProductCache fetches product information from the Fastly Product Entitlement API
// and stores results in a local cache.
type ProductCache struct {
	client HTTPClient
	token  string
	logger log.Logger

	mtx      sync.Mutex
	products map[string]bool
}

// NewProductCache returns an empty cache of Product information. Use the Refresh method
// to populate with data.
func NewProductCache(client HTTPClient, token string, logger log.Logger) *ProductCache {
	return &ProductCache{
		client:   client,
		token:    token,
		logger:   logger,
		products: make(map[string]bool),
	}
}

// Refresh requests data from the Fastly API and stores data in the cache.
func (p *ProductCache) Refresh(ctx context.Context) error {
	for _, product := range Products {
		if product == Default {
			continue
		}
		uri := fmt.Sprintf("https://api.fastly.com/entitled-products/%s", product)

		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)

		if err != nil {
			return fmt.Errorf("error constructing API product request: %w", err)
		}

		req.Header.Set("Fastly-Key", p.token)
		req.Header.Set("Accept", "application/json")
		resp, err := p.client.Do(req)

		if err != nil {
			return fmt.Errorf("error executing API product request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return NewError(resp)
		}

		var response Product

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("error decoding API product response: %w", err)
		}

		level.Debug(p.logger).Log("product", response.Meta.Name, "hasAccess", response.HasAccess)

		p.mtx.Lock()
		p.products[response.Meta.Name] = response.HasAccess
		p.mtx.Unlock()

	}

	return nil
}

// HasAccess takes a product as a string and returns a boolean
// based on the response from the Product API.
func (p *ProductCache) HasAccess(product string) bool {
	if product == Default {
		return true
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if v, ok := p.products[product]; ok {
		return v
	} else {
		return false
	}
}
