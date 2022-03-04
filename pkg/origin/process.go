package origin

// Process updates the metrics with data from the API response.
func Process(response *Response, serviceID, serviceName, serviceVersion string, m *Metrics) {
	for _, d := range response.Data {
		for datacenter, byOrigin := range d.Datacenter {
			for origin, stats := range byOrigin {
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin).Add(float64(stats.RespBodyBytes))
				m.RespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin).Add(float64(stats.RespHeaderBytes))
				m.ResponsesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin).Add(float64(stats.Responses))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "1xx").Add(float64(stats.Status1xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "200").Add(float64(stats.Status200))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "204").Add(float64(stats.Status204))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "2xx").Add(float64(stats.Status2xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "301").Add(float64(stats.Status301))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "302").Add(float64(stats.Status302))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "304").Add(float64(stats.Status304))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "3xx").Add(float64(stats.Status3xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "400").Add(float64(stats.Status400))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "401").Add(float64(stats.Status401))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "403").Add(float64(stats.Status403))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "404").Add(float64(stats.Status404))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "416").Add(float64(stats.Status416))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "4xx").Add(float64(stats.Status4xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "500").Add(float64(stats.Status500))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "501").Add(float64(stats.Status501))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "502").Add(float64(stats.Status502))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "503").Add(float64(stats.Status503))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "504").Add(float64(stats.Status504))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "505").Add(float64(stats.Status505))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, "5xx").Add(float64(stats.Status5xx))
			}
		}
	}
}
