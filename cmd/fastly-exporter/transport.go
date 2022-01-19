package main

import "net/http"

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func userAgentTransport(next http.RoundTripper, userAgent string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.Header.Set("User-Agent", userAgent)
		return next.RoundTrip(req)
	})
}
