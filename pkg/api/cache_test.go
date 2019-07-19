package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
)

func TestCache(t *testing.T) {
	var (
		s1 = api.Service{ID: "AbcDef123ghiJKlmnOPsq", Name: "My first service", Version: 5}
		s2 = api.Service{ID: "XXXXXXXXXXXXXXXXXXXXXX", Name: "Dummy service", Version: 1}
	)
	for _, testcase := range []struct {
		name    string
		options []api.CacheOption
		want    []api.Service
	}{
		{
			name:    "no options",
			options: nil,
			want:    []api.Service{s1, s2},
		},
		{
			name:    "whitelist both",
			options: []api.CacheOption{api.WithExplicitServiceIDs(s1.ID, s2.ID, "additional service ID")},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "whitelist one",
			options: []api.CacheOption{api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{s1},
		},
		{
			name:    "whitelist none",
			options: []api.CacheOption{api.WithExplicitServiceIDs("nonexistant service ID")},
			want:    []api.Service{},
		},
		{
			name:    "exact name match",
			options: []api.CacheOption{api.WithNameMatching(regexp.MustCompile(`^` + s1.Name + `$`))},
			want:    []api.Service{s1},
		},
		{
			name:    "partial name match",
			options: []api.CacheOption{api.WithNameMatching(regexp.MustCompile(`mmy`))},
			want:    []api.Service{s2},
		},
		{
			name:    "generous name match",
			options: []api.CacheOption{api.WithNameMatching(regexp.MustCompile(`.*e.*`))},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "no name match",
			options: []api.CacheOption{api.WithNameMatching(regexp.MustCompile(`not found`))},
			want:    []api.Service{},
		},
		{
			name:    "single shard",
			options: []api.CacheOption{api.WithShard(1, 1)},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "shard n0 m3",
			options: []api.CacheOption{api.WithShard(1, 3)},
			want:    []api.Service{s1}, // verified experimentally
		},
		{
			name:    "shard n1 m3",
			options: []api.CacheOption{api.WithShard(2, 3)},
			want:    []api.Service{s2}, // verified experimentally
		},
		{
			name:    "shard n2 m3",
			options: []api.CacheOption{api.WithShard(3, 3)},
			want:    []api.Service{}, // verified experimentally
		},
		{
			name:    "shard and service ID passing",
			options: []api.CacheOption{api.WithShard(1, 3), api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{s1},
		},
		{
			name:    "shard and service ID failing",
			options: []api.CacheOption{api.WithShard(2, 3), api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				cache  = api.NewCache("irrelevant_token", testcase.options...)
				client = fixedResponseClient{code: 200, response: serviceResponseFixture}
			)
			if err := cache.Refresh(client); err != nil {
				t.Fatal(err)
			}
			var (
				serviceIDs = cache.ServiceIDs()
				services   = make([]api.Service, len(serviceIDs))
			)
			for i, id := range serviceIDs {
				name, version, _ := cache.Metadata(id)
				services[i] = api.Service{ID: id, Name: name, Version: version}
			}
			if want, have := testcase.want, services; !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}
		})
	}
}

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

const serviceResponseFixture = `[
	{
		"version": 5,
		"name": "My first service",
		"created_at": "2018-07-26T06:13:51Z",
		"versions": [
			{
				"testing": false,
				"locked": true,
				"number": 1,
				"active": false,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-07-26T06:13:51Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T06:17:27Z",
				"deployed": false
			},
			{
				"testing": false,
				"locked": true,
				"number": 2,
				"active": false,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-07-26T06:15:47Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T20:30:44Z",
				"deployed": false
			},
			{
				"testing": false,
				"locked": true,
				"number": 3,
				"active": false,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-07-26T20:28:26Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T20:48:40Z",
				"deployed": false
			},
			{
				"testing": false,
				"locked": true,
				"number": 4,
				"active": false,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-07-26T20:47:58Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T21:35:32Z",
				"deployed": false
			},
			{
				"testing": false,
				"locked": true,
				"number": 5,
				"active": true,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-07-26T21:35:23Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T21:35:33Z",
				"deployed": false
			},
			{
				"testing": false,
				"locked": false,
				"number": 6,
				"active": false,
				"service_id": "AbcDef123ghiJKlmnOPsq",
				"staging": false,
				"created_at": "2018-09-28T04:02:22Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-09-28T04:05:33Z",
				"deployed": false
			}
		],
		"comment": "",
		"customer_id": "1a2a3a4azzzzzzzzzzzzzz",
		"updated_at": "2018-10-24T06:31:41Z",
		"id": "AbcDef123ghiJKlmnOPsq"
	},
	{
		"version": 1,
		"name": "Dummy service",
		"created_at": "2018-09-20T16:42:20Z",
		"versions": [
			{
				"testing": false,
				"locked": true,
				"number": 1,
				"active": true,
				"service_id": "XXXXXXXXXXXXXXXXXXXXXX",
				"staging": false,
				"created_at": "2018-09-20T16:42:20Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-09-20T16:42:21Z",
				"deployed": false
			}
		],
		"comment": "",
		"customer_id": "1a2a3a4azzzzzzzzzzzzzz",
		"updated_at": "2018-09-20T16:42:20Z",
		"id": "XXXXXXXXXXXXXXXXXXXXXX"
	}
]`
