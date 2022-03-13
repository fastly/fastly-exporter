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

				// Latency stats are clearly from xxx_bucket{le="v"} metrics,
				// but I don't see a good way to re-populate a histogram from
				// those numbers. (If I'm missing something, file an issue!)
				//
				// Our clue is the final bucket, which says it's observations
				// "of 60s and above". Based on that we use the lower bound of
				// each stat as the observed value, except for the first bucket
				// which we yolo as 500us because 0 doesn't really make sense??
				for v, n := range map[float64]int{
					60.00:  stats.Latency60000plus,
					10.00:  stats.Latency10000to60000,
					5.000:  stats.Latency5000to10000,
					1.000:  stats.Latency1000to5000,
					0.500:  stats.Latency500to1000,
					0.250:  stats.Latency250to500,
					0.100:  stats.Latency100to250,
					0.050:  stats.Latency50to100,
					0.010:  stats.Latency10to50,
					0.005:  stats.Latency5to10,
					0.001:  stats.Latency1to5,
					0.0005: stats.Latency0to1, // yolo
				} {
					for i := 0; i < n; i++ {
						m.LatencySeconds.WithLabelValues(serviceID, serviceName, datacenter, origin).Observe(v)
					}
				}
			}
		}
	}
}
