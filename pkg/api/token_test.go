package api_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/fastly/fastly-exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestTokenMetric(t *testing.T) {
	var (
		namespace = "fastly"
		subsystem = "rt"
	)
	client := api.NewTokenRecorder(fixedResponseClient{code: http.StatusOK, response: tokenReponseExpiresAt}, "")
	gatherer, _ := client.Gatherer(namespace, subsystem)
	client.Set(context.Background())

	expected := `
# HELP fastly_rt_token_expiration Unix timestamp of the expiration time of the Fastly API Token
# TYPE fastly_rt_token_expiration gauge
fastly_rt_token_expiration{token_id="id1234",user_id="user456"} 1.7764704e+09
`
	err := testutil.GatherAndCompare(gatherer, strings.NewReader(expected), "fastly_rt_token_expiration")
	if err != nil {
		t.Error(err)
	}
}

func TestTokenMetricWithoutExpiration(t *testing.T) {
	var (
		namespace = "fastly"
		subsystem = "rt"
	)
	client := api.NewTokenRecorder(fixedResponseClient{code: http.StatusOK, response: tokenReponseNoExpiry}, "")
	gatherer, _ := client.Gatherer(namespace, subsystem)
	client.Set(context.Background())

	expected := `
# HELP fastly_rt_token_expiration Unix timestamp of the expiration time of the Fastly API Token
# TYPE fastly_rt_token_expiration gauge
`
	err := testutil.GatherAndCompare(gatherer, strings.NewReader(expected), "fastly_rt_token_expiration")
	if err != nil {
		t.Error(err)
	}
}

const tokenReponseExpiresAt = `
{
  "id": "id1234",
  "user_id": "user456",
  "customer_id": "customer987",
  "name": "Fastly API Token",
  "last_used_at": "2024-04-18T13:37:06Z",
  "created_at": "2016-10-11T18:36:35Z",
  "expires_at": "2026-04-18T00:00:00Z"
}`

const tokenReponseNoExpiry = `
{
  "id": "id1234",
  "user_id": "user456",
  "customer_id": "customer987",
  "name": "Fastly API Token",
  "last_used_at": "2024-04-18T13:37:06Z",
  "created_at": "2016-10-11T18:36:35Z",
  "expires_at": null
}`
