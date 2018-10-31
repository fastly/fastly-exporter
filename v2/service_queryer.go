package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type serviceQueryer struct {
	token    string
	ids      []string // optional
	resolver nameUpdater
	manager  idUpdater
}

type nameUpdater interface {
	update(names map[string]string)
}

type idUpdater interface {
	update(ids []string)
}

func newServiceQueryer(token string, ids []string, resolver nameUpdater, manager idUpdater) *serviceQueryer {
	return &serviceQueryer{
		token:    token,
		ids:      ids,
		resolver: resolver,
		manager:  manager,
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

	var response serviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "error decuding API services response")
	}

	filter := map[string]bool{}
	for _, id := range q.ids {
		filter[id] = true
	}

	var (
		names = map[string]string{}
		ids   []string
	)
	for _, pair := range response {
		var (
			allowAll  = len(filter) == 0
			allowThis = filter[pair.ID]
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
