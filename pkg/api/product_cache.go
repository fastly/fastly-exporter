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
	// ProductDefault represents the standard real-time stats available to all services.
	ProductDefault = "default"

	// ProductOriginInspector represents the origin inspector stats available via the
	// entitlement API.
	ProductOriginInspector = "origin_inspector"

	// ProductDomainInspector represents the domain inspector stats available via the
	// entitlement API.
	ProductDomainInspector = "domain_inspector"
)

// Products is the slice of available products supported by real-time stats.
var Products = []string{ProductDefault, ProductOriginInspector, ProductDomainInspector}

type entitlementsResponse struct {
	Customers []struct {
		Contracts []struct {
			Items []struct {
				ProductID *string `json:"product_id,omitempty"`
			} `json:"items,omitempty"`
		} `json:"contracts"`
	} `json:"customers"`
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

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.fastly.com/entitlements", nil)
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

	var response entitlementsResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding API product response: %w", err)
	}

	activeProducts := make(map[string]interface{})

	for _, customer := range response.Customers {
		for _, contract := range customer.Contracts {
			for _, item := range contract.Items {
				if item.ProductID == nil {
					continue
				}
				activeProducts[*item.ProductID] = true
			}
		}
	}

	for _, product := range Products {
		if product == ProductDefault {
			continue
		}
		_, hasAccess := activeProducts[product]
		level.Debug(p.logger).Log("product", product, "hasAccess", hasAccess)
		p.mtx.Lock()
		p.products[product] = hasAccess
		p.mtx.Unlock()
	}

	return nil
}

// HasAccess takes a product as a string and returns a boolean
// based on the response from the Product API.
func (p *ProductCache) HasAccess(product string) bool {
	if product == ProductDefault {
		return true
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if v, ok := p.products[product]; ok {
		return v
	}
	return true
}
