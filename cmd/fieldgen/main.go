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
		validateHelp = fs.Bool("validate-help", true, "validate metric definition Help field")
	)
	fs.Parse(os.Args[1:])

	var metrics []exporterMetric
	if err := readJSON(*metricsFile, &metrics); err != nil {
		return err
	}

	if *validateHelp {
		var noHelp []string
		for _, m := range metrics {
			if m.Help == "" || strings.Contains(m.Help, "TODO") {
				noHelp = append(noHelp, m.FieldName)
			}
		}
		if len(noHelp) > 0 {
			return fmt.Errorf("metrics with undefined Help strings: %s", strings.Join(noHelp, ", "))
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

	{
		unmapped := map[string]bool{}
		for _, f := range fields {
			unmapped[f.FieldName] = true
		}
		for _, m := range mappings {
			delete(unmapped, m.APIField)
			for _, pair := range m.APIFieldLabelMapping {
				delete(unmapped, pair[0])
			}
			for _, f := range m.APIFieldSizes {
				delete(unmapped, f)
			}
		}
		if len(unmapped) > 0 {
			var names []string
			for f := range unmapped {
				names = append(names, f)
			}
			return fmt.Errorf("unmapped API fields: %s", strings.Join(names, ", "))
		}
	}

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
	for _, m := range metrics {
		fmt.Printf("\t%s *prometheus.%sVec\n", m.FieldName, m.Type)
	}
	fmt.Printf("}\n")
	fmt.Printf("\n")
	fmt.Printf("// NewMetrics returns a new set of metrics registered to the registerer.\n")
	fmt.Printf("// Only metrics whose names pass the name filter are registered.\n")
	fmt.Printf("func NewMetrics(namespace, subsystem string, nameFilter filter.Filter, r prometheus.Registerer) (*Metrics, error) {\n")
	fmt.Printf("\tm := &Metrics{\n")
	fmt.Printf("\t\t" + `RealtimeAPIRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Subsystem: subsystem, Name: "realtime_api_requests_total", Help: "Total requests made to the real-time stats API.", }, []string{"service_id", "service_name", "result"}),` + "\n")
	fmt.Printf("\t\t" + `ServiceInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Subsystem: subsystem, Name: "service_info", Help: "Static gauge with service ID, name, and version information.", }, []string{"service_id", "service_name", "service_version"}),` + "\n")
	for _, m := range metrics {
		fmt.Printf("\t\t%s: %s,\n", m.FieldName, m.create())
	}
	fmt.Printf("\t}\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", registerBlock)
	fmt.Printf("\treturn m, nil\n")
	fmt.Printf("}\n\n")
	fmt.Printf("%s\n", getNameBlock)
	fmt.Printf("\n")
	fmt.Printf("// Process updates the metrics with data from the API response.\n")
	fmt.Printf("func Process(resp *APIResponse, serviceID, serviceName, serviceVersion string, m *Metrics) {\n")
	fmt.Printf("\tfor _, d := range resp.Data {\n")
	fmt.Printf("\t\tfor datacenter, stats := range d.Datacenter {\n")
	fmt.Printf("\t\t\tm.ServiceInfo.WithLabelValues(serviceID, serviceName, serviceVersion).Set(1)\n")
	fmt.Printf("\t\t\tm.RequestsTotal.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.Requests))\n")
	for _, m := range mappings {
		switch m.Kind {
		case "Counter":
			fmt.Printf("\t\t\tm.%s.WithLabelValues(serviceID, serviceName, datacenter).Add(float64(stats.%s))\n", m.ExporterMetric, m.APIField)
		case "CounterLabels":
			for _, pair := range m.APIFieldLabelMapping {
				fmt.Printf("\t\t\tm.%s.WithLabelValues(serviceID, serviceName, datacenter, \"%s\").Add(float64(stats.%s))\n", m.ExporterMetric, pair[1], pair[0])
			}
		case "Histogram":
			fmt.Printf("\t\t\tprocessHistogram(stats.%s, m.%s.WithLabelValues(serviceID, serviceName, datacenter))\n", m.APIField, m.ExporterMetric)
		case "ObjectSize":
			fmt.Printf("\t\t\tprocessObjectSizes(stats.ObjectSize1k, stats.ObjectSize10k, stats.ObjectSize100k, stats.ObjectSize1m, stats.ObjectSize10m, stats.ObjectSize100m, stats.ObjectSize1g, m.%s.WithLabelValues(serviceID, serviceName, datacenter))\n", m.ExporterMetric) // hacky hack
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
	FieldName   string   `json:"field_name"`
	Type        string   `json:"type"`
	MetricName  string   `json:"metric_name"`
	ExtraLabels []string `json:"extra_labels"`
	Help        string   `json:"help"`
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
	var (
		constructor = fmt.Sprintf("prometheus.New%sVec", m.Type)
		options     = fmt.Sprintf("prometheus.%sOpts", m.Type)
		constLabels = quoteList(append(standardLabels, m.ExtraLabels...))
	)
	return fmt.Sprintf(
		`%s(%s{Namespace: namespace, Subsystem: subsystem, Name: "%s", Help: "%s"}, []string{%s})`,
		constructor, options, m.MetricName, m.Help, constLabels,
	)
}

type apiField struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
	Key       string `json:"key"`
}

type mapping struct {
	ExporterMetric       string      `json:"exporter_metric"`
	Kind                 string      `json:"kind"`
	APIField             string      `json:"api_field"`               // kind=Counter, kind=Histogram
	APIFieldLabelMapping [][2]string `json:"api_field_label_mapping"` // Kind=CounterLabels
	APIFieldSizes        []string    `json:"api_field_sizes"`         // kind=ObjectSize
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

	"github.com/peterbourgon/fastly-exporter/pkg/filter"
	"github.com/prometheus/client_golang/prometheus"
)
`

var apiResponseBlock = strings.Replace(`
// APIResponse models the response from rt.fastly.com. It can get quite large;
// when there are lots of services being monitored, unmarshaling to this type is
// the CPU bottleneck of the program.
type APIResponse struct {
	Timestamp uint64 ·json:"Timestamp"·
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
		panic(fmt.Sprintf("programmer error: field %d/%d in Metrics type isn't a prometheus.Collector", i+1, v.NumField()))
	}
	if name := getName(c); !nameFilter.Permit(name) {
		continue
	}
	if err := r.Register(c); err != nil {
		return nil, fmt.Errorf("error registering metric %d/%d: %w", i+1, v.NumField(), err)
	}
}
`
