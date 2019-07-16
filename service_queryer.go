package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// serviceQueryer asks the Fastly API to resolve a set of service IDs to their
// names. This is necessary because names can be changed dynamically, and the
// fastly-exporter should reflect the most recent name.
type serviceQueryer struct {
	token     string
	whitelist map[string]bool // service IDs to use (optional; if not specified, allow all)
	updater   serviceUpdater
	manager   idUpdater
}

// nameUpdater is a consumer contract for the write side of the name cache.
// Whenever the service queryer gets a new mapping of service IDs to names,
// it will call this method to save that latest mapping.
type serviceUpdater interface {
	update(names map[string]nameVersion)
}

// idUpdater is a consumer contract for the write side of the monitor manager.
// Whenever the service queryer gets a new set of service IDs that should be
// monitored, it will call this method to save those latest IDs.
//
// Note that while this method is called regardless, it only has a meaningful
// effect when the fastly-exporter and service queryer are configured without
// any explicit service IDs, and thus should monitor *all* service IDs available
// to a Fastly token.
type idUpdater interface {
	update(ids []string)
}

func newServiceQueryer(token string, ids []string, updater serviceUpdater, manager idUpdater) *serviceQueryer {
	whitelist := map[string]bool{}
	for _, id := range ids {
		whitelist[id] = true
	}

	return &serviceQueryer{
		token:     token,
		whitelist: whitelist,
		updater:   updater,
		manager:   manager,
	}
}

// refresh the service ID to name mapping, updating the name updater (the name
// cache) and the ID updater (the monitor manager).
func (q *serviceQueryer) refresh(client httpClient) error {
	req, err := http.NewRequest("GET", "https://api.fastly.com/service", nil)
	if err != nil {
		return errors.Wrap(err, "error constructing API services request")
	}

	req.Header.Set("Fastly-Key", q.token)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making API services request")
	}
	defer resp.Body.Close()

	var response serviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "error decoding API services response")
	}

	var (
		services = map[string]nameVersion{}
		ids      []string
	)
	for _, tuple := range response {
		var (
			allowAll  = len(q.whitelist) == 0
			allowThis = q.whitelist[tuple.ID]
		)
		if allowAll || allowThis {
			services[tuple.ID] = nameVersion{tuple.Name, strconv.Itoa(tuple.Version)}
			ids = append(ids, tuple.ID)
		}
	}

	q.updater.update(services)
	q.manager.update(ids)
	return nil
}

type serviceResponse []struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version int    `json:"version"`
}
