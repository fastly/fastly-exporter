package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
)

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
	ServiceCacheConfig
	cache serviceCache
}

// ServiceCacheConfig collects the parameters required for a service cache.
type ServiceCacheConfig struct {
	// Client used by managed subscribers to query rt.fastly.com. If not
	// provided, http.DefaultClient is used, which may not include the desired
	// User-Agent, among other things.
	Client HTTPClient

	// Token provided as the Fastly-Key when querying rt.fastly.com.
	Token string

	// IDFilter is a set of specific service IDs that are permitted. If nil or
	// empty, all service IDs will be considered.
	IDFilter StringSet

	// NameFilter filters services based on their name. The zero value for a
	// filter is valid and permits all names.
	NameFilter filter.Filter

	// ShardFilter filters services based on the sharding rules described in the
	// README. The zero value is valid and permits all services.
	ShardFilter Shard

	// Logger is used for runtime diagnostic information.
	// If not provided, a no-op logger is used.
	Logger log.Logger
}

func (c *ServiceCacheConfig) validate() {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	if c.Logger == nil {
		c.Logger = nopLogger
	}
}

// NewServiceCache returns a new, empty cache of services.
func NewServiceCache(c ServiceCacheConfig) *ServiceCache {
	c.validate()
	return &ServiceCache{ServiceCacheConfig: c}
}

// Refresh services and their metadata.
func (c *ServiceCache) Refresh(ctx context.Context) error {
	begin := time.Now()

	var (
		uri     = "https://api.fastly.com/service?page=1&per_page=100"
		total   = 0
		nextgen = map[string]Service{}
	)

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
		if err != nil {
			return fmt.Errorf("error constructing API services request: %w", err)
		}

		req.Header.Set("Fastly-Key", c.Token)
		req.Header.Set("Accept", "application/json")
		resp, err := c.Client.Do(req)
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
			debug := level.Debug(log.With(c.Logger,
				"service_id", s.ID,
				"service_name", s.Name,
				"service_version", s.Version,
			))

			if reject := !c.IDFilter.Empty() && !c.IDFilter.Has(s.ID); reject {
				debug.Log("result", "rejected", "reason", "service ID not explicitly allowed")
				continue
			}

			if reject := !c.NameFilter.Permit(s.Name); reject {
				debug.Log("result", "rejected", "reason", "service name rejected by name filter")
				continue
			}

			if reject := !c.ShardFilter.match(s.ID); reject {
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

	level.Debug(c.Logger).Log(
		"refresh_took", time.Since(begin),
		"total_service_count", total,
		"accepted_service_count", len(nextgen),
	)

	c.cache.update(nextgen, c.Logger)

	return nil
}

// ServiceIDs currently being monitored by the cache.
// The set can change over time.
func (c *ServiceCache) ServiceIDs() (ids []string) {
	return c.cache.getAll()
}

// Metadata returns selected metadata associated with a given service ID.
// If the cache doesn't contain that service ID, found will be false.
func (c *ServiceCache) Metadata(id string) (name string, version int, found bool) {
	return c.cache.getOne(id)
}

//
//
//

type serviceCache struct {
	mtx      sync.RWMutex
	services map[string]Service
}

func (c *serviceCache) update(nextgen map[string]Service, logger log.Logger) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for id, next := range nextgen {
		_, ok := c.services[id]
		if created := !ok; created {
			level.Info(logger).Log("service", "found", "service_id", id, "name", next.Name, "version", next.Version)
		}
	}

	for id, prev := range c.services {
		next, ok := nextgen[id]
		if removed := !ok; removed {
			level.Info(logger).Log("service", "removed", "service_id", id, "name", prev.Name, "version", prev.Version)
		}
		if renamed := ok && prev.Name != next.Name; renamed {
			level.Info(logger).Log("service", "renamed", "service_id", id, "from", prev.Name, "to", next.Name)
		}
		if updated := ok && prev.Version != next.Version; updated {
			level.Info(logger).Log("service", "updated", "service_id", id, "from", prev.Version, "to", next.Version)
		}
	}

	c.services = nextgen
}

func (c *serviceCache) getAll() (ids []string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	ids = make([]string, 0, len(c.services))
	for _, s := range c.services {
		ids = append(ids, s.ID)
	}
	sort.Strings(ids) // mostly for tests
	return ids
}

func (c *serviceCache) getOne(id string) (name string, version int, found bool) {
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

// StringSet is an add-only set of strings.
type StringSet struct{ m map[string]struct{} }

// StringSetWith constructs a string set with an initial set of strings.
func StringSetWith(strs ...string) StringSet {
	var ss StringSet
	ss.Add(strs...)
	return ss
}

// Set adds the value to the set. It's meant to implement flag.Value.
func (ss *StringSet) Set(value string) error {
	ss.Add(value)
	return nil
}

// Add the strs to the set.
func (ss *StringSet) Add(strs ...string) {
	if ss.m == nil {
		ss.m = map[string]struct{}{}
	}
	for _, s := range strs {
		ss.m[s] = struct{}{}
	}
}

// Empty returns true if the set is empty.
func (ss StringSet) Empty() bool {
	return len(ss.m) == 0
}

// Has returns true if `s` is in the set.
func (ss StringSet) Has(s string) bool {
	_, ok := ss.m[s]
	return ok
}

//
//
//

// Shard identifies one exporter instance among many.
type Shard struct{ N, M uint64 }

// ParseShard parses a string of the form "N/M" where N > 0, M > 0, and N < M.
func ParseShard(str string) (s Shard, err error) {
	if i := strings.Index(str, "/"); i >= 0 {
		s.N, _ = strconv.ParseUint(strings.TrimSpace(str[:i]), 10, 64)
		s.M, _ = strconv.ParseUint(strings.TrimSpace(str[i+1:]), 10, 64)
	}
	if s.M <= 0 || s.N <= 0 && s.N > s.M {
		err = fmt.Errorf("%q: invalid format", str)
	}
	return s, err
}

func (s Shard) match(serviceID string) bool {
	switch {
	case s.M == 0:
		return true // the zero value of the type matches all IDs
	case s.M > 0 && s.N > 0 && s.N <= s.M:
		return xxhash.Sum64String(serviceID)%s.M == (s.N - 1)
	default:
		panic(fmt.Errorf("programmer error: invalid shard %v", s))
	}
}
