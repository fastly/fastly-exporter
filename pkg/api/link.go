package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetNextLink reads the `Link` response headers, and returns the first URI
// tagged as `rel="next"`. Relative URIs are evaluated against the original
// request if possible. See RFC 5899.
func GetNextLink(resp *http.Response) (*url.URL, error) {
	rawuri, ok := uriFromLinks(resp.Header.Values("Link"), `rel`, `next`)
	if !ok {
		return nil, fmt.Errorf(`rel="next": no match`)
	}

	if resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL.Parse(rawuri)
	}

	return url.Parse(rawuri)
}

func uriFromLinks(links []string, k, v string) (string, bool) {
	for _, link := range links {
		if rawuri, ok := uriFromLink(link, k, v); ok {
			return rawuri, true
		}
	}
	return "", false
}

func uriFromLink(link, k, v string) (string, bool) {
	var rawuri string
	for _, link := range strings.Split(link, ",") {
		for _, param := range strings.Split(link, ";") {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}

			if param[0] == '<' && param[len(param)-1] == '>' {
				rawuri = strings.Trim(param, "<>")
				continue
			}

			keyval := strings.SplitN(param, "=", 2)
			if len(keyval) != 2 {
				continue
			}

			key := strings.TrimSpace(keyval[0])
			if !strings.EqualFold(key, k) {
				continue
			}

			var val string
			val = keyval[1]
			val = strings.TrimSpace(val)
			val = strings.Trim(val, `"`)
			val = strings.TrimSpace(val)
			if !strings.EqualFold(val, v) {
				continue
			}

			if rawuri != "" {
				return rawuri, true
			}
		}
	}
	return "", false
}
