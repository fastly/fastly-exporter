package api

import (
	"net/http"
)

type CustomUserAgent struct {
	agentIdentifier string
	rt              http.RoundTripper
}

// DefaultUserAgent passed to HTTPClient when the User-Agent
// HTTP header is not set.
const DefaultUserAgent = "Fastly-Exporter (unknown version)"

// RoundTrip satisfies the http.RoundTripper interface for CustomUserAgent
func (cua *CustomUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	if cua.agentIdentifier == "" {
		cua.agentIdentifier = DefaultUserAgent
	}
	req.Header.Set("User-Agent", cua.agentIdentifier)
	return cua.rt.RoundTrip(req)
}

// NewCustomUserAgent the provided http.Roundtriper will set
// the User-Agent header for each request.
func NewCustomUserAgent(rt http.RoundTripper, agentIdentifier string) http.RoundTripper {
	return &CustomUserAgent{
		agentIdentifier,
		rt,
	}
}
