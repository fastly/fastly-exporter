package main

import (
	"reflect"
	"testing"
)

func TestServiceQueryerFixture(t *testing.T) {
	var (
		token   = "irrelevant"
		id      = "AbcDef123ghiJKlmnOPsq"
		ids     = []string{id}
		cache   = newServiceCache()
		manager = &mockManager{}
		queryer = newServiceQueryer(token, ids, cache, manager)
		client  = fixedResponseClient{200, serviceResponseFixture}
	)

	if err := queryer.refresh(client); err != nil {
		t.Fatalf("queryer.refresh: %v", err)
	}

	name, version := cache.resolve(id)
	if want, have := "My first service", name; want != have {
		t.Fatalf("service cache: name: want %q, have %q", want, have)
	}
	if want, have := "5", version; want != have {
		t.Fatalf("service cache: version: want %q, have %q", want, have)
	}

	if want, have := []string{id}, manager.ids; !reflect.DeepEqual(want, have) {
		t.Fatalf("monitor manager: want %v, have %v", want, have)
	}
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
