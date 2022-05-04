package api_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/fastly/fastly-exporter/pkg/filter"
)

func TestServiceCache(t *testing.T) {
	t.Parallel()

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
				ctx    = context.Background()
				client = fixedResponseClient{code: 200, response: serviceResponseLarge}
				cache  = api.NewServiceCache(client, "irrelevant_token", testcase.options...)
			)
			if err := cache.Refresh(ctx); err != nil {
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

func TestServiceCachePagination(t *testing.T) {
	t.Parallel()

	responses := []string{
		`[
			{ "version": 6, "name": "Service 1/1", "id": "c9407d61ae888d" },
			{ "version": 1, "name": "Service 1/2", "id": "cb32a38adf2e00" },
			{ "version": 6, "name": "Service 1/3", "id": "82de5396a46629" },
			{ "version": 2, "name": "Service 1/4", "id": "4200f01763cff9" }
		]`,
		`[
			{ "version": 7, "name": "Service 2/1", "id": "ce2976ac5a3e24" },
			{ "version": 3, "name": "Service 2/2", "id": "e1c2f1aa5fc341" }
		]`,
		`[
			{ "version": 7, "name": "Service 3/1", "id": "65544b504189bf" },
			{ "version": 5, "name": "Service 3/2", "id": "686ec4e72a836a" }
		]`,
	}

	var (
		ctx    = context.Background()
		client = paginatedResponseClient{responses}
		cache  = api.NewServiceCache(client, "irrelevant_token")
	)

	if err := cache.Refresh(ctx); err != nil {
		t.Fatal(err)
	}

	if want, have := []string{
		"4200f01763cff9", "65544b504189bf", "686ec4e72a836a", "82de5396a46629",
		"c9407d61ae888d", "cb32a38adf2e00", "ce2976ac5a3e24", "e1c2f1aa5fc341",
	}, cache.ServiceIDs(); !cmp.Equal(want, have) {
		t.Fatal(cmp.Diff(want, have))
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

const serviceResponseLarge = `[
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
