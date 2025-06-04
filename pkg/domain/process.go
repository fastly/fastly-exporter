package domain

// Process updates the metrics with data from the API response.
func Process(response *Response, serviceID, serviceName, _ string, m *Metrics, aggregateOnly bool) {
	const aggregateDC = "aggregate"

	for _, d := range response.Data {
		if aggregateOnly {
			for domain, stats := range d.Aggregated {
				process(serviceID, serviceName, aggregateDC, domain, stats, m)
			}

			continue
		}

		for datacenter, byDomain := range d.Datacenter {
			for domain, stats := range byDomain {
				process(serviceID, serviceName, datacenter, domain, stats, m)
			}
		}
	}
}

func process(serviceID, serviceName, datacenter, domain string, stats Stats, m *Metrics) {
	m.BackendReqBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.BereqBodyBytes))
	m.BackendReqHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.BereqHeaderBytes))
	m.EdgeHitRatio.WithLabelValues(serviceID, serviceName, datacenter, domain).Set(stats.EdgeHitRatio)
	m.EdgeHitRequestsTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.EdgeHitRequests))
	m.EdgeMissRequestsTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.EdgeMissRequests))
	m.EdgeRequestsTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.EdgeRequests))
	m.EdgeResponseBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.EdgeRespBodyBytes))
	m.EdgeResponseHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.EdgeRespHeaderBytes))
	m.OriginFetchRespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.OriginFetchRespBodyBytes))
	m.OriginFetchRespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.OriginFetchRespHeaderBytes))
	m.OriginFetches.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.OriginFetches))
	m.OriginOffload.WithLabelValues(serviceID, serviceName, datacenter, domain).Set(stats.OriginOffload)
	m.OriginStatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "1xx").Add(float64(stats.OriginStatus1xx))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "200").Add(float64(stats.OriginStatus200))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "204").Add(float64(stats.OriginStatus204))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "206").Add(float64(stats.OriginStatus206))
	m.OriginStatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "2xx").Add(float64(stats.OriginStatus2xx))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "301").Add(float64(stats.OriginStatus301))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "302").Add(float64(stats.OriginStatus302))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "304").Add(float64(stats.OriginStatus304))
	m.OriginStatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "3xx").Add(float64(stats.OriginStatus3xx))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "400").Add(float64(stats.OriginStatus400))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "401").Add(float64(stats.OriginStatus401))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "403").Add(float64(stats.OriginStatus403))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "404").Add(float64(stats.OriginStatus404))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "416").Add(float64(stats.OriginStatus416))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "429").Add(float64(stats.OriginStatus429))
	m.OriginStatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "4xx").Add(float64(stats.OriginStatus4xx))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "500").Add(float64(stats.OriginStatus500))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "501").Add(float64(stats.OriginStatus501))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "502").Add(float64(stats.OriginStatus502))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "503").Add(float64(stats.OriginStatus503))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "504").Add(float64(stats.OriginStatus504))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "505").Add(float64(stats.OriginStatus505))
	m.OriginStatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "530").Add(float64(stats.OriginStatus530))
	m.OriginStatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "5xx").Add(float64(stats.OriginStatus5xx))
	m.RequestsTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.Requests))
	m.RespBodyBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.RespBodyBytes))
	m.RespHeaderBytesTotal.WithLabelValues(serviceID, serviceName, datacenter, domain).Add(float64(stats.RespHeaderBytes))
	m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "1xx").Add(float64(stats.Status1xx))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "200").Add(float64(stats.Status200))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "204").Add(float64(stats.Status204))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "206").Add(float64(stats.Status206))
	m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "2xx").Add(float64(stats.Status2xx))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "301").Add(float64(stats.Status301))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "302").Add(float64(stats.Status302))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "304").Add(float64(stats.Status304))
	m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "3xx").Add(float64(stats.Status3xx))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "400").Add(float64(stats.Status400))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "401").Add(float64(stats.Status401))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "403").Add(float64(stats.Status403))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "404").Add(float64(stats.Status404))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "416").Add(float64(stats.Status416))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "429").Add(float64(stats.Status429))
	m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "4xx").Add(float64(stats.Status4xx))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "500").Add(float64(stats.Status500))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "501").Add(float64(stats.Status501))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "502").Add(float64(stats.Status502))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "503").Add(float64(stats.Status503))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "504").Add(float64(stats.Status504))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "505").Add(float64(stats.Status505))
	m.StatusCodeTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "530").Add(float64(stats.Status530))
	m.StatusGroupTotal.WithLabelValues(serviceID, serviceName, datacenter, domain, "5xx").Add(float64(stats.Status5xx))
}
