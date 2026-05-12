package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// DefaultBaseDomain is the default Fastly base domain.
// API and RT URLs are derived as https://api.{domain} and https://rt.{domain}.
const DefaultBaseDomain = "fastly.com"

// BaseURL returns the full API base URL for a given base domain.
func BaseURL(baseDomain string) string {
	if baseDomain == "" {
		baseDomain = DefaultBaseDomain
	}
	return "https://api." + baseDomain
}

// RTBaseURL returns the full real-time stats base URL for a given base domain.
func RTBaseURL(baseDomain string) string {
	if baseDomain == "" {
		baseDomain = DefaultBaseDomain
	}
	return "https://rt." + baseDomain
}

// HTTPClient is a consumer contract for components in this package.
// It models a concrete http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Error represents an error received from api.fastly.com.
type Error struct {
	Code int
	Msg  string `json:"msg"`
}

// NewError returns an error derived from the provided response.
func NewError(resp *http.Response) *Error {
	e := &Error{Code: resp.StatusCode}
	json.NewDecoder(resp.Body).Decode(e)
	return e
}

// Error implements the error interface.
func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString("api.fastly.com responded with ")
	sb.WriteString(http.StatusText(e.Code))
	if e.Msg != "" {
		sb.WriteString(" (" + e.Msg + ")")
	}
	return sb.String()
}
