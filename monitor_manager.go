package main

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// monitorManager maintains the goroutines running the monitor function for each
// service ID. It's possible that a service ID is added or removed during
// runtime; for example, if we're configured to monitor all service IDs
// accessible to a token, and the admin for that account adds or removes a
// service. In that case, the monitor manager updates the monitors accordingly.
type monitorManager struct {
	mtx     sync.Mutex
	running map[string]interrupt

	client      httpClient
	token       string
	resolver    nameResolver
	metrics     prometheusMetrics
	postprocess func()
	logger      log.Logger
}

type interrupt struct {
	cancel func()
	done   <-chan struct{}
}

type nameResolver interface {
	resolve(id string) (name string)
}

// newMonitorManager returns an empty, usable monitor manager.
func newMonitorManager(client httpClient, token string, resolver nameResolver, metrics prometheusMetrics, postprocess func(), logger log.Logger) *monitorManager {
	return &monitorManager{
		running: map[string]interrupt{},

		client:      client,
		token:       token,
		resolver:    resolver,
		metrics:     metrics,
		postprocess: postprocess,
		logger:      logger,
	}
}

// update the set of service IDs that the monitor manager should be managing.
// New service IDs spawn new monitors; existing service IDs leave their monitors
// unchanged; service IDs that were in the manager but aren't in the incoming
// set of IDs have their monitors canceled and reaped.
func (m *monitorManager) update(ids []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	nextgen := map[string]interrupt{}
	for _, id := range ids {
		if irq, ok := m.running[id]; ok {
			delete(m.running, id)
			nextgen[id] = irq
		} else {
			level.Info(m.logger).Log("service_id", id, "service_name", m.resolver.resolve(id), "monitor", "start")
			nextgen[id] = m.spawn(id)
		}
	}

	for id, irq := range m.running {
		level.Info(m.logger).Log("service_id", id, "service_name", m.resolver.resolve(id), "monitor", "stop")
		irq.cancel()
		<-irq.done
	}

	m.running = nextgen
}

// spawn a new monitor watching the provided service ID. Return an interrupt,
// which allows the monitor to be canceled and reaped.
func (m *monitorManager) spawn(id string) interrupt {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
	)
	go func() {
		monitor(ctx, m.client, m.token, id, m.resolver, m.metrics, m.postprocess, log.With(m.logger, "service_id", id))
		close(done)
	}()
	return interrupt{cancel, done}
}

// stopAll cancels and reaps all running monitors in sequence.
// When it returns, the monitor manager is empty.
func (m *monitorManager) stopAll() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for id, irq := range m.running {
		level.Info(m.logger).Log("service_id", id, "service_name", m.resolver.resolve(id), "monitor", "stop")
		irq.cancel()
		<-irq.done
		delete(m.running, id)
	}
}

// currentRunning returns all service IDs that are currently being monitored.
func (m *monitorManager) currentlyRunning() (ids []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for id := range m.running {
		ids = append(ids, id)
	}
	return ids
}
