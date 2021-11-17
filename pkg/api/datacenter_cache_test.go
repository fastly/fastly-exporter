package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbourgon/fastly-exporter/pkg/api"
)

func TestDatacenterCache(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name    string
		client  api.HTTPClient
		wantDCs []api.Datacenter
		wantErr error
	}{
		{
			name:    "small",
			client:  fixedResponseClient{code: http.StatusOK, response: datacentersResponseSmall},
			wantErr: nil,
			wantDCs: []api.Datacenter{
				{Code: "AMS", Name: "Amsterdam", Group: "Europe"},
				{Code: "WLG", Name: "Wellington", Group: "Asia/Pacific"},
			},
		},
		{
			name:    "error",
			client:  fixedResponseClient{code: http.StatusUnauthorized},
			wantErr: &api.Error{Code: http.StatusUnauthorized},
			wantDCs: []api.Datacenter{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				client = testcase.client
				cache  = api.NewDatacenterCache(client, "irrelevant token")
			)

			if want, have := testcase.wantErr, cache.Refresh(ctx); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}

			if want, have := testcase.wantDCs, cache.Datacenters(); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}
		})
	}
}

const datacentersResponseSmall = `
[
  {
    "code": "AMS",
    "name": "Amsterdam",
    "group": "Europe",
    "coordinates": {
      "x": 0,
      "y": 0,
      "latitude": 52.308613,
      "longitude": 4.763889
    },
    "shield": "amsterdam-nl"
  },
  {
    "code": "WLG",
    "name": "Wellington",
    "group": "Asia/Pacific",
    "coordinates": {
      "x": 0,
      "y": 0,
      "latitude": -41.327221,
      "longitude": 174.805278
    }
  }
]
`
