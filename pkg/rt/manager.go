package rt

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/fastly/fastly-exporter/pkg/prom"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ServiceIdentifier is a consumer contract for a subscriber manager.
// It models the service ID listing and lookup methods of an api.Cache.
type ServiceIdentifier interface {
	ServiceIDs() []string
}

// MetricsProvider is a consumer contract for a subscriber manager. It models
// the method of the prom.Registry which yields a set of Prometheus metrics for
// a specific service ID.
type MetricsProvider interface {
	MetricsFor(serviceID string) *prom.Metrics
}

// Manager owns a set of subscribers. On refresh, it asks the ServiceIdentifier
// for a set of service IDs that should be active, and manages the lifecycles of
// the corresponding subscribers.
type Manager struct {
	ids               ServiceIdentifier
	client            HTTPClient
	token             string
	metrics           MetricsProvider
	subscriberOptions []SubscriberOption
	productCache      *api.ProductCache
	logger            log.Logger

	mtx     sync.RWMutex
	managed map[string]interrupt
}

// NewManager returns a usable manager. Callers should invoke Refresh on a
// regular schedule to keep the set of managed subscribers up-to-date. The HTTP
// client, token, metrics, and subscriber options parameters are passed thru to
// constructed subscribers.
func NewManager(ids ServiceIdentifier, client HTTPClient, token string, metrics MetricsProvider, subscriberOptions []SubscriberOption, productCache *api.ProductCache, logger log.Logger) *Manager {
	return &Manager{
		ids:               ids,
		client:            client,
		token:             token,
		metrics:           metrics,
		subscriberOptions: subscriberOptions,
		productCache:      productCache,
		logger:            logger,

		managed: map[string]interrupt{},
	}
}

// Refresh the set of subscribers managed by the manager, by asking the
// authority provided in the constructor for the current set of service IDs, and
// comparing those IDs with the ones already under management. If a service ID
// was not previously managed, start a new subscriber. If a service ID was
// previously managed but isn't in the latest set of IDs, terminate the
// subscriber. Finally, if a service ID was both previously managed and is in
// the latest set of IDs, simply keep the existing subscriber.
func (m *Manager) Refresh() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	nextgen := map[string]interrupt{}
	for _, product := range api.Products {
		if m.productCache.HasAccess(product) {
			for _, id := range m.ids.ServiceIDs() {
				key := id + sep + product

				if irq, ok := m.managed[key]; ok {
					level.Debug(m.logger).Log("service_id", id, "type", product, "subscriber", "maintain")
					nextgen[key] = irq // move
					delete(m.managed, key)
				} else {
					level.Info(m.logger).Log("service_id", id, "type", product, "subscriber", "create")
					nextgen[key] = m.spawn(id, product)
				}
			}
		}

		for _, key := range m.managedIDsWithLock() {
			id, managedProduct := parseKey(key)
			if managedProduct != product {
				continue
			}

			level.Info(m.logger).Log("service_id", id, "type", product, "subscriber", "stop")
			irq := m.managed[key]
			irq.cancel()
			err := <-irq.done
			delete(m.managed, key)
			level.Debug(m.logger).Log("service_id", id, "type", product, "interrupt", err)
		}

		for key, irq := range nextgen {
			id, managedProduct := parseKey(key)
			if managedProduct != product {
				continue
			}

			select {
			default: // still running (good)
			case err := <-irq.done: // exited (bad)
				level.Error(m.logger).Log("service_id", id, "type", product, "interrupt", err, "err", "premature termination", "msg", "will attempt to reconnect on next refresh")
				delete(nextgen, key)
			}
		}
	}

	m.managed = nextgen
}

// Active returns the set of service IDs currently being managed.
// Mostly useful for tests.
func (m *Manager) Active() (serviceIDs []string) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.managedIDsWithLock()
}

// StopAll terminates and cleans up all active subscribers.
func (m *Manager) StopAll() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, key := range m.managedIDsWithLock() {
		id, product := parseKey(key)
		level.Info(m.logger).Log("service_id", id, "type", product, "subscriber", "stop")
		irq := m.managed[key]
		irq.cancel()
		for i := 0; i < cap(irq.done); i++ {
			err := <-irq.done
			level.Debug(m.logger).Log("service_id", id, "goroutine", i+1, "of", cap(irq.done), "interrupt", err)
		}
		delete(m.managed, key)
	}
}

func (m *Manager) spawn(serviceID string, product string) interrupt {
	var (
		subscriber  = NewSubscriber(m.client, m.token, serviceID, m.metrics.MetricsFor(serviceID), m.subscriberOptions...)
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan error, 1)
	)
	switch product {
	case api.OriginInspector:
		go func() { done <- fmt.Errorf("origins: %w", subscriber.RunOrigins(ctx)) }()
	default:
		go func() { done <- fmt.Errorf("realtime: %w", subscriber.RunRealtime(ctx)) }()
	}

	return interrupt{cancel, done}
}

func (m *Manager) managedIDsWithLock() []string {
	ids := make([]string, 0, len(m.managed))
	for id := range m.managed {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

type interrupt struct {
	cancel func()
	done   <-chan error
}

func parseKey(key string) (string, string) {
	ks := strings.Split(key, sep)
	return ks[0], ks[1]
}

const sep = "|"
