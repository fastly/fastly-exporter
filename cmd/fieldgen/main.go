package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := exec(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func exec() error {
	fs := flag.NewFlagSet("fieldgen", flag.ExitOnError)
	var (
		metricsFile  = fs.String("metrics", "exporter_metrics.json", "JSON file containing metric definitions")
		fieldsFile   = fs.String("fields", "api_fields.json", "JSON file containing API field definitions")
		mappingsFile = fs.String("mappings", "mappings.json", "JSON file containing metric-to-field mappings")
		invalid      = fs.Bool("invalid", false, "skip validations")
	)
	fs.Parse(os.Args[1:])

	var metrics []exporterMetric
	if err := readJSON(*metricsFile, &metrics); err != nil {
		return err
	}

	if !*invalid {
		var (
			noHelp   []string
			todoHelp []string
		)
		for _, m := range metrics {
			if m.Help == "" {
				noHelp = append(noHelp, m.FieldName)
			}
			if strings.Contains(m.Help, "TODO") {
				todoHelp = append(todoHelp, m.FieldName)
			}
		}
		if len(noHelp) > 0 {
			return fmt.Errorf("metrics with undefined Help strings: %s", strings.Join(noHelp, ", "))
		}
		if len(todoHelp) > 0 {
			fmt.Fprintf(os.Stderr, "warning: metrics with TODO Help strings: %s\n", strings.Join(todoHelp, ", "))
		}
	}

	var fields []apiField
	if err := readJSON(*fieldsFile, &fields); err != nil {
		return err
	}

	var mappings []mapping
	if err := readJSON(*mappingsFile, &mappings); err != nil {
		return err
	}

	if !*invalid {
		unmappedAPIFields := map[string]bool{}
		for _, f := range fields {
			unmappedAPIFields[f.FieldName] = true
		}
		for _, m := range mappings {
			delete(unmappedAPIFields, m.APIField)
			for _, pair := range m.APIFieldLabels {
				delete(unmappedAPIFields, pair[0])
			}
			for _, f := range m.APIFieldSizes {
				delete(unmappedAPIFields, f)
			}
		}
		if len(unmappedAPIFields) > 0 {
			var names []string
			for f := range unmappedAPIFields {
				names = append(names, f)
			}
			return fmt.Errorf("unmapped API fields: %s", strings.Join(names, ", "))
		}
	}

	fmt.Printf("// Code generated by fieldgen; DO NOT EDIT.\n")
	fmt.Printf("\n")
	fmt.Printf("package gen\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", importBlock)
	fmt.Printf("\n")
	fmt.Printf("%s\n", apiResponseBlock)
	fmt.Printf("\n")
	fmt.Printf("// Datacenter models the per-datacenter portion of the rt.fastly.com response.\n")
	fmt.Printf("type Datacenter struct {\n")
	for _, f := range fields {
		fmt.Printf("\t%s %s `json:\"%s\"`\n", f.FieldName, f.Type, f.Key)
	}
	fmt.Printf("}\n")
	fmt.Printf("\n")
	fmt.Printf("// Metrics collects all of the Prometheus metrics exported by this service.\n")
	fmt.Printf("type Metrics struct {\n")
	fmt.Printf("\tRealtimeAPIRequestsTotal *prometheus.CounterVec\n")
	fmt.Printf("\tServiceInfo *prometheus.GaugeVec\n")
	fmt.Printf("\tLastSuccessfulResponse *prometheus.GaugeVec\n")
	for _, m := range metrics {
		fmt.Printf("\t%s *prometheus.%sVec\n", m.FieldName, m.Type)
	}
	fmt.Printf("}\n")
	fmt.Printf("\n")
	fmt.Printf("// NewMetrics returns a new set of metrics registered to the registerer.\n")
	fmt.Printf("// Only metrics whose names pass the name filter are registered.\n")
	fmt.Printf("func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) *Metrics {\n")
	fmt.Printf("\tm := Metrics{\n")
	fmt.Printf("\t\t" + `RealtimeAPIRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "realtime_api_requests_total", Help: "Total requests made to the real-time stats API.", }, []string{"service_id", "service_name", "result"}),` + "\n")
	fmt.Printf("\t\t" + `ServiceInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "service_info", Help: "Static gauge with service ID, name, and version information.", }, []string{"service_id", "service_name", "service_version"}),` + "\n")
	fmt.Printf("\t\t" + `LastSuccessfulResponse: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "last_successful_response", Help: "Unix timestamp of the last successful response received from the real-time stats API.", }, []string{"service_id", "service_name"}),` + "\n")
	for _, m := range metrics {
		fmt.Printf("\t\t%s: %s,\n", m.FieldName, m.create())
	}
	fmt.Printf("\t}\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", registerBlock)
	fmt.Printf("\treturn &m\n")
	fmt.Printf("}\n\n")
	fmt.Printf("%s\n", getNameBlock)
	fmt.Printf("\n")
	fmt.Printf("// Process updates the metrics with data from the API response.\n")
	fmt.Printf("func Process(response *APIResponse, serviceID, serviceName, serviceVersion string, m *Metrics) {\n")
	fmt.Printf("\tfor _, d := range response.Data {\n")
	fmt.Printf("\t\tfor datacenter, stats := range d.Datacenter {\n")
	for _, m := range mappings {
		switch m.Kind {
		case "Counter":
			fmt.Printf("\t\t\tm.%s.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.%s))\n", m.ExporterMetric, m.APIField)
		case "Counter1000":
			fmt.Printf("\t\t\tm.%s.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.%s) / 10000.0)\n", m.ExporterMetric, m.APIField)
		case "CounterLabels":
			for _, pair := range m.APIFieldLabels {
				fmt.Printf("\t\t\tm.%s.WithLabelValues(serviceID, serviceName, datacenter, \"%s\").Add(float64(stats.%s))\n", m.ExporterMetric, pair[1], pair[0])
			}
		case "Histogram":
			fmt.Printf("\t\t\tprocessHistogram(stats.%s, m.%s.WithLabelValues(serviceID, serviceName, datacenter))\n", m.APIField, m.ExporterMetric)
		case "ObjectSize":
			fmt.Printf("\t\t\tprocessObjectSizes(stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k, stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m, stats.ObjectSize1g, m.%s.WithLabelValues(serviceID, serviceName, datacenter))\n", m.ExporterMetric) // hacky hack
		case "Ignored":
			//
		default:
			fmt.Printf("\t\t\t// %s: unknown mapping kind %q\n", m.ExporterMetric, m.Kind)
		}
	}
	fmt.Printf("\t\t}\n")
	fmt.Printf("\t}\n")
	fmt.Printf("}\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", processBlock)
	return nil
}

func readJSON(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(data)
}

type exporterMetric struct {
	FieldName   string    `json:"field_name"`
	Type        string    `json:"type"`
	MetricName  string    `json:"metric_name"`
	ExtraLabels []string  `json:"extra_labels"`
	Help        string    `json:"help"`
	Buckets     []float64 `json:"buckets"`
}

var standardLabels = []string{"service_id", "service_name", "datacenter"}

func quoteList(a []string) string {
	b := make([]string, len(a))
	for i, s := range a {
		b[i] = `"` + s + `"`
	}
	return strings.Join(b, ", ")
}

func (m exporterMetric) create() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, `prometheus.New%sVec(prometheus.%sOpts{`, m.Type, m.Type)
	fmt.Fprintf(&sb, `Namespace: namespace`)
	fmt.Fprintf(&sb, `, Subsystem: subsystem`)
	fmt.Fprintf(&sb, `, Name: "%s"`, m.MetricName)
	fmt.Fprintf(&sb, `, Help: "%s"`, m.Help)
	if len(m.Buckets) > 0 {
		fmt.Fprintf(&sb, `, Buckets: []float64{%s}`, renderFloats(m.Buckets))
	}
	fmt.Fprintf(&sb, `}, []string{%s}`, quoteList(append(standardLabels, m.ExtraLabels...)))
	fmt.Fprintf(&sb, `)`)

	return sb.String()
}

type apiField struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
	Key       string `json:"key"`
}

type mapping struct {
	ExporterMetric string      `json:"exporter_metric"`
	Kind           string      `json:"kind"`
	APIField       string      `json:"api_field"`        // kind=Counter, kind=Histogram
	APIFieldLabels [][2]string `json:"api_field_labels"` // Kind=CounterLabels
	APIFieldSizes  []string    `json:"api_field_sizes"`  // kind=ObjectSize
}

func renderFloats(a []float64) string {
	s := make([]string, len(a))
	for i, f := range a {
		s[i] = fmt.Sprint(f)
	}
	return strings.Join(s, ", ")
}

//
//
//

const importBlock = `
import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/fastly/fastly-exporter/pkg/filter"
	"github.com/prometheus/client_golang/prometheus"
)
`

var apiResponseBlock = strings.Replace(`
// APIResponse models the response from rt.fastly.com. It can get quite large;
// when there are lots of services being monitored, unmarshaling to this type is
// the CPU bottleneck of the program.
type APIResponse struct {
	Timestamp uint64 ·json:"Timestamp"·
	AggregateDelay int64 ·json:"AggregateDelay"·
	Data      []struct {
		Datacenter map[string]Datacenter ·json:"datacenter"·
		Aggregated Datacenter            ·json:"aggregated"·
		Recorded   uint64                ·json:"recorded"·
	} ·json:"Data"·
	Error string ·json:"error"·
}
`, "·", "`", -1)

const getNameBlock = `
var descNameRegex = regexp.MustCompile("fqName: \"([^\"]+)\"")

func getName(c prometheus.Collector) string {
	d := make(chan *prometheus.Desc, 1)
	c.Describe(d)
	desc := (<-d).String()
	matches := descNameRegex.FindAllStringSubmatch(desc, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		return matches[0][1]
	}
	return ""
}
`

const processBlock = `
func processHistogram(src map[string]uint64, obs prometheus.Observer) {
	for str, count := range src {
		ms, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		s := float64(ms) / 1e3
		for i := 0; i < int(count); i++ {
			obs.Observe(s)
		}
	}
}

func processObjectSizes(n1k, n10k, n100k, n1m, n10m, n100m, n1g uint64, obs prometheus.Observer) {
	for v, n := range map[uint64]uint64{
		1 * 1024:           n1k,
		10 * 1024:          n10k,
		100 * 1024:         n100k,
		1 * 1000 * 1024:    n1m,
		10 * 1000 * 1024:   n10m,
		100 * 1000 * 1024:  n100m,
		1000 * 1000 * 1024: n1g,
	} {
		for i := uint64(0); i < n; i++ {
			obs.Observe(float64(v))
		}
	}
}
`

const registerBlock = `
for i, v := 0, reflect.ValueOf(m); i < v.NumField(); i++ {
	c, ok := v.Field(i).Interface().(prometheus.Collector)
	if !ok {
		panic(fmt.Errorf("field %d/%d in Metrics type isn't a prometheus.Collector", i+1, v.NumField()))
	}
	if name := getName(c); !nameFilter.Permit(name) {
		continue
	}
	if err := r.Register(c); err != nil {
		panic(fmt.Errorf("error registering metric %d/%d: %w", i+1, v.NumField(), err))
	}
}
`
