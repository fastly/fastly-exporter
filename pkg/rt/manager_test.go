package rt_test

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/fastly/fastly-exporter/pkg/prom"
	"github.com/fastly/fastly-exporter/pkg/rt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/go-cmp/cmp"
)

func TestManager(t *testing.T) {
	var (
		cache    = &mockCache{}
		s1       = api.Service{ID: "101010", Name: "service 1", Version: 1}
		s2       = api.Service{ID: "2f2f2f", Name: "service 2", Version: 2}
		s3       = api.Service{ID: "3a3b3c", Name: "service 3", Version: 3}
		client   = newMockRealtimeClient(`{}`)
		token    = "irrelevant-token"
		registry = prom.NewRegistry("v0.0.0-DEV", "namespace", "subsystem", filter.Filter{})
		logbuf   = &bytes.Buffer{}
		logger   = log.NewLogfmtLogger(logbuf)
		options  = []rt.SubscriberOption{rt.WithMetadataProvider(cache)}
		products = newMockProductCache()
		manager  = rt.NewManager(cache, client, token, registry, options, products, level.NewFilter(logger, level.AllowInfo()))
	)

	assertStringSliceEqual(t, []string{}, sortedServiceIDs(manager))

	products.update(api.ProductOriginInspector, false)
	products.update(api.ProductDomainInspector, false)

	cache.update([]api.Service{s1, s2})
	manager.Refresh() // create s1, create s2
	assertStringSliceEqual(t, []string{s1.ID, s2.ID}, sortedServiceIDs(manager))

	cache.update([]api.Service{s3, s1, s2})
	manager.Refresh() // create s3
	assertStringSliceEqual(t, []string{s1.ID, s2.ID, s3.ID}, sortedServiceIDs(manager))

	manager.Refresh() // no effect
	assertStringSliceEqual(t, []string{s1.ID, s2.ID, s3.ID}, sortedServiceIDs(manager))

	cache.update([]api.Service{s3, s2})
	manager.Refresh() // stop s1
	assertStringSliceEqual(t, []string{s2.ID, s3.ID}, sortedServiceIDs(manager))

	cache.update([]api.Service{s2})
	manager.Refresh() // stop s3
	assertStringSliceEqual(t, []string{s2.ID}, sortedServiceIDs(manager))

	cache.update([]api.Service{})
	manager.Refresh() // stop s2
	assertStringSliceEqual(t, []string{}, sortedServiceIDs(manager))

	cache.update([]api.Service{s2, s3})
	manager.Refresh() // create s2, create s3
	assertStringSliceEqual(t, []string{s2.ID, s3.ID}, sortedServiceIDs(manager))

	manager.StopAll() // stop s2, stop s3
	assertStringSliceEqual(t, []string{}, sortedServiceIDs(manager))

	products.update(api.ProductOriginInspector, true)
	cache.update([]api.Service{s1})
	manager.Refresh() // create s1 with origin inspector
	// expecting the ID twice -- one for each product
	assertStringSliceEqual(t, []string{s1.ID, s1.ID}, sortedServiceIDs(manager))

	products.update(api.ProductDomainInspector, true)
	cache.update([]api.Service{s1})
	manager.Refresh() // create s1 with domain inspector
	// expecting the ID thrice -- one for each product
	assertStringSliceEqual(t, []string{s1.ID, s1.ID, s1.ID}, sortedServiceIDs(manager))

	manager.StopAll() // stop s1
	assertStringSliceEqual(t, []string{}, sortedServiceIDs(manager))

	want := []string{
		`level=info service_id=101010 type=default subscriber=create`,
		`level=info service_id=2f2f2f type=default subscriber=create`,
		`level=info service_id=3a3b3c type=default subscriber=create`,
		`level=info service_id=101010 type=default subscriber=stop`,
		`level=info service_id=3a3b3c type=default subscriber=stop`,
		`level=info service_id=2f2f2f type=default subscriber=stop`,
		`level=info service_id=2f2f2f type=default subscriber=create`,
		`level=info service_id=3a3b3c type=default subscriber=create`,
		`level=info service_id=2f2f2f type=default subscriber=stop`,
		`level=info service_id=3a3b3c type=default subscriber=stop`,
		`level=info service_id=101010 type=default subscriber=create`,
		`level=info service_id=101010 type=origin_inspector subscriber=create`,
		`level=info service_id=101010 type=domain_inspector subscriber=create`,
		`level=info service_id=101010 type=default subscriber=stop`,
		`level=info service_id=101010 type=origin_inspector subscriber=stop`,
		`level=info service_id=101010 type=domain_inspector subscriber=stop`,
	}
	have := strings.Split(strings.TrimSpace(logbuf.String()), "\n")
	sort.Strings(want)
	sort.Strings(have)

	if !cmp.Equal(want, have) {
		t.Error(cmp.Diff(want, have))
	}
}

func sortedServiceIDs(m *rt.Manager) []string {
	serviceIDs := m.Active()
	sort.Strings(serviceIDs)
	return serviceIDs
}
