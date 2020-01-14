package rt

import (
	"context"
	"sort"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
)

// ServiceProvider is a consumer contract for a subscriber manager.
// It models the service ID listing and lookup methods of an api.Cache.
type ServiceProvider interface {
	ServiceIDs() []string
}

// Manager owns a set of subscribers. When refreshed, it will ask a
// ServiceProvider for a set of service IDs that should be active, and manage
// the lifecycles of the corresponding subscribers.
type Manager struct {
	provider          ServiceProvider
	client            HTTPClient
	token             string
	metrics           *prom.Metrics
	subscriberOptions []SubscriberOption
	logger            log.Logger

	mtx     sync.RWMutex
	managed map[string]interrupt
}

// NewManager returns a usable manager. Callers should invoke Refresh on a
// regular schedule to keep the set of managed subscribers up-to-date. The HTTP
// client, token, metrics, and subscriber options parameters are passed thru to
// constructed subscribers.
func NewManager(provider ServiceProvider, client HTTPClient, token string, metrics *prom.Metrics, subscriberOptions []SubscriberOption, logger log.Logger) *Manager {
	return &Manager{
		provider:          provider,
		client:            client,
		token:             token,
		metrics:           metrics,
		subscriberOptions: subscriberOptions,
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
	for _, id := range m.provider.ServiceIDs() {
		if irq, ok := m.managed[id]; ok {
			level.Debug(m.logger).Log("service_id", id, "subscriber", "maintain")
			nextgen[id] = irq // move
			delete(m.managed, id)
		} else {
			level.Info(m.logger).Log("service_id", id, "subscriber", "create")
			nextgen[id] = m.spawn(id)
		}
	}

	for _, id := range m.managedIDsWithLock() {
		level.Info(m.logger).Log("service_id", id, "subscriber", "stop")
		irq := m.managed[id]
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		level.Debug(m.logger).Log("service_id", id, "interrupt", err)
	}

	for id, irq := range nextgen {
		select {
		default: // still running (good)
		case err := <-irq.done: // exited (bad)
			level.Error(m.logger).Log("service_id", id, "interrupt", err, "err", "premature termination", "msg", "will attempt to reconnect on next refresh")
			delete(nextgen, id)
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

	for _, id := range m.managedIDsWithLock() {
		level.Info(m.logger).Log("service_id", id, "subscriber", "stop")
		irq := m.managed[id]
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		level.Debug(m.logger).Log("service_id", id, "interrupt", err)
	}
}

func (m *Manager) spawn(serviceID string) interrupt {
	var (
		subscriber  = NewSubscriber(m.client, m.token, serviceID, m.metrics, m.subscriberOptions...)
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan error, 1)
	)
	go func() {
		done <- subscriber.Run(ctx)
	}()
	return interrupt{cancel, done}
}

func (m *Manager) managedIDsWithLock() (ids []string) {
	ids = make([]string, 0, len(m.managed))
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
