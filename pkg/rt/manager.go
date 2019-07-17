package rt

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/peterbourgon/fastly-exporter/pkg/prom"
)

// Servicer is a consumer contract for a subscriber manager.
// It models the service ID listing and lookup methods of an api.Cache.
type Servicer interface {
	ServiceIDs() []string
	Metadata(id string) (name string, version int, found bool)
}

// Manager owns a set of subscribers, polling a service authority to determine
// which service IDs should be active.
type Manager struct {
	servicer          Servicer
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
func NewManager(servicer Servicer, client HTTPClient, token string, metrics *prom.Metrics, subscriberOptions []SubscriberOption, logger log.Logger) *Manager {
	return &Manager{
		servicer:          servicer,
		client:            client,
		token:             token,
		metrics:           metrics,
		subscriberOptions: subscriberOptions,
		logger:            logger,

		managed: map[string]interrupt{},
	}
}

// Refresh the set of subscribers managed by the manager, by asking the servicer
// provided in the constructor for the authoritative set of service IDs, and
// comparing those IDs with the ones already under management. If a new service
// ID was not previously managed, start a new subscriber. If a service ID was
// previously managed but isn't in the latest set of IDs, terminate the
// subscriber. Finally, if a service ID was both previously managed and is in
// the latest set of IDs, simply keep the existing subscriber.
func (m *Manager) Refresh() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var (
		ids     = m.servicer.ServiceIDs()
		nextgen = map[string]interrupt{}
	)
	for _, id := range ids {
		name, _, _ := m.servicer.Metadata(id)
		if irq, ok := m.managed[id]; ok {
			level.Debug(m.logger).Log("service_id", id, "service_name", name, "subscriber", "maintain")
			nextgen[id] = irq // move
			delete(m.managed, id)
		} else {
			level.Info(m.logger).Log("service_id", id, "service_name", name, "subscriber", "create")
			nextgen[id] = m.spawn(id)
		}
	}
	for id, irq := range m.managed {
		name, _, _ := m.servicer.Metadata(id)
		level.Info(m.logger).Log("service_id", id, "service_name", name, "subscriber", "stop")
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		level.Debug(m.logger).Log("service_id", id, "service_name", name, "interrupt", err)
	}

	for id, irq := range nextgen {
		select {
		default: // good
		case err := <-irq.done: // bad
			name, _, _ := m.servicer.Metadata(id)
			level.Error(m.logger).Log("service_id", id, "service_name", name, "interrupt", err, "err", "premature termination", "msg", "will attempt reconnect on next refresh")
			delete(nextgen, id) // should get restarted on next refresh
		}
	}

	m.managed = nextgen
}

// Active returns the set of service IDs currently being managed.
func (m *Manager) Active() (serviceIDs []string) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	serviceIDs = make([]string, 0, len(m.managed))
	for id := range m.managed {
		serviceIDs = append(serviceIDs, id)
	}
	return serviceIDs
}

// StopAll terminates and cleans up all active subscribers.
func (m *Manager) StopAll() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for id, irq := range m.managed {
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		name, _, _ := m.servicer.Metadata(id)
		level.Debug(m.logger).Log("service_id", id, "service_name", name, "interrupt", err)
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

type interrupt struct {
	cancel func()
	done   <-chan error
}
