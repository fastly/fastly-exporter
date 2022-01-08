package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

//
//
//

type paginatedResponseClient struct {
	responses []string
}

func (c paginatedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageIndex := page - 1
		if pageIndex >= len(c.responses) {
			http.Error(w, "page too large", 400)
			return
		}

		if pageIndex+1 < len(c.responses) {
			values := r.URL.Query()
			values.Set("page", strconv.Itoa(page+1))
			r.URL.RawQuery = values.Encode()
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, r.URL.String()))
		}

		fmt.Fprint(w, c.responses[pageIndex])
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

//
//
//

func TestGetNextLink(t *testing.T) {
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
		{
			name:  `link 1`,
			input: `<https://example.com/?page=2>; rel="next"; title="foo"`,
			want:  `https://example.com/?page=2`,
		},
		{
			name:  `link 2`,
			input: `<https://example.com/?page=2>; rel="next"`,
			want:  `https://example.com/?page=2`,
		},
		{
			name:  `link 3`,
			input: `<https://example.com/?page=2>; rel="next",<https://example.com/?page=34>; rel="last"`,
			want:  `https://example.com/?page=2`,
		},
		{
			name:  `link 4`,
			input: `<//www.w3.org/wiki/LinkHeader>; rel="original latest-version",<//www.w3.org/wiki/Special:TimeGate/LinkHeader>; rel="timegate",<//www.w3.org/wiki/Special:TimeMap/LinkHeader>; rel="timemap"; type="application/link-format"; from="Mon, 03 Sep 2007 14:52:48 GMT"; until="Tue, 16 Jun 2015 22:59:23 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=84697>; rel="last memento"; datetime="Tue, 16 Jun 2015 22:59:23 GMT"`,
			want:  `//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			h := http.Header{"Link": strings.Split(testcase.input, ",")}
			if want, have := testcase.want, api.GetNextLink(h); want != have {
				t.Fatalf("want %q, have %q", want, have)
			}
		})
	}
}
