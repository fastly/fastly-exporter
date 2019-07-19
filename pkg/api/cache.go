package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

// HTTPClient is a consumer contract for the cache.
// It models a concrete http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Service metadata associated with a single service.
// Also serves as a DTO for api.fastly.com/service.
type Service struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version int    `json:"version"`
}

// Cache polls api.fastly.com/service to keep metadata about
// one or more service IDs up-to-date.
type Cache struct {
	token     string
	whitelist stringset
	match     *regexp.Regexp
	shard     shardSlice
	logger    log.Logger

	mtx      sync.RWMutex
	services map[string]Service
}

// NewCache returns an empty cache of service metadata. By default, it will
// fetch metadata about all services available to the provided token. Use
// options to restrict which services the cache should manage.
func NewCache(token string, options ...CacheOption) *Cache {
	c := &Cache{
		token:  token,
		logger: log.NewNopLogger(),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// CacheOption provides some additional behavior to a cache. Options that
// restrict which services are cached combine with AND semantics.
type CacheOption func(*Cache)

// WithExplicitServiceIDs restricts the cache to fetch metadata only for the
// provided service IDs. By default, all service IDs available to the provided
// token are allowed.
func WithExplicitServiceIDs(ids ...string) CacheOption {
	return func(c *Cache) { c.whitelist = newStringSet(ids) }
}

// WithNameMatching restricts the cache to fetch metadata only for the
// services whose names match the provided regexp. By default,
// no name filtering occurs.
func WithNameMatching(re *regexp.Regexp) CacheOption {
	return func(c *Cache) { c.match = re }
}

// WithShard restricts the cache to fetch metadata only for those services whose
// IDs, when hashed and taken modulo m, equal (n-1). By default, no sharding
// occurs.
//
// This option is designed to allow users to split accounts (tokens) that have a
// large number of services across multiple exporter processes. For example, to
// split across 3 processes, each process would set n={1,2,3} and m=3.
func WithShard(n, m uint64) CacheOption {
	return func(c *Cache) { c.shard = shardSlice{n, m} }
}

// WithLogger sets the logger used by the cache during refresh.
// By default, no log events are emitted.
func WithLogger(logger log.Logger) CacheOption {
	return func(c *Cache) { c.logger = logger }
}

// Refresh services and their metadata.
func (c *Cache) Refresh(client HTTPClient) error {
	begin := time.Now()

	req, err := http.NewRequest("GET", "https://api.fastly.com/service", nil)
	if err != nil {
		return errors.Wrap(err, "error constructing API services request")
	}

	req.Header.Set("Fastly-Key", c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing API services request")
	}
	defer resp.Body.Close()

	var response []Service
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "error decoding API services response")
	}

	nextgen := map[string]Service{}
	for _, s := range response {
		debug := level.Debug(log.With(c.logger,
			"service_id", s.ID,
			"service_name", s.Name,
			"service_version", s.Version,
		))

		if reject := !c.whitelist.empty() && !c.whitelist.has(s.ID); reject {
			debug.Log("result", "rejected", "cause", "not in service ID whitelist")
			continue
		}

		if reject := c.match != nil && !c.match.MatchString(s.Name); reject {
			debug.Log("result", "rejected", "cause", "failed name regexp")
			continue
		}

		if reject := !c.shard.match(s.ID); reject {
			debug.Log("result", "rejected", "cause", "service ID in different shard")
			continue
		}

		debug.Log("result", "accepted")
		nextgen[s.ID] = s
	}

	level.Debug(c.logger).Log(
		"refresh_took", time.Since(begin),
		"total_service_count", len(response),
		"accepted_service_count", len(nextgen),
	)

	c.mtx.Lock()
	c.services = nextgen
	c.mtx.Unlock()

	return nil
}

// ServiceIDs currently being monitored by the cache.
// The set can change over time.
func (c *Cache) ServiceIDs() (ids []string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	ids = make([]string, 0, len(c.services))
	for _, s := range c.services {
		ids = append(ids, s.ID)
	}

	sort.Strings(ids) // mostly for tests

	return ids
}

// Metadata returns selected metadata associated with a given service ID.
// If the cache doesn't contain that service ID, found will be false.
func (c *Cache) Metadata(id string) (name string, version int, found bool) {
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

type stringset map[string]struct{}

func newStringSet(initial []string) stringset {
	ss := stringset{}
	for _, s := range initial {
		ss[s] = struct{}{}
	}
	return ss
}

func (ss stringset) empty() bool {
	return len(ss) == 0
}

func (ss stringset) has(s string) bool {
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
