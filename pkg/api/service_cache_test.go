package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
)

func TestCache(t *testing.T) {
	var (
		s1 = api.Service{ID: "AbcDef123ghiJKlmnOPsq", Name: "My first service", Version: 5}
		s2 = api.Service{ID: "XXXXXXXXXXXXXXXXXXXXXX", Name: "Dummy service", Version: 1}
	)
	for _, testcase := range []struct {
		name    string
		options []api.ServiceCacheOption
		want    []api.Service
	}{
		{
			name:    "no options",
			options: nil,
			want:    []api.Service{s1, s2},
		},
		{
			name:    "allowlist both",
			options: []api.ServiceCacheOption{api.WithExplicitServiceIDs(s1.ID, s2.ID, "additional service ID")},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "allowlist one",
			options: []api.ServiceCacheOption{api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{s1},
		},
		{
			name:    "allowlist none",
			options: []api.ServiceCacheOption{api.WithExplicitServiceIDs("nonexistant service ID")},
			want:    []api.Service{},
		},
		{
			name:    "exact name include match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterAllowlist(`^` + s1.Name + `$`))},
			want:    []api.Service{s1},
		},
		{
			name:    "partial name include match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterAllowlist(`mmy`))},
			want:    []api.Service{s2},
		},
		{
			name:    "generous name include match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterAllowlist(`.*e.*`))},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "no name include match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterAllowlist(`not found`))},
			want:    []api.Service{},
		},
		{
			name:    "exact name exclude match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterBlocklist(`^` + s1.Name + `$`))},
			want:    []api.Service{s2},
		},
		{
			name:    "partial name exclude match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterBlocklist(`mmy`))},
			want:    []api.Service{s1},
		},
		{
			name:    "generous name exclude match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterBlocklist(`.*e.*`))},
			want:    []api.Service{},
		},
		{
			name:    "no name exclude match",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterBlocklist(`not found`))},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "name exclude and include",
			options: []api.ServiceCacheOption{api.WithNameFilter(filterAllowlistBlocklist(`.*e.*`, `mmy`))},
			want:    []api.Service{s1},
		},
		{
			name:    "single shard",
			options: []api.ServiceCacheOption{api.WithShard(1, 1)},
			want:    []api.Service{s1, s2},
		},
		{
			name:    "shard n0 m3",
			options: []api.ServiceCacheOption{api.WithShard(1, 3)},
			want:    []api.Service{s1}, // verified experimentally
		},
		{
			name:    "shard n1 m3",
			options: []api.ServiceCacheOption{api.WithShard(2, 3)},
			want:    []api.Service{s2}, // verified experimentally
		},
		{
			name:    "shard n2 m3",
			options: []api.ServiceCacheOption{api.WithShard(3, 3)},
			want:    []api.Service{}, // verified experimentally
		},
		{
			name:    "shard and service ID passing",
			options: []api.ServiceCacheOption{api.WithShard(1, 3), api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{s1},
		},
		{
			name:    "shard and service ID failing",
			options: []api.ServiceCacheOption{api.WithShard(2, 3), api.WithExplicitServiceIDs(s1.ID)},
			want:    []api.Service{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				client = fixedResponseClient{code: 200, response: serviceResponseFixture}
				cache  = api.NewServiceCache(client, "irrelevant_token", testcase.options...)
			)
			if err := cache.Refresh(); err != nil {
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

func filterAllowlist(a string) (f filter.Filter) {
	f.Allow(a)
	return f
}

func filterBlocklist(b string) (f filter.Filter) {
	f.Block(b)
	return f
}

func filterAllowlistBlocklist(a, b string) (f filter.Filter) {
	f.Allow(a)
	f.Block(b)
	return f
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
