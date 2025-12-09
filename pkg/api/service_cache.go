package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// maxServicePageSize is the maximum amount of results that can be requested
// from the api.fastly.com/service endpoint.
const maxServicePageSize = 1000

// Service metadata associated with a single service.
// Also serves as a DTO for api.fastly.com/service.
type Service struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version int    `json:"version"`
}

// ServiceCache polls api.fastly.com/service to keep metadata about
// one or more service IDs up-to-date.
type ServiceCache struct {
	client HTTPClient
	token  string

	serviceIDs stringSet
	nameFilter filter.Filter
	shard      shardSlice
	logger     log.Logger

	mtx      sync.RWMutex
	services map[string]Service
}

// NewServiceCache returns an empty cache of service metadata. By default, it
// will fetch metadata about all services available to the provided token. Use
// options to restrict which services the cache should manage.
func NewServiceCache(client HTTPClient, token string, options ...ServiceCacheOption) *ServiceCache {
	c := &ServiceCache{
		client: client,
		token:  token,
		logger: log.NewNopLogger(),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// ServiceCacheOption provides some additional behavior to a service cache.
// Options that restrict which services are cached combine with AND semantics.
type ServiceCacheOption func(*ServiceCache)

// WithExplicitServiceIDs restricts the cache to fetch metadata only for the
// provided service IDs. By default, all service IDs available to the provided
// token are allowed.
func WithExplicitServiceIDs(ids ...string) ServiceCacheOption {
	return func(c *ServiceCache) { c.serviceIDs = newStringSet(ids) }
}

// WithNameFilter restricts the cache to fetch metadata only for the services
// whose names pass the provided filter. By default, no name filtering occurs.
func WithNameFilter(f filter.Filter) ServiceCacheOption {
	return func(c *ServiceCache) { c.nameFilter = f }
}

// WithShard restricts the cache to fetch metadata only for those services whose
// IDs, when hashed and taken modulo m, equal (n-1). By default, no sharding
// occurs.
//
// This option is designed to allow users to split accounts (tokens) that have a
// large number of services across multiple exporter processes. For example, to
// split across 3 processes, each process would set n={1,2,3} and m=3.
func WithShard(n, m uint64) ServiceCacheOption {
	return func(c *ServiceCache) { c.shard = shardSlice{n, m} }
}

// WithLogger sets the logger used by the cache during refresh.
// By default, no log events are emitted.
func WithLogger(logger log.Logger) ServiceCacheOption {
	return func(c *ServiceCache) { c.logger = logger }
}

// Refresh services and their metadata.
func (c *ServiceCache) Refresh(ctx context.Context) error {
	begin := time.Now()

	var (
		uri     = fmt.Sprintf("https://api.fastly.com/service?page=1&per_page=%d&filter%%5Binclude_versions%%5D=false", maxServicePageSize)
		total   = 0
		nextgen = map[string]Service{}
	)

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
		if err != nil {
			return fmt.Errorf("error constructing API services request: %w", err)
		}

		req.Header.Set("Fastly-Key", c.token)
		req.Header.Set("Accept", "application/json")
		resp, err := c.client.Do(req)
		if err != nil {
			return fmt.Errorf("error executing API services request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return NewError(resp)
		}

		var response []Service
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("error decoding API services response: %w", err)
		}
		total += len(response)

		for _, s := range response {
			debug := level.Debug(log.With(c.logger,
				"service_id", s.ID,
				"service_name", s.Name,
				"service_version", s.Version,
			))

			if reject := !c.serviceIDs.empty() && !c.serviceIDs.has(s.ID); reject {
				debug.Log("result", "rejected", "reason", "service ID not explicitly allowed")
				continue
			}

			if reject := !c.nameFilter.Permit(s.Name); reject {
				debug.Log("result", "rejected", "reason", "service name rejected by name filter")
				continue
			}

			if reject := !c.shard.match(s.ID); reject {
				debug.Log("result", "rejected", "reason", "service ID in different shard")
				continue
			}

			debug.Log("result", "accepted")
			nextgen[s.ID] = s
		}

		next, err := GetNextLink(resp)
		if err != nil {
			break
		}

		uri = next.String()
	}

	level.Debug(c.logger).Log(
		"refresh_took", time.Since(begin),
		"total_service_count", total,
		"accepted_service_count", len(nextgen),
	)

	c.mtx.Lock()
	for id, next := range nextgen {
		_, ok := c.services[id]
		if created := !ok; created {
			level.Info(c.logger).Log("service", "found", "service_id", id, "name", next.Name, "version", next.Version)
		}
	}
	for id, prev := range c.services {
		next, ok := nextgen[id]
		if removed := !ok; removed {
			level.Info(c.logger).Log("service", "removed", "service_id", id, "name", prev.Name, "version", prev.Version)
		}
		if renamed := ok && prev.Name != next.Name; renamed {
			level.Info(c.logger).Log("service", "renamed", "service_id", id, "from", prev.Name, "to", next.Name)
		}
		if updated := ok && prev.Version != next.Version; updated {
			level.Info(c.logger).Log("service", "updated", "service_id", id, "from", prev.Version, "to", next.Version)
		}
	}
	c.services = nextgen
	c.mtx.Unlock()

	return nil
}

// ServiceIDs currently being monitored by the cache.
// The set can change over time.
func (c *ServiceCache) ServiceIDs() (ids []string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	ids = make([]string, 0, len(c.services))
	for _, s := range c.services {
		ids = append(ids, s.ID)
	}
	sort.Strings(ids) // mostly for tests
	return ids
}

// Services returns all services currently being monitored by the cache.
func (c *ServiceCache) Services() []Service {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	services := make([]Service, 0, len(c.services))
	for _, s := range c.services {
		services = append(services, s)
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].ID < services[j].ID
	})

	return services
}

// Metadata returns selected metadata associated with a given service ID.
// If the cache doesn't contain that service ID, found will be false.
func (c *ServiceCache) Metadata(id string) (name string, version int, found bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if s, ok := c.services[id]; ok {
		name, version, found = s.Name, s.Version, true
	}
	return name, version, found
}

//
//
//

type stringSet map[string]struct{}

func newStringSet(initial []string) stringSet {
	ss := stringSet{}
	for _, s := range initial {
		ss[s] = struct{}{}
	}
	return ss
}

func (ss stringSet) empty() bool {
	return len(ss) == 0
}

func (ss stringSet) has(s string) bool {
	_, ok := ss[s]
	return ok
}

type shardSlice struct{ n, m uint64 }

func (ss shardSlice) match(serviceID string) bool {
	if ss.m == 0 {
		return true // the zero value of the type matches all IDs
	}

	if ss.n == 0 {
		panic("programmer error: shard with n = 0, m != 0")
	}

	h := xxhash.New()
	fmt.Fprint(h, serviceID)
	return h.Sum64()%uint64(ss.m) == uint64(ss.n-1)
}
