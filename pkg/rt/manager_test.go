package rt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
	"github.com/peterbourgon/fastly-exporter/pkg/rt"
	"github.com/prometheus/client_golang/prometheus"
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
		manager  = rt.NewManager(cache, client, token, registry, options, level.NewFilter(logger, level.AllowInfo()))
	)

	assertStringSliceEqual(t, []string{}, manager.Active())

	cache.update([]api.Service{s1, s2})
	manager.Refresh() // create s1, create s2
	assertStringSliceEqual(t, []string{s1.ID, s2.ID}, manager.Active())

	cache.update([]api.Service{s3, s1, s2})
	manager.Refresh() // create s3
	assertStringSliceEqual(t, []string{s1.ID, s2.ID, s3.ID}, manager.Active())

	manager.Refresh() // no effect
	assertStringSliceEqual(t, []string{s1.ID, s2.ID, s3.ID}, manager.Active())

	cache.update([]api.Service{s3, s2})
	manager.Refresh() // stop s1
	assertStringSliceEqual(t, []string{s2.ID, s3.ID}, manager.Active())

	cache.update([]api.Service{s2})
	manager.Refresh() // stop s3
	assertStringSliceEqual(t, []string{s2.ID}, manager.Active())

	cache.update([]api.Service{})
	manager.Refresh() // stop s2
	assertStringSliceEqual(t, []string{}, manager.Active())

	cache.update([]api.Service{s2, s3})
	manager.Refresh() // create s2, create s3
	assertStringSliceEqual(t, []string{s2.ID, s3.ID}, manager.Active())

	manager.StopAll() // stop s2, stop s3
	assertStringSliceEqual(t, []string{}, manager.Active())

	if want, have := []string{
		`level=info service_id=101010 subscriber=create`,
		`level=info service_id=2f2f2f subscriber=create`,
		`level=info service_id=3a3b3c subscriber=create`,
		`level=info service_id=101010 subscriber=stop`,
		`level=info service_id=3a3b3c subscriber=stop`,
		`level=info service_id=2f2f2f subscriber=stop`,
		`level=info service_id=2f2f2f subscriber=create`,
		`level=info service_id=3a3b3c subscriber=create`,
		`level=info service_id=2f2f2f subscriber=stop`,
		`level=info service_id=3a3b3c subscriber=stop`,
	}, strings.Split(strings.TrimSpace(logbuf.String()), "\n"); !cmp.Equal(want, have) {
		t.Error(cmp.Diff(want, have))
	}
}

//
//
//

type nopCollector struct{}

func (nopCollector) Describe(ch chan<- *prometheus.Desc) {}
func (nopCollector) Collect(ch chan<- prometheus.Metric) {}
