package origin

// Response models the origin inspector real-time data from rt.fastly.com.
type Response struct {
	Data           []Data `json:"Data"`
	Timestamp      uint64 `json:"Timestamp"`
	AggregateDelay int    `json:"AggregateDelay"`
	Error          string `json:"Error"`
}

// Data is the top-level grouping of real-time origin inspector stats.
type Data struct {
	Datacenter ByDatacenter `json:"datacenter"`
	Aggregated ByOrigin     `json:"aggregated"`
}

// ByDatacenter groups origin inspector stats by datacenter.
type ByDatacenter map[string]ByOrigin

// ByOrigin groups origin inspector stats by origin.
type ByOrigin map[string]Stats

// Stats for a specific datacenter and origin.
type Stats struct {
	RespBodyBytes       uint64 `json:"resp_body_bytes"`          // Number of body bytes from origin.
	RespHeaderBytes     uint64 `json:"resp_header_bytes"`        // Number of header bytes from origin.
	Responses           uint64 `json:"responses"`                // Number of responses from origin.
	Status1xx           uint64 `json:"status_1xx"`               // Number of 1xx "Informational" category status codes delivered from origin.
	Status200           uint64 `json:"status_200"`               // Number of responses received with status code 200 (Success) from origin.
	Status204           uint64 `json:"status_204"`               // Number of responses received with status code 204 (No Content) from origin.
	Status2xx           uint64 `json:"status_2xx"`               // Number of 2xx "Success" status codes delivered from origin.
	Status301           uint64 `json:"status_301"`               // Number of responses received with status code 301 (Moved Permanently) from origin.
	Status302           uint64 `json:"status_302"`               // Number of responses received with status code 302 (Found) from origin.
	Status304           uint64 `json:"status_304"`               // Number of responses received with status code 304 (Not Modified) from origin.
	Status3xx           uint64 `json:"status_3xx"`               // Number of 3xx "Redirection" codes delivered from origin.
	Status400           uint64 `json:"status_400"`               // Number of responses received with status code 400 (Bad Request) from origin.
	Status401           uint64 `json:"status_401"`               // Number of responses received with status code 401 (Unauthorized) from origin.
	Status403           uint64 `json:"status_403"`               // Number of responses received with status code 403 (Forbidden) from origin.
	Status404           uint64 `json:"status_404"`               // Number of responses received with status code 404 (Not Found) from origin.
	Status416           uint64 `json:"status_416"`               // Number of responses received with status code 416 (Range Not Satisfiable) from origin.
	Status4xx           uint64 `json:"status_4xx"`               // Number of 4xx "Client Error" codes delivered from origin.
	Status500           uint64 `json:"status_500"`               // Number of responses received with status code 500 (Internal Server Error) from origin.
	Status501           uint64 `json:"status_501"`               // Number of responses received with status code 501 (Not Implemented) from origin.
	Status502           uint64 `json:"status_502"`               // Number of responses received with status code 502 (Bad Gateway) from origin.
	Status503           uint64 `json:"status_503"`               // Number of responses received with status code 503 (Service Unavailable) from origin.
	Status504           uint64 `json:"status_504"`               // Number of responses received with status code 504 (Gateway Timeout) from origin.
	Status505           uint64 `json:"status_505"`               // Number of responses received with status code 505 (HTTP Version Not Supported) from origin.
	Status5xx           uint64 `json:"status_5xx"`               // Number of 5xx "Server Error" codes delivered from origin.
	Latency0to1         uint64 `json:"latency_0_to_1ms"`         // Number of responses from origin with latency between 0 and 1 millisecond.
	Latency1to5         uint64 `json:"latency_1_to_5ms"`         // Number of responses from origin with latency between 1 and 5 milliseconds.
	Latency5to10        uint64 `json:"latency_5_to_10ms"`        // Number of responses from origin with latency between 5 and 10 milliseconds.
	Latency10to50       uint64 `json:"latency_10_to_50ms"`       // Number of responses from origin with latency between 10 and 50 milliseconds.
	Latency50to100      uint64 `json:"latency_50_to_100ms"`      // Number of responses from origin with latency between 50 and 100 milliseconds.
	Latency100to250     uint64 `json:"latency_100_to_250ms"`     // Number of responses from origin with latency between 100 and 250 milliseconds.
	Latency250to500     uint64 `json:"latency_250_to_500ms"`     // Number of responses from origin with latency between 250 and 500 milliseconds.
	Latency500to1000    uint64 `json:"latency_500_to_1000ms"`    // Number of responses from origin with latency between 500 and 1,000 milliseconds.
	Latency1000to5000   uint64 `json:"latency_1000_to_5000ms"`   // Number of responses from origin with latency between 1,000 and 5,000 milliseconds.
	Latency5000to10000  uint64 `json:"latency_5000_to_10000ms"`  // Number of responses from origin with latency between 5,000 and 10,000 milliseconds.
	Latency10000to60000 uint64 `json:"latency_10000_to_60000ms"` // Number of responses from origin with latency between 10,000 and 60,000 milliseconds.
	Latency60000plus    uint64 `json:"latency_60000ms"`          // Number of responses from origin with latency of 60,000 milliseconds and above.
}
