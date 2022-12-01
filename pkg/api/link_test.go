package api_test

import (
	"net/http"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
)

func TestGetNextLink(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name  string
		req   string
		links []string
		want  string
	}{
		{
			name:  "basic",
			req:   "http://zombo.com",
			links: nil,
			want:  "",
		},
		{
			name:  "full",
			req:   "http://zombo.com",
			links: []string{`<http://zombo.com/2> ; rel="next" `},
			want:  "http://zombo.com/2",
		},

		{
			name:  `RFC 5988 1`,
			links: []string{` <http://example.com/TheBook/chapter2>; rel="previous";			title="previous chapter"`},
			want:  ``,
		},
		{
			name:  `RFC 5988 2`,
			links: []string{`</>; rel="http://example.net/foo"`},
			want:  ``,
		},
		{
			name: `RFC 5988 3`,
			links: []string{
				`</TheBook/chapter2>;
					  rel="previous"; title*=UTF-8'de'letztes%20Kapitel`,
				`</TheBook/chapter4>;
					  rel="next"; title*=UTF-8'de'n%c3%a4chstes%20Kapitel`,
			},
			want: `/TheBook/chapter4`,
		},
		{
			name: `linkheader 1`,
			links: []string{
				`<https://api.github.com/user/9287/repos?page=3&per_page=100>; rel="next"`,
				`<https://api.github.com/user/9287/repos?page=1&per_page=100>; rel="prev"; pet="cat"`,
				`<https://api.github.com/user/9287/repos?page=5&per_page=100>; rel="last""`,
			},
			want: `https://api.github.com/user/9287/repos?page=3&per_page=100`,
		},
		{
			name: `linkheader 2`,
			links: []string{
				`<https://api.github.com/user/9287/repos?page=1&per_page=100>; rel="stylesheet"; pet="cat"`,
				`<https://api.github.com/user/9287/repos?page=5&per_page=100>; rel="stylesheet"`,
				`<https://api.github.com/user/9287/repos?page=3&per_page=100>; rel="next"`,
			},
			want: `https://api.github.com/user/9287/repos?page=3&per_page=100`,
		},
		{
			name:  `linkheader 3`,
			links: []string{"<https://api.github.com/user/58276/repos?page=9>; rel=\"last\"", "<https://api.github.com/user/58276/repos?page=2>; rel=\"next\" "},
			want:  `https://api.github.com/user/58276/repos?page=2`,
		},
		{
			name:  `linkheader 4`,
			links: []string{"<https://api.github.com/user/58276/repos?page=9>; rel=\"last\", <https://api.github.com/user/58276/repos?page=2>; rel=\"next\""},
			want:  `https://api.github.com/user/58276/repos?page=2`,
		},
		{
			name:  `link 1`,
			links: []string{`<https://example.com/?page=2>; rel="next"; title="foo"`},
			want:  `https://example.com/?page=2`,
		},
		{
			name:  `link 2`,
			links: []string{`<https://example.com/?page=2>; rel="next"`},
			want:  `https://example.com/?page=2`,
		},
		{
			name: `link 3`,
			links: []string{
				`<https://example.com/?page=2>; rel="next"`,
				`<https://example.com/?page=34>; rel="last"`,
			},
			want: `https://example.com/?page=2`,
		},
		{
			name: `link 4`,
			links: []string{
				`<//www.w3.org/wiki/LinkHeader>; rel="original latest-version"`,
				`<//www.w3.org/wiki/Special:TimeGate/LinkHeader>; rel="timegate"`,
				`<//www.w3.org/wiki/Special:TimeMap/LinkHeader>; rel="timemap"; type="application/link-format"; from="Mon, 03 Sep 2007 14:52:48 GMT"; until="Tue, 16 Jun 2015 22:59:23 GMT"`,
				`<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`,
				`<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=84697>; rel="last memento"; datetime="Tue, 16 Jun 2015 22:59:23 GMT"`,
			},
			want: `//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152`,
		},
		{
			name:  "whitespace 1",
			links: []string{`<https://api.com/search?page=2> ; rel="      next "`},
			want:  "https://api.com/search?page=2",
		},
		{
			name:  `rel 1`,
			req:   `https://www.w3.org/somewhere_else`,
			links: []string{`<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`},
			want:  `https://www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152`,
		},
		{
			name:  `rel 2`,
			req:   `https://zombo.com/path/which/will/be/ignored`,
			links: []string{`</abc>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`},
			want:  `https://zombo.com/abc`,
		},
		{
			name:  `rel 3`,
			req:   `https://zombo.com/retained/path/except_last_element`,
			links: []string{`<abc/def?page=2>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`},
			want:  `https://zombo.com/retained/path/abc/def?page=2`,
		},
		{
			name:  `rel 4`,
			req:   `https://zombo.com/retained/path/including_last_element/`,
			links: []string{`<abc/def?page=2>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`},
			want:  `https://zombo.com/retained/path/including_last_element/abc/def?page=2`,
		},
		{
			name:  `rel 5 with []`,
			req:   `https://zombo.com/with_brackets?filter[foo]=bar`,
			links: []string{`</with_brackets?filter%5Bfoo%5D=bar&page=2>; rel="next"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`},
			want:  `https://zombo.com/with_brackets?filter[foo]=bar&page=2`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				req, _ = http.NewRequest("GET", testcase.req, nil)
				header = http.Header{"Link": testcase.links}
				resp   = &http.Response{Request: req, Header: header}
				want   = testcase.want
			)

			var have string
			u, err := api.GetNextLink(resp)
			if err != nil {
				t.Log(err)
			}
			if u != nil {
				have = u.String()
			}
			if want != have {
				t.Fatalf("want %q, have %q", want, have)
			}
		})
	}
}
