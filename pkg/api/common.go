package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
	sb.WriteString("api.fastly.com responded with")
	sb.WriteString(http.StatusText(e.Code))
	if e.Msg != "" {
		sb.WriteString(" (" + e.Msg + ")")
	}
	return sb.String()
}
