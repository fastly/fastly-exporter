package filter_test

import (
	"testing"

	"github.com/fastly/fastly-exporter/pkg/filter"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name      string
		allowlist []string
		blocklist []string
		inputs    map[string]bool // to Permit
	}{
		{
			name: "default allow",
			inputs: map[string]bool{
				"anything": true,
				"":         true,
			},
		},
		{
			name:      "single allowlist",
			allowlist: []string{"foo"},
			inputs: map[string]bool{
				"foo": true,
				"bar": false,
				"":    false,
			},
		},
		{
			name:      "multiple allowlist",
			allowlist: []string{"foo", "bar"},
			inputs: map[string]bool{
				"foo": true,
				"bar": true,
				"baz": false,
				"":    false,
			},
		},
		{
			name:      "single blocklist",
			blocklist: []string{"foo"},
			inputs: map[string]bool{
				"foo": false,
				"bar": true,
				"":    true,
			},
		},
		{
			name:      "multiple blocklist",
			blocklist: []string{"foo", "bar"},
			inputs: map[string]bool{
				"foo": false,
				"bar": false,
				"baz": true,
				"":    true,
			},
		},
		{
			name:      "allowlist and blocklist",
			allowlist: []string{"foo", "bar"},
			blocklist: []string{"baz", "qux"},
			inputs: map[string]bool{
				"foo":           true,
				"foo bar":       true,
				"some bar blah": true,
				"foo bar baz":   false,
				"bar baz":       false,
				"baz":           false,
				"fo ba":         false,
				"":              false,
			},
		},
		{
			name:      "actual regex",
			allowlist: []string{"[123]xx"},
			blocklist: []string{"bad$"},
			inputs: map[string]bool{
				"1xx":     true,
				"2xx_ok":  true,
				"3xx_bad": false,
				"4xx":     false,
				"5xx":     false,
				"bad_2xx": true,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			testcase := testcase
			t.Parallel()

			var f filter.Filter
			for _, s := range testcase.allowlist {
				if err := f.Allow(s); err != nil {
					t.Fatalf("Allow(%s): %v", s, err)
				}
			}
			for _, s := range testcase.blocklist {
				if err := f.Block(s); err != nil {
					t.Fatalf("Block(%s): %v", s, err)
				}
			}
			for input, want := range testcase.inputs {
				if have := f.Permit(input); want != have {
					t.Errorf("Allow(%q): want %v, have %v", input, want, have)
				}
			}
		})
	}
}
