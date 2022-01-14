package rt

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/peterbourgon/fastly-exporter/pkg/gen"
)

// HTTPClient is a consumer contract for the subscriber.
// It models a concrete http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ServiceProvider is a consumer contract for a subscriber manager. It models
// the method of the api.ServiceCache that returns service IDs to export.
type ServiceProvider interface{ ServiceIDs() []string }

// MetricsProvider is a consumer contract for a subscriber manager. It models
// the method of the prom.Registry which yields a set of Prometheus metrics for
// a specific service ID.
type MetricsProvider interface{ MetricsFor(id string) *gen.Metrics }

// MetadataProvider is a consumer contract for the manager and the subscribers.
// It models the service lookup method of an api.ServiceCache.
type MetadataProvider interface {
	Metadata(id string) (name string, version int, found bool)
}

// ManagerConfig collects the parameters used to construct a Manager.
type ManagerConfig struct {
	// Client used by managed subscribers to query rt.fastly.com. If not
	// provided, http.DefaultClient is used, which may not include the desired
	// User-Agent, among other things.
	Client HTTPClient

	// Token provided as the Fastly-Key when querying rt.fastly.com.
	Token string

	// Services yields the current set of service IDs which should be managed.
	// If not provided, the manager wouldn't ever start any subscribers, and no
	// stats would be fetched, which kind of defeats the purpose? So, required.
	Services ServiceProvider

	// Metrics yields per-service gen.Metrics structs for subscribers to update.
	// If not provided, all of the data we fetch from the real-time stats API
	// would get dumped to /dev/null, rendering the program kind of pointless.
	// So, required.
	Metrics MetricsProvider

	// Metadata yields per-service metadata like service name, which ends up in
	// Prometheus labels. If not provided, relevant labels will have "unknown"
	// or zero values.
	Metadata MetadataProvider

	// Logger is used for runtime diagnostic information.
	// If not provided, a no-op logger is used.
	Logger log.Logger
}

func (c *ManagerConfig) validate() error {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	if c.Services == nil {
		return fmt.Errorf("manager: a source of service IDs is required")
	}

	if c.Metrics == nil {
		return fmt.Errorf("manager: a source of Prometheus metrics is required")
	}

	if c.Metadata == nil {
		c.Metadata = nopMetadataProvider{}
	}

	if c.Logger == nil {
		c.Logger = nopLogger
	}

	return nil
}

// Manager owns a set of subscribers. When refreshed, it will ask a
// ServiceIdentifier for a set of service IDs that should be active, and manage
// the lifecycles of the corresponding subscribers.
type Manager struct {
	ManagerConfig

	mtx     sync.RWMutex
	managed map[string]interrupt
}

// NewManager returns a usable manager. Callers should invoke Refresh on a
// regular schedule to keep the set of managed subscribers up-to-date. The HTTP
// client, token, metrics, and subscriber options parameters are passed thru to
// constructed subscribers.
func NewManager(c ManagerConfig) (*Manager, error) {
	return &Manager{
		ManagerConfig: c,
	}, c.validate()
}

// Refresh the manager by fetching the current set of service IDs from the
// metadata provider. If a service ID was not previously managed, start a new
// subscriber. If a service ID was previously managed but isn't in the latest
// set of IDs, terminate the subscriber. Finally, if a service ID was both
// previously managed and is in the latest set of IDs, simply keep the existing
// subscriber.
func (m *Manager) Refresh() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	nextgen := map[string]interrupt{}

	for _, id := range m.Services.ServiceIDs() {
		if irq, found := m.managed[id]; found {
			level.Debug(m.Logger).Log("service_id", id, "subscriber", "maintain")
			nextgen[id] = irq // move
			delete(m.managed, id)
		} else {
			level.Info(m.Logger).Log("service_id", id, "subscriber", "create")
			nextgen[id] = m.spawn(id)
		}
	}

	for id, irq := range m.managed {
		level.Info(m.Logger).Log("service_id", id, "subscriber", "stop")
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		level.Debug(m.Logger).Log("service_id", id, "interrupt", err)
	}

	for id, irq := range nextgen {
		select {
		default: // still running (good)
		case err := <-irq.done: // exited (bad)
			level.Error(m.Logger).Log("service_id", id, "interrupt", err, "msg", "premature termination, will attempt to reconnect on next refresh")
			delete(nextgen, id)
		}
	}

	m.managed = nextgen
}

// Active service IDs being managed. Mostly useful for tests.
func (m *Manager) Active() []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	ids := make([]string, 0, len(m.managed))
	for id := range m.managed {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Close terminates all subscribers.
func (m *Manager) Close() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for id, irq := range m.managed {
		level.Info(m.Logger).Log("service_id", id, "subscriber", "stop")
		irq.cancel()
		err := <-irq.done
		delete(m.managed, id)
		level.Debug(m.Logger).Log("service_id", id, "interrupt", err)
	}
}

func (m *Manager) spawn(serviceID string) interrupt {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		s, err := NewSubscriber(SubscriberConfig{
			Client:    m.Client,
			Token:     m.Token,
			ServiceID: serviceID,
			Metrics:   m.Metrics.MetricsFor(serviceID),
			Metadata:  m.Metadata,
			Logger:    m.Logger,
		})
		if err == nil {
			done <- s.Run(ctx)
		} else {
			done <- err
		}
	}()
	return interrupt{cancel, done}
}

type interrupt struct {
	cancel func()
	done   <-chan error
}
