package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
)

func TestProductCache(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name      string
		client    api.HTTPClient
		wantProds map[string]bool
		wantErr   error
	}{
		{
			name:    "success",
			client:  fixedResponseClient{response: productsResponse, code: http.StatusOK},
			wantErr: nil,
			wantProds: map[string]bool{
				"origin_inspector": true,
			},
		},
		{
			name:      "noCustomerSuccess",
			client:    fixedResponseClient{response: noCustomerResponse, code: http.StatusOK},
			wantErr:   nil,
			wantProds: map[string]bool{},
		},
		{
			name:      "error",
			client:    fixedResponseClient{code: http.StatusUnauthorized},
			wantErr:   &api.Error{Code: http.StatusUnauthorized},
			wantProds: map[string]bool{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				client = testcase.client
				cache  = api.NewProductCache(client, "irrelevant token", log.NewNopLogger())
			)

			// err
			if want, have := testcase.wantErr, cache.Refresh(ctx); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}

			for k, v := range testcase.wantProds {
				if v != cache.HasAccess(k) {
					t.Fatalf("expected %v, got %v for %v", v, cache.HasAccess(k), k)
				}
			}
		})
	}
}

// only products that the customer is entitled to are returned from the /entitlements endpoint
const productsResponse = `
{
  "customers": [
    {
      "contracts": [
        {
          "items": [
            {
              "product_id": "origin_inspector"
            },
            {
              "other_key": "other"
            }
          ]
        },
        {
          "other_key_2": [
          ]
        }
      ]
    },
    {
      "other key": {}
    }
  ]
}
`

const noCustomerResponse = `
{
  "metadata": {}
}
`
