package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
)

type fixedResponseClient struct {
	code     int
	response string
}

func (c fixedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c.code)
		fmt.Fprint(w, c.response)
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

//
//
//

type paginatedResponseClient struct {
	responses []string
}

func (c paginatedResponseClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageIndex := page - 1
		if pageIndex >= len(c.responses) {
			http.Error(w, "page too large", 400)
			return
		}

		if pageIndex+1 < len(c.responses) {
			values := r.URL.Query()
			values.Set("page", strconv.Itoa(page+1))
			r.URL.RawQuery = values.Encode()
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, r.URL.String()))
		}

		fmt.Fprint(w, c.responses[pageIndex])
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}

//
//
//

type sequentialResponseClient struct {
	responses []string
}

func newSequentialResponseClient(responses ...string) *sequentialResponseClient {
	return &sequentialResponseClient{
		responses: responses,
	}
}

func (c *sequentialResponseClient) Do(req *http.Request) (*http.Response, error) {
	var response string
	if len(c.responses) <= 1 {
		response = c.responses[0]
	} else {
		response, c.responses = c.responses[0], c.responses[1:]
	}

	rec := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	}).ServeHTTP(rec, req)
	return rec.Result(), nil
}
