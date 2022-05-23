package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAgentTransport(t *testing.T) {
	c := make(chan string, 1)
	handler := func(_ http.ResponseWriter, r *http.Request) { c <- r.Header.Get("User-Agent") }
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	userAgent := "something"
	transport := userAgentTransport(http.DefaultTransport, userAgent)
	client := &http.Client{Transport: transport}
	if _, err := client.Get(server.URL); err != nil {
		t.Fatal(err)
	}

	if want, have := userAgent, <-c; want != have {
		t.Fatalf("want %q, have %q", want, have)
	}
}
