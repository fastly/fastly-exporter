package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"
)

type mockManager struct {
	ids []string
}

func (m *mockManager) update(ids []string) {
	m.ids = ids
}

type fixedResponseClient struct {
	code     int
	response string
}

func (c fixedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c.code)
		fmt.Fprint(w, c.response)
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

type mockRealtimeClient struct {
	response string
	served   uint64
}

func (c *mockRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	// First request immediately returns real data.
	if atomic.AddUint64(&(c.served), 1) == 1 {
		return fixedResponseClient{200, c.response}.Do(req)
	}

	// Subsequent requests block a bit and then return empty JSON.
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case <-time.After(time.Second):
		return fixedResponseClient{200, "{}"}.Do(req)
	}
}

type countingRealtimeClient struct {
	code     int
	response string
	served   uint64
}

func (c *countingRealtimeClient) Do(req *http.Request) (*http.Response, error) {
	atomic.AddUint64(&(c.served), 1)
	return fixedResponseClient{c.code, c.response}.Do(req)
}

type userAgentCapturingClient struct {
	userAgent atomic.Value
}

func (c *userAgentCapturingClient) Do(req *http.Request) (*http.Response, error) {
	c.userAgent.Store(req.Header.Get("User-Agent"))
	return fixedResponseClient{200, "{}"}.Do(req)
}

func within(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() { // 🔥
			return true
		}
	}
	return false
}
