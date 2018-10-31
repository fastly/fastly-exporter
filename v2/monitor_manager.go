package main

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type monitorManager struct {
	mtx      sync.Mutex
	token    string
	running  map[string]interrupt
	resolver nameResolver
	metrics  prometheusMetrics
	logger   log.Logger
}

type interrupt struct {
	cancel func()
	done   <-chan struct{}
}

type nameResolver interface {
	resolve(id string) (name string)
}

func newMonitorManager(token string, resolver nameResolver, metrics prometheusMetrics, logger log.Logger) *monitorManager {
	return &monitorManager{
		token:    token,
		running:  map[string]interrupt{},
		resolver: resolver,
		metrics:  metrics,
		logger:   logger,
	}
}

func (m *monitorManager) update(ids []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	nextgen := map[string]interrupt{}
	for _, id := range ids {
		if irq, found := m.running[id]; found {
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
}

func (m *monitorManager) spawn(id string) interrupt {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		done        = make(chan struct{})
	)
	go func() {
		monitor(ctx, m.token, id, m.resolver, m.metrics, log.With(m.logger, "service_id", id))
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
