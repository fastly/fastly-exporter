package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
)

func TestCertificateCache(t *testing.T) {
	t.Parallel()

	for _, testcase := range []struct {
		name        string
		client      api.HTTPClient
		wantCerts   []api.Certificate
		wantErr     error
		wantEnabled bool
	}{
		{
			name:    "success",
			client:  fixedResponseClient{code: http.StatusOK, response: certificatesResponseLarge},
			wantErr: nil,
			wantCerts: []api.Certificate{
				{
					ID: "ZfkhTtm4LdaOprVcdsffx4",
					Attributes: api.Attributes{
						CN:       "first.example1.com",
						Name:     "first.example1.com",
						Issuer:   "First CA",
						NotAfter: "2023-06-25T01:09:23.000Z",
						SN:       "52135557897532112355784498781325912334",
					},
				},
				{
					ID: "YkUe3r6S3zN4m6lVCd3sGc",
					Attributes: api.Attributes{
						CN:       "second.example2.com",
						Name:     "My Testing Cert",
						Issuer:   "Second CA",
						NotAfter: "2024-08-29T11:07:33.000Z",
						SN:       "11106091125671337225612345678987654321",
					},
				},
			},
			wantEnabled: true,
		},
		{
			name:        "success_and_empty",
			client:      fixedResponseClient{code: http.StatusOK, response: certificatesResponseEmpty},
			wantErr:     nil,
			wantCerts:   []api.Certificate{},
			wantEnabled: true,
		},
		{
			name:        "forbidden",
			client:      fixedResponseClient{code: http.StatusForbidden},
			wantErr:     &api.Error{Code: http.StatusForbidden},
			wantCerts:   nil,
			wantEnabled: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				client = testcase.client
				cache  = api.NewCertificateCache(client, "irrelevant token", true, log.NewNopLogger())
			)

			if want, have := testcase.wantErr, cache.Refresh(ctx); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}

			if want, have := testcase.wantCerts, cache.Certificates(); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}

			if want, have := testcase.wantEnabled, cache.Enabled(); !cmp.Equal(want, have) {
				t.Fatal(cmp.Diff(want, have))
			}
		})
	}
}

const certificatesResponseEmpty = `
{
  "data": [],
  "links": {
    "self": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=20&sort=created_at",
    "first": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=20&sort=created_at",
    "prev": null,
    "next": null,
    "last": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=27&sort=created_at"
  },
  "meta": {
    "per_page": 20,
    "current_page": 1,
    "record_count": 0,
    "total_pages": 1
  }
}
`

const certificatesResponseLarge = `
{
  "meta": {
    "total_pages": 1,
    "record_count": 3,
    "current_page": 1,
    "per_page": 20
  },
  "links": {
    "last": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=20&sort=created_at",
    "next": null,
    "prev": null,
    "first": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=20&sort=created_at",
    "self": "https://api.fastly.com/tls/certificates?page%5Bnumber%5D=1&page%5Bsize%5D=20&sort=created_at"
  },
  "data": [
    { 
      "relationships": {
        "tls_domains": {
          "data": [
            { 
              "type": "tls_domain",
              "id": "abcd.first.example1.com"
            },
            { 
              "type": "tls_domain",
              "id": "1234.first.example1.com"
            },
            { 
              "type": "tls_domain",
              "id": "sub-domain.first.example1.com"
            }
          ]
        }
      },
      "attributes": {
        "updated_at": "2022-07-05T04:48:21.000Z",
        "signature_algorithm": "SHA256-RSA",
        "created_at": "2022-07-05T04:48:21.000Z",
        "issued_to": "first.example1.com",
        "issuer": "First CA",
        "name": "first.example1.com",
        "not_after": "2023-06-25T01:09:23.000Z",
        "not_before": "2022-06-25T01:09:24.000Z",
        "replace": false,
        "serial_number": "52135557897532112355784498781325912334"
      },
      "type": "tls_certificate",
      "id": "ZfkhTtm4LdaOprVcdsffx4"
    },
    { 
      "relationships": {
        "tls_domains": {
          "data": [
            { 
              "type": "tls_domain",
              "id": "abcd1234.second.example2.com"
            }
          ]
        }
      },
      "attributes": {
        "updated_at": "2023-09-09T14:46:31.000Z",
        "signature_algorithm": "SHA256-RSA",
        "created_at": "2023-09-09T14:46:31.000Z",
        "issued_to": "second.example2.com",
        "issuer": "Second CA",
        "name": "My Testing Cert",
        "not_after": "2024-08-29T11:07:33.000Z",
        "not_before": "2023-08-29T11:07:34.000Z",
        "replace": false,
        "serial_number": "11106091125671337225612345678987654321"
      },
      "type": "tls_certificate",
      "id": "YkUe3r6S3zN4m6lVCd3sGc"
    }
  ]
}
`
