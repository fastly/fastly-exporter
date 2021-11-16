package api

import "net/http"

// HTTPClient is a consumer contract for components in this package.
// It models a concrete http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
