package prom

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/peterbourgon/fastly-exporter/pkg/gen"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry collects Prometheus metrics on a per-service basis.
//
// Writers (i.e. rt.Subscribers) should call MetricsFor with their specific
// service ID, and update the returned set of Prometheus metrics. Readers (i.e.
// Prometheus) can scrape metrics for all services via the `/metrics` endpoint,
// or a single service via `/metrics?target=<service ID>`.
//
// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config
type Registry struct {
	mtx              sync.Mutex
	version          string
	namespace        string
	subsystem        string
	metricNameFilter filter.Filter
	byServiceID      map[string]*metricsRegistry
	defaultGatherers []prometheus.Gatherer

	http.Handler
}

// NewRegistry returns a new and empty registry for Prometheus metrics.
func NewRegistry(version, namespace, subsystem string, metricNameFilter filter.Filter, defaultGatherers ...prometheus.Gatherer) *Registry {
	r := &Registry{
		version:          version,
		namespace:        namespace,
		subsystem:        subsystem,
		metricNameFilter: metricNameFilter,
		byServiceID:      map[string]*metricsRegistry{},
		defaultGatherers: defaultGatherers,
	}

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Methods("GET").Path("/").HandlerFunc(r.handleIndex)
	router.Methods("GET").Path("/sd").HandlerFunc(r.handleServiceDiscovery)
	router.Methods("GET").Path("/metrics").HandlerFunc(r.handleMetrics)
	r.Handler = router

	return r
}

// metricsRegistry combines a set of metrics for a single Fastly service with a
// Prometheus registry that yields those metrics. The registry can be combined
// with other registries and served as a single set of metrics via the
// prometheus.Gatherers helper type.
type metricsRegistry struct {
	metrics  *gen.Metrics
	registry *prometheus.Registry
}

// MetricsFor returns a set of Prometheus metrics for a specific service, with
// the expectation that callers will update those metrics with data retrieved
// from the Fastly real-time stats API.
func (r *Registry) MetricsFor(serviceID string) *gen.Metrics {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	mr, ok := r.byServiceID[serviceID]
	if !ok {
		registry := prometheus.NewRegistry()
		metrics := gen.NewMetrics(r.namespace, r.subsystem, r.metricNameFilter, registry)
		mr = &metricsRegistry{metrics, registry}
		r.byServiceID[serviceID] = mr // TODO(pb): at some point, expire and remove?
	}

	return mr.metrics
}

func (r *Registry) handleIndex(w http.ResponseWriter, req *http.Request) {
	type link struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	links := []link{
		{"/sd", "Service discovery"},
		{"/metrics", "Metrics for all services"},
	}

	for _, serviceID := range r.serviceIDs() {
		query := url.Values{"target": []string{serviceID}}.Encode()
		path := "/metrics?" + query
		name := "Metrics for service " + serviceID
		links = append(links, link{path, name})
	}

	accept := req.Header.Get("accept")
	switch {
	case strings.Contains(accept, "text/html"):
		w.Header().Set("content-type", "text/html; charset=utf-8")
		indexTemplate.Execute(w, struct {
			Version string
			Links   []link
		}{
			Version: r.version,
			Links:   links,
		})

	case strings.Contains(accept, "application/json"):
		w.Header().Set("content-type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(links)

	default:
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		for _, link := range links {
			fmt.Fprintf(w, "%s: %s\n", link.Path, link.Name)
		}
	}
}

func (r *Registry) handleServiceDiscovery(w http.ResponseWriter, req *http.Request) {
	targets := r.serviceIDs()

	response := []struct {
		Targets []string `json:"targets"`
	}{
		{Targets: targets},
	}

	buf, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write(buf)
}

func (r *Registry) handleMetrics(w http.ResponseWriter, req *http.Request) {
	var (
		target    = req.URL.Query().Get("target") // empty target string means all targets
		gatherers = prometheus.Gatherers(append(r.defaultGatherers, r.servicesGathererFor(target)))
		handler   = promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
	)
	handler.ServeHTTP(w, req)
}

func (r *Registry) serviceIDs() []string {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	serviceIDs := make([]string, 0, len(r.byServiceID))
	for serviceID := range r.byServiceID {
		serviceIDs = append(serviceIDs, serviceID)
	}

	sort.Strings(serviceIDs)

	return serviceIDs
}

func (r *Registry) servicesGathererFor(target string) prometheus.Gatherer {
	var allow func(candidate string) bool
	switch target {
	case "":
		allow = func(candidate string) bool { return true }
	default:
		allow = func(candidate string) bool { return candidate == target }
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var gatherers prometheus.Gatherers
	for serviceID, mr := range r.byServiceID {
		if allow(serviceID) {
			gatherers = append(gatherers, mr.registry)
		}
	}

	return gatherers
}

var indexTemplate = template.Must(template.New("").Parse(`
<html>
<head>
<title>fastly-exporter</title>
<style>body { margin: 1em 2em; font-size: 16pt; }</style>
</head>
<body>
<p>fastly-exporter{{ if .Version }} version <strong>{{ .Version }}</strong>{{ end }}</p>
<ul>
{{ range .Links -}}
<li><a href="{{ .Path }}">{{ .Name }}</a></li>
{{ end -}}
</ul>
</body>
</html>
`))
