package filter_test

import (
	"testing"

	"github.com/peterbourgon/fastly-exporter/pkg/filter"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name      string
		whitelist []string
		blacklist []string
		inputs    map[string]bool
	}{
		{
			name: "default allow",
			inputs: map[string]bool{
				"anything": true,
				"":         true,
			},
		},
		{
			name:      "single whitelist",
			whitelist: []string{"foo"},
			inputs: map[string]bool{
				"foo": true,
				"bar": false,
				"":    false,
			},
		},
		{
			name:      "multiple whitelist",
			whitelist: []string{"foo", "bar"},
			inputs: map[string]bool{
				"foo": true,
				"bar": true,
				"baz": false,
				"":    false,
			},
		},
		{
			name:      "single blacklist",
			blacklist: []string{"foo"},
			inputs: map[string]bool{
				"foo": false,
				"bar": true,
				"":    true,
			},
		},
		{
			name:      "multiple blacklist",
			blacklist: []string{"foo", "bar"},
			inputs: map[string]bool{
				"foo": false,
				"bar": false,
				"baz": true,
				"":    true,
			},
		},
		{
			name:      "whitelist and blacklist",
			whitelist: []string{"foo", "bar"},
			blacklist: []string{"baz", "qux"},
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
			whitelist: []string{"[123]xx"},
			blacklist: []string{"bad$"},
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
			t.Parallel()

			var f filter.Filter
			for _, s := range testcase.whitelist {
				if err := f.Whitelist(s); err != nil {
					t.Fatalf("Whitelist(%s): %v", s, err)
				}
			}
			for _, s := range testcase.blacklist {
				if err := f.Blacklist(s); err != nil {
					t.Fatalf("Blacklist(%s): %v", s, err)
				}
			}
			for input, want := range testcase.inputs {
				if have := f.Allow(input); want != have {
					t.Errorf("Allow(%q): want %v, have %v", input, want, have)
				}
			}
		})
	}
}
