package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func TestMonitorManager(t *testing.T) {
	var (
		client      = fixedResponseClient{"{}"}
		token       = "irrelevant-token"
		cache       = newNameCache()
		metrics     = prometheusMetrics{}
		postprocess = func() {}
		logbuf      = bytes.Buffer{}
		logger      = level.NewFilter(log.NewLogfmtLogger(&logbuf), level.AllowInfo())
		manager     = newMonitorManager(client, token, cache, metrics, postprocess, logger)
	)

	manager.update([]string{"c", "b", "a"})

	if want, have := []string{"a", "b", "c"}, manager.currentlyRunning(); !setEqual(want, have) {
		t.Errorf("first gen: running: want %v, have %v", want, have)
	}
	if want, have := []string{
		"level=info service_id=a service_name=a monitor=start",
		"level=info service_id=b service_name=b monitor=start",
		"level=info service_id=c service_name=c monitor=start",
	}, strings.Split(strings.TrimSpace(logbuf.String()), "\n"); !setEqual(want, have) {
		t.Errorf("first gen: logs: want\n%s\nhave\n%s\n", strings.Join(want, "\n"), strings.Join(have, "\n"))
	}
	logbuf.Reset()

	manager.update([]string{"c", "d", "f", "e"})

	if want, have := []string{"c", "d", "e", "f"}, manager.currentlyRunning(); !setEqual(want, have) {
		t.Errorf("second gen: running: want %v, have %v", want, have)
	}
	if want, have := []string{
		"level=info service_id=a service_name=a monitor=stop",
		"level=info service_id=b service_name=b monitor=stop",
		// c should be untouched
		"level=info service_id=d service_name=d monitor=start",
		"level=info service_id=e service_name=e monitor=start",
		"level=info service_id=f service_name=f monitor=start",
	}, strings.Split(strings.TrimSpace(logbuf.String()), "\n"); !setEqual(want, have) {
		t.Errorf("second gen: logs: want\n%s\nhave\n%s\n", strings.Join(want, "\n"), strings.Join(have, "\n"))
	}
	logbuf.Reset()

	manager.stopAll()

	if want, have := []string{}, manager.currentlyRunning(); !setEqual(want, have) {
		t.Errorf("stopAll: running: want %v, have %v", want, have)
	}
	if want, have := []string{
		"level=info service_id=c service_name=c monitor=stop",
		"level=info service_id=d service_name=d monitor=stop",
		"level=info service_id=e service_name=e monitor=stop",
		"level=info service_id=f service_name=f monitor=stop",
	}, strings.Split(strings.TrimSpace(logbuf.String()), "\n"); !setEqual(want, have) {
		t.Errorf("stopAll: logs: want\n%s\nhave\n%s\n", strings.Join(want, "\n"), strings.Join(have, "\n"))
	}
}

func setEqual(a, b []string) bool {
	m := map[string]bool{}
	for _, s := range a {
		m[s] = true
	}
	for _, s := range b {
		if _, ok := m[s]; !ok {
			return false
		}
		delete(m, s)
	}
	return len(m) == 0
}
