package api_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
)

func TestCustomUserAgent(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ua   string
		want string
	}{
		{
			name: "no UA provided",
			ua:   "",
			want: api.DefaultUserAgent,
		},
		{
			name: "UA provided",
			ua:   "someclient/v1.2.0",
			want: "someclient/v1.2.0",
		},
	}

	for _, testcase := range tt {
		t.Run(testcase.name, func(t *testing.T) {
			mockres := io.NopCloser(bytes.NewBuffer([]byte(`{}`)))
			transporter := newMockTransport(mockres)
			uaTransporter := api.NewCustomUserAgent(transporter, testcase.ua)
			c := http.Client{
				Transport: uaTransporter,
			}

			res, err := c.Get("/")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if want, have := testcase.want, res.Request.Header.Get("User-Agent"); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}
		})
	}
}

type uaTransport struct {
	body io.ReadCloser
}

func (t *uaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp http.Response
	{
		resp = http.Response{
			Body:    t.body,
			Request: req,
		}
	}

	return &resp, nil
}

func newMockTransport(body io.ReadCloser) http.RoundTripper {
	return &uaTransport{body: body}
}
