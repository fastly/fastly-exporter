package origin

const (
	srcDelivery = "delivery"
	srcCompute  = "compute"
	srcWaf      = "waf"
)

// Process updates the metrics with data from the API response.
func Process(response *Response, serviceID, serviceName, serviceVersion string, m *Metrics) {
	for _, d := range response.Data {
		for datacenter, byOrigin := range d.Datacenter {
			for origin, stats := range byOrigin {
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery).Add(float64(stats.RespBodyBytes))
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute).Add(float64(stats.ComputeRespBodyBytes))
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf).Add(float64(stats.WafRespBodyBytes))
				m.RespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery).Add(float64(stats.RespHeaderBytes))
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute).Add(float64(stats.ComputeRespHeaderBytes))
				m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf).Add(float64(stats.WafRespHeaderBytes))
				m.ResponsesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery).Add(float64(stats.Responses))
				m.ResponsesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute).Add(float64(stats.ComputeResponses))
				m.ResponsesTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf).Add(float64(stats.WafResponses))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "1xx").Add(float64(stats.Status1xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "200").Add(float64(stats.Status200))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "204").Add(float64(stats.Status204))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "206").Add(float64(stats.Status206))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "2xx").Add(float64(stats.Status2xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "301").Add(float64(stats.Status301))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "302").Add(float64(stats.Status302))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "304").Add(float64(stats.Status304))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "3xx").Add(float64(stats.Status3xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "400").Add(float64(stats.Status400))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "401").Add(float64(stats.Status401))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "403").Add(float64(stats.Status403))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "404").Add(float64(stats.Status404))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "416").Add(float64(stats.Status416))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "429").Add(float64(stats.Status429))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "4xx").Add(float64(stats.Status4xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "500").Add(float64(stats.Status500))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "501").Add(float64(stats.Status501))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "502").Add(float64(stats.Status502))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "503").Add(float64(stats.Status503))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "504").Add(float64(stats.Status504))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "505").Add(float64(stats.Status505))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcDelivery, "5xx").Add(float64(stats.Status5xx))

				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "1xx").Add(float64(stats.ComputeStatus1xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "200").Add(float64(stats.ComputeStatus200))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "204").Add(float64(stats.ComputeStatus204))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "2xx").Add(float64(stats.ComputeStatus2xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "301").Add(float64(stats.ComputeStatus301))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "302").Add(float64(stats.ComputeStatus302))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "304").Add(float64(stats.ComputeStatus304))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "3xx").Add(float64(stats.ComputeStatus3xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "400").Add(float64(stats.ComputeStatus400))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "401").Add(float64(stats.ComputeStatus401))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "403").Add(float64(stats.ComputeStatus403))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "404").Add(float64(stats.ComputeStatus404))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "416").Add(float64(stats.ComputeStatus416))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "4xx").Add(float64(stats.ComputeStatus4xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "500").Add(float64(stats.ComputeStatus500))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "501").Add(float64(stats.ComputeStatus501))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "502").Add(float64(stats.ComputeStatus502))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "503").Add(float64(stats.ComputeStatus503))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "504").Add(float64(stats.ComputeStatus504))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "505").Add(float64(stats.ComputeStatus505))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcCompute, "5xx").Add(float64(stats.ComputeStatus5xx))

				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "1xx").Add(float64(stats.WafStatus1xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "200").Add(float64(stats.WafStatus200))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "204").Add(float64(stats.WafStatus204))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "206").Add(float64(stats.WafStatus206))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "2xx").Add(float64(stats.WafStatus2xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "301").Add(float64(stats.WafStatus301))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "302").Add(float64(stats.WafStatus302))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "304").Add(float64(stats.WafStatus304))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "3xx").Add(float64(stats.WafStatus3xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "400").Add(float64(stats.WafStatus400))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "401").Add(float64(stats.WafStatus401))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "403").Add(float64(stats.WafStatus403))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "404").Add(float64(stats.WafStatus404))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "416").Add(float64(stats.WafStatus416))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "429").Add(float64(stats.WafStatus429))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "4xx").Add(float64(stats.WafStatus4xx))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "500").Add(float64(stats.WafStatus500))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "501").Add(float64(stats.WafStatus501))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "502").Add(float64(stats.WafStatus502))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "503").Add(float64(stats.WafStatus503))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "504").Add(float64(stats.WafStatus504))
				m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "505").Add(float64(stats.WafStatus505))
				m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, origin, srcWaf, "5xx").Add(float64(stats.WafStatus5xx))

				// Latency stats are clearly from xxx_bucket{le="v"} metrics,
				// but I don't see a good way to re-populate a histogram from
				// those numbers. (If I'm missing something, file an issue!)
				//
				// Our clue is the final bucket, which says it's observations
				// "of 60s and above". Based on that we use the lower bound of
				// each stat as the observed value, except for the first bucket
				// which we yolo as 500us because 0 doesn't really make sense??
				for v, n := range map[float64]uint64{
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
					for i := uint64(0); i < n; i++ {
						m.LatencySeconds.WithLabelValues(serviceID, serviceName, datacenter, origin).Observe(v)
					}
				}

				for v, n := range map[float64]uint64{
					60.00:  stats.WafLatency60000plus,
					10.00:  stats.WafLatency10000to60000,
					5.000:  stats.WafLatency5000to10000,
					1.000:  stats.WafLatency1000to5000,
					0.500:  stats.WafLatency500to1000,
					0.250:  stats.WafLatency250to500,
					0.100:  stats.WafLatency100to250,
					0.050:  stats.WafLatency50to100,
					0.010:  stats.WafLatency10to50,
					0.005:  stats.WafLatency5to10,
					0.001:  stats.WafLatency1to5,
					0.0005: stats.WafLatency0to1, // yolo
				} {
					for i := uint64(0); i < n; i++ {
						m.LatencySeconds.WithLabelValues(serviceID, serviceName, datacenter, origin).Observe(v)
					}
				}

				for v, n := range map[float64]uint64{
					60.00:  stats.ComputeLatency60000plus,
					10.00:  stats.ComputeLatency10000to60000,
					5.000:  stats.ComputeLatency5000to10000,
					1.000:  stats.ComputeLatency1000to5000,
					0.500:  stats.ComputeLatency500to1000,
					0.250:  stats.ComputeLatency250to500,
					0.100:  stats.ComputeLatency100to250,
					0.050:  stats.ComputeLatency50to100,
					0.010:  stats.ComputeLatency10to50,
					0.005:  stats.ComputeLatency5to10,
					0.001:  stats.ComputeLatency1to5,
					0.0005: stats.ComputeLatency0to1, // yolo
				} {
					for i := uint64(0); i < n; i++ {
						m.LatencySeconds.WithLabelValues(serviceID, serviceName, datacenter, origin).Observe(v)
					}
				}
			}
		}
	}
}
