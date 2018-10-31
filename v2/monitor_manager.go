package main

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

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

func (m *monitorManager) currentlyRunning() (ids []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for id := range m.running {
		ids = append(ids, id)
	}
	return ids
}
