package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type serviceQueryer struct {
	token     string
	whitelist map[string]bool // service IDs to use (optional, if not specified, allow all)
	resolver  nameUpdater
	manager   idUpdater
}

type nameUpdater interface {
	update(names map[string]string)
}

type idUpdater interface {
	update(ids []string)
}

func newServiceQueryer(token string, ids []string, resolver nameUpdater, manager idUpdater) *serviceQueryer {
	whitelist := map[string]bool{}
	for _, id := range ids {
		whitelist[id] = true
	}

	return &serviceQueryer{
		token:     token,
		whitelist: whitelist,
		resolver:  resolver,
		manager:   manager,
	}
}

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
		names = map[string]string{}
		ids   []string
	)
	for _, pair := range response {
		var (
			allowAll  = len(q.whitelist) == 0
			allowThis = q.whitelist[pair.ID]
		)
		if allowAll || allowThis {
			names[pair.ID] = pair.Name
			ids = append(ids, pair.ID)
		}
	}

	q.resolver.update(names)
	q.manager.update(ids)
	return nil
}

type serviceResponse []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
