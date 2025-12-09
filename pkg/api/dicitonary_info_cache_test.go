package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
)

func TestDictionaryInfoCache(t *testing.T) {
	t.Parallel()

	svcClient := fixedResponseClient{code: 200, response: serviceResponseForDictionaryInfo}
	serviceCache := api.NewServiceCache(svcClient, "irrelevant_token")
	serviceCache.Refresh(context.Background())

	for _, testcase := range []struct {
		name        string
		client      api.HTTPClient
		wantDicts   []api.Dictionary
		wantErr     error
		wantEnabled bool
	}{
		{
			name:    "success",
			client:  newSequentialResponseClient(dictionaryListResponse, dictionaryInfoResponseOne, dictionaryInfoResponseTwo),
			wantErr: nil,
			wantDicts: []api.Dictionary{
				{
					ServiceID:      "qwertMz1ncwA0KC3TBloku",
					ServiceName:    "test",
					Version:        575,
					DictionaryID:   "AsDfAOwejFaNQfAaFRxv2J",
					DictionaryName: "dict1",
					Digest:         "a889640be09d865b91194a896de19deea22823d707977807f42da717d85d372e",
					ItemCount:      9,
					LastUpdatedTS:  1724975476,
				},
				{
					ServiceID:      "qwertMz1ncwA0KC3TBloku",
					ServiceName:    "test",
					Version:        575,
					DictionaryID:   "foobGF34VOKSM0yVzBaQh7",
					DictionaryName: "dict2",
					Digest:         "785b9d481ca4455b734a9781ee73c07a44fe81f67f32ba286babd50afaa841d5",
					ItemCount:      11,
					LastUpdatedTS:  1765231819,
				},
			},
			wantEnabled: true,
		},
		{
			name:   "success_and_empty",
			client: fixedResponseClient{code: http.StatusOK, response: dictionaryListResponseEmpty},

			wantErr:     nil,
			wantDicts:   []api.Dictionary{},
			wantEnabled: true,
		},
		{
			name:        "forbidden",
			client:      fixedResponseClient{code: http.StatusForbidden},
			wantErr:     &api.Error{Code: http.StatusForbidden},
			wantDicts:   nil,
			wantEnabled: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				client = testcase.client
				cache  = api.NewDictionaryInfoCache(client, "irrelevant token", log.NewNopLogger(), serviceCache, true)
			)

			if want, have := testcase.wantErr, cache.Refresh(ctx); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}

			if want, have := testcase.wantDicts, cache.Dictionaries(); !cmp.Equal(want, have) {
				for _, x := range want {
					fmt.Println(x)
				}
				t.Fatal(cmp.Diff(want, have))
			}

			if want, have := testcase.wantEnabled, cache.Enabled(); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}
		})
	}
}

const dictionaryListResponseEmpty = `[]`

const dictionaryListResponse = `
[
  {
    "name": "dict1",
    "id": "AsDfAOwejFaNQfAaFRxv2J",
    "updated_at": "2025-10-31T13:27:40Z",
    "created_at": "2021-11-05T16:10:59Z",
    "service_id": "qwertMz1ncwA0KC3TBloku",
    "version": 575,
    "write_only": false,
    "deleted_at": null
  },
  {
    "name": "dict2",
    "id": "foobGF34VOKSM0yVzBaQh7",
    "updated_at": "2025-10-31T13:27:40Z",
    "deleted_at": null,
    "service_id": "qwertMz1ncwA0KC3TBloku",
    "created_at": "2024-01-29T17:51:55Z",
    "write_only": false,
    "version": 575
  }
]
`

const dictionaryInfoResponseOne = `
{
  "digest": "a889640be09d865b91194a896de19deea22823d707977807f42da717d85d372e",
  "last_updated": "2024-08-29 23:51:16",
  "item_count": 9
}
`

const dictionaryInfoResponseTwo = `
{
  "item_count": 11,
  "digest": "785b9d481ca4455b734a9781ee73c07a44fe81f67f32ba286babd50afaa841d5",
  "last_updated": "2025-12-08 22:10:19"
}
`

const serviceResponseForDictionaryInfo = `[
	{
		"version": 575,
		"name": "test",
		"created_at": "2018-07-26T06:13:51Z",
		"versions": [
			{
				"testing": false,
				"locked": true,
				"number": 1,
				"active": false,
				"service_id": "qwertMz1ncwA0KC3TBloku",
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
				"number": 575,
				"active": true,
				"service_id": "qwertMz1ncwA0KC3TBloku",
				"staging": false,
				"created_at": "2018-07-26T21:35:23Z",
				"deleted_at": null,
				"comment": "",
				"updated_at": "2018-07-26T21:35:33Z",
				"deployed": false
			}
		],
		"comment": "",
		"customer_id": "1a2a3a4azzzzzzzzzzzzzz",
		"updated_at": "2018-10-24T06:31:41Z",
		"id": "qwertMz1ncwA0KC3TBloku"
	}
]`
