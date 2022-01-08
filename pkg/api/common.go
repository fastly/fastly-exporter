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

// GetNextLink returns the `rel="next"` URI from the `Link` header, if one
// exists. The URI is not validated. See RFC 5899.
func GetNextLink(h http.Header) string {
	for _, link := range h.Values("Link") {
		var (
			linkURI   string
			isRelNext bool
		)
		for _, token := range strings.Split(link, ";") {
			// Ignore empties.
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}

			// Each link should include a URI.
			var (
				isURI   = token[0] == '<' && token[len(token)-1] == '>'
				haveURI = linkURI != ""
				canTake = isURI && !haveURI
			)
			if canTake {
				linkURI = strings.Trim(token, "<>")
			}
			if isURI {
				continue
			}

			// Ensure it's a key/value pair.
			params := strings.SplitN(token, "=", 2)
			if len(params) != 2 {
				continue
			}

			// Only care about key `rel`.
			key := strings.Trim(params[0], ` `)
			if key != "rel" {
				continue
			}

			// Only care about val `"next"``.
			val := strings.Trim(params[1], ` "'`) // TODO(pb): OK?
			if !strings.EqualFold(val, "next") {
				continue
			}

			// The URI we captured previously is the next link.
			isRelNext = true
			break
		}

		if isRelNext && linkURI != "" {
			return linkURI
		}
	}

	return ""
}
