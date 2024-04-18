package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TokenRecorder requests api.fastly.com/tokens/self once and sets a gauge metric
type TokenRecorder struct {
	client HTTPClient
	token  string
	metric *prometheus.GaugeVec
}

// NewTokenRecorder returns an empty token recorder. Use the
// Set method to get token data and set the gauge metric.
func NewTokenRecorder(client HTTPClient, token string) *TokenRecorder {
	return &TokenRecorder{
		client: client,
		token:  token,
	}
}

// Gatherer returns a Prometheus gatherer which will yield current
// token expiration as a gauge metric.
func (t *TokenRecorder) Gatherer(namespace, subsystem string) (prometheus.Gatherer, error) {
	registry := prometheus.NewRegistry()
	tokenExpiration := prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "token_expiration", Help: "Unix timestamp of the expiration time of the Fastly API Token"}, []string{"token_id", "user_id"})
	err := registry.Register(tokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("registering token collector: %w", err)
	}
	t.metric = tokenExpiration
	return registry, nil
}

// Set retreives token metadata from the Fastly API and sets the gauge metric
func (t *TokenRecorder) Set(ctx context.Context) error {
	token, err := t.getToken(ctx)
	if err != nil {
		return err
	}

	if !token.Expiration.IsZero() {
		t.metric.WithLabelValues(token.ID, token.UserID).Set(float64(token.Expiration.Unix()))
	}
	return nil
}

type token struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Expiration time.Time `json:"expires_at,omitempty"`
}

func (t *TokenRecorder) getToken(ctx context.Context) (*token, error) {
	uri := "https://api.fastly.com/tokens/self"

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("error constructing API tokens request: %w", err)
	}

	req.Header.Set("Fastly-Key", t.token)
	req.Header.Set("Accept", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing API tokens request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp)
	}

	var response token

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding API Tokens response: %w", err)
	}

	return &response, nil
}
