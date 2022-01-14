package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
)

var nopLogger = log.NewNopLogger()

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
	s := fmt.Sprintf("api.fastly.com responded with %d %s", e.Code, http.StatusText(e.Code))
	if e.Msg != "" {
		s += " (" + e.Msg + ")"
	}
	return s
}
