package api

import (
	"net/http"
)

type CustomUserAgent struct {
	agentIdentifier string
	rt              http.RoundTripper
}

// RoundTrip satisfies the http.RoundTripper interface for CustomUserAgent
func (cua *CustomUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", cua.agentIdentifier)
	return cua.rt.RoundTrip(req)
}

// NewCustomUserAgent the provided http.Roundtriper will set
// the User-Agent header for each request.
func NewCustomUserAgent(rt http.RoundTripper, agentIdentifier string) *CustomUserAgent {
	return &CustomUserAgent{
		agentIdentifier,
		rt,
	}
}
