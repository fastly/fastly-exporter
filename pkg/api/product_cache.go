package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ProductAccess models the response from the Fastly Product Entitlement API
type ProductAccess struct {
	Product   `json:"product"`
	HasAccess bool `json:"has_access"`
}

// Product models product metadata return from the Fastly Product Entitlement API
type Product struct {
	Name string `json:"id"`
}

// ProductCache fetches product information from the Fastly Product Entitlement API
// and stores results in a local cache.
type ProductCache struct {
	client HTTPClient
	token  string
	logger log.Logger

	products map[string]struct{}
}

// NewProductCache returns an empty cache of Product information. Use the Fetch method
// to populate with data.
func NewProductCache(client HTTPClient, token string, logger log.Logger) *ProductCache {
	return &ProductCache{
		client:   client,
		token:    token,
		logger:   logger,
		products: make(map[string]struct{}),
	}
}

// Refresh requests data from the Fastly API and stores data in the cache.
func (p *ProductCache) Refresh(ctx context.Context) error {
	var products = []string{"origin_inspector", "domain_inspector"}

	for _, product := range products {
		uri := fmt.Sprintf("https://api.fastly.com/entitled-products/flonk/%s", product)

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

		var response ProductAccess

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("error decoding API product response: %w", err)
		}

		level.Debug(p.logger).Log("product", response.Name, "hasAccess", response.HasAccess)

		if response.HasAccess {
			p.products[response.Name] = struct{}{}
		}

	}
	return nil
}

// Products returns the list of products
func (p *ProductCache) HasAccess(product string) bool {
	_, ok := p.products[product]
	return ok
}
