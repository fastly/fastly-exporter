package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/peterbourgon/fastly-exporter/pkg/api"
)

type fixedResponseClient struct {
	code     int
	response string
}

func (c fixedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c.code)
		fmt.Fprint(w, c.response)
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

func TestGetRelNext(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name: `RFC 5988 1`,
			input: `   Link: <http://example.com/TheBook/chapter2>; rel="previous";			title="previous chapter"`,
			want: ``,
		},
		{
			name:  `RFC 5988 2`,
			input: `   Link: </>; rel="http://example.net/foo"`,
			want:  ``,
		},
		{
			name: `RFC 5988 3`,
			input: `
				Link: </TheBook/chapter2>;
					  rel="previous"; title*=UTF-8'de'letztes%20Kapitel,
					  </TheBook/chapter4>;
					  rel="next"; title*=UTF-8'de'n%c3%a4chstes%20Kapitel
			`,
			want: `/TheBook/chapter4`,
		},
		{
			name: `linkheader 1`,
			input: "<https://api.github.com/user/9287/repos?page=3&per_page=100>; rel=\"next\", " +
				"<https://api.github.com/user/9287/repos?page=1&per_page=100>; rel=\"prev\"; pet=\"cat\", " +
				"<https://api.github.com/user/9287/repos?page=5&per_page=100>; rel=\"last\"",
			want: `https://api.github.com/user/9287/repos?page=3&per_page=100`,
		},
		{
			name: `linkheader 2`,
			input: "<https://api.github.com/user/9287/repos?page=3&per_page=100>; rel=\"next\", " +
				"<https://api.github.com/user/9287/repos?page=1&per_page=100>; rel=\"stylesheet\"; pet=\"cat\", " +
				"<https://api.github.com/user/9287/repos?page=5&per_page=100>; rel=\"stylesheet\"",
			want: `https://api.github.com/user/9287/repos?page=3&per_page=100`,
		},
		{
			name:  `linkheader 3`,
			input: "<https://api.github.com/user/58276/repos?page=9>; rel=\"last\", <https://api.github.com/user/58276/repos?page=2>; rel=\"next\" ",
			want:  `https://api.github.com/user/58276/repos?page=2`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			h := http.Header{"Link": strings.Split(testcase.input, ",")}
			have, _ := api.GetNextLink(h)
			if want, have := testcase.want, have; want != have {
				t.Fatalf("want %q, have %q", want, have)
			}
		})
	}
}
