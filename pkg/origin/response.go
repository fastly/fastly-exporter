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
	RespBodyBytes              uint64 `json:"resp_body_bytes"`                  // Number of body bytes from origin.
	RespHeaderBytes            uint64 `json:"resp_header_bytes"`                // Number of header bytes from origin.
	Responses                  uint64 `json:"responses"`                        // Number of responses from origin.
	Status1xx                  uint64 `json:"status_1xx"`                       // Number of 1xx "Informational" status codes delivered from origin.
	Status200                  uint64 `json:"status_200"`                       // Number of responses received with status code 200 (Success) from origin.
	Status204                  uint64 `json:"status_204"`                       // Number of responses received with status code 204 (No Content) from origin.
	Status206                  uint64 `json:"status_206"`                       // Number of responses received with status code 206 (Partial Content) from origin.
	Status2xx                  uint64 `json:"status_2xx"`                       // Number of 2xx "Success" status codes delivered from origin.
	Status301                  uint64 `json:"status_301"`                       // Number of responses received with status code 301 (Moved Permanently) from origin.
	Status302                  uint64 `json:"status_302"`                       // Number of responses received with status code 302 (Found) from origin.
	Status304                  uint64 `json:"status_304"`                       // Number of responses received with status code 304 (Not Modified) from origin.
	Status3xx                  uint64 `json:"status_3xx"`                       // Number of 3xx "Redirection" codes delivered from origin.
	Status400                  uint64 `json:"status_400"`                       // Number of responses received with status code 400 (Bad Request) from origin.
	Status401                  uint64 `json:"status_401"`                       // Number of responses received with status code 401 (Unauthorized) from origin.
	Status403                  uint64 `json:"status_403"`                       // Number of responses received with status code 403 (Forbidden) from origin.
	Status404                  uint64 `json:"status_404"`                       // Number of responses received with status code 404 (Not Found) from origin.
	Status416                  uint64 `json:"status_416"`                       // Number of responses received with status code 416 (Range Not Satisfiable) from origin.
	Status429                  uint64 `json:"status_429"`                       // Number of responses received with status code 429 (Too Many Requests) from origin.
	Status4xx                  uint64 `json:"status_4xx"`                       // Number of 4xx "Client Error" codes delivered from origin.
	Status500                  uint64 `json:"status_500"`                       // Number of responses received with status code 500 (Internal Server Error) from origin.
	Status501                  uint64 `json:"status_501"`                       // Number of responses received with status code 501 (Not Implemented) from origin.
	Status502                  uint64 `json:"status_502"`                       // Number of responses received with status code 502 (Bad Gateway) from origin.
	Status503                  uint64 `json:"status_503"`                       // Number of responses received with status code 503 (Service Unavailable) from origin.
	Status504                  uint64 `json:"status_504"`                       // Number of responses received with status code 504 (Gateway Timeout) from origin.
	Status505                  uint64 `json:"status_505"`                       // Number of responses received with status code 505 (HTTP Version Not Supported) from origin.
	Status530                  uint64 `json:"status_530"`                       // Number of responses received from origin with status code 530.
	Status5xx                  uint64 `json:"status_5xx"`                       // Number of 5xx "Server Error" codes delivered from origin.
	Latency0to1                uint64 `json:"latency_0_to_1ms"`                 // Number of responses from origin with latency between 0 and 1 millisecond.
	Latency1to5                uint64 `json:"latency_1_to_5ms"`                 // Number of responses from origin with latency between 1 and 5 milliseconds.
	Latency5to10               uint64 `json:"latency_5_to_10ms"`                // Number of responses from origin with latency between 5 and 10 milliseconds.
	Latency10to50              uint64 `json:"latency_10_to_50ms"`               // Number of responses from origin with latency between 10 and 50 milliseconds.
	Latency50to100             uint64 `json:"latency_50_to_100ms"`              // Number of responses from origin with latency between 50 and 100 milliseconds.
	Latency100to250            uint64 `json:"latency_100_to_250ms"`             // Number of responses from origin with latency between 100 and 250 milliseconds.
	Latency250to500            uint64 `json:"latency_250_to_500ms"`             // Number of responses from origin with latency between 250 and 500 milliseconds.
	Latency500to1000           uint64 `json:"latency_500_to_1000ms"`            // Number of responses from origin with latency between 500 and 1,000 milliseconds.
	Latency1000to5000          uint64 `json:"latency_1000_to_5000ms"`           // Number of responses from origin with latency between 1,000 and 5,000 milliseconds.
	Latency5000to10000         uint64 `json:"latency_5000_to_10000ms"`          // Number of responses from origin with latency between 5,000 and 10,000 milliseconds.
	Latency10000to60000        uint64 `json:"latency_10000_to_60000ms"`         // Number of responses from origin with latency between 10,000 and 60,000 milliseconds.
	Latency60000plus           uint64 `json:"latency_60000ms"`                  // Number of responses from origin with latency of 60,000 milliseconds and above.
	WafResponses               uint64 `json:"waf_responses"`                    // Number of responses received for origin requests made by the Fastly WAF.
	WafRespHeaderBytes         uint64 `json:"waf_resp_header_bytes"`            // Number of header bytes received for origin requests made by the Fastly WAF.
	WafRespBodyBytes           uint64 `json:"waf_resp_body_bytes"`              // Number of body bytes received for origin requests made by the Fastly WAF.
	WafStatus1xx               uint64 `json:"waf_status_1xx"`                   // Number of 1xx "Informational" status codes received for origin requests made by the Fastly WAF.
	WafStatus2xx               uint64 `json:"waf_status_2xx"`                   // Number of 2xx "Success" status codes received for origin requests made by the Fastly WAF.
	WafStatus3xx               uint64 `json:"waf_status_3xx"`                   // Number of 3xx "Redirection" codes received for origin requests made by the Fastly WAF.
	WafStatus4xx               uint64 `json:"waf_status_4xx"`                   // Number of 4xx "Client Error" codes received for origin requests made by the Fastly WAF.
	WafStatus5xx               uint64 `json:"waf_status_5xx"`                   // Number of 5xx "Server Error" codes received for origin requests made by the Fastly WAF.
	WafStatus200               uint64 `json:"waf_status_200"`                   // Number of responses received with status code 200 (Success) received for origin requests made by the Fastly WAF.
	WafStatus204               uint64 `json:"waf_status_204"`                   // Number of responses received with status code 204 (No Content) received for origin requests made by the Fastly WAF.
	WafStatus206               uint64 `json:"waf_status_206"`                   // Number of responses received with status code 206 (Partial Content) received for origin requests made by the Fastly WAF.
	WafStatus301               uint64 `json:"waf_status_301"`                   // Number of responses received with status code 301 (Moved Permanently) received for origin requests made by the Fastly WAF.
	WafStatus302               uint64 `json:"waf_status_302"`                   // Number of responses received with status code 302 (Found) received for origin requests made by the Fastly WAF.
	WafStatus304               uint64 `json:"waf_status_304"`                   // Number of responses received with status code 304 (Not Modified) received for origin requests made by the Fastly WAF.
	WafStatus400               uint64 `json:"waf_status_400"`                   // Number of responses received with status code 400 (Bad Request) received for origin requests made by the Fastly WAF.
	WafStatus401               uint64 `json:"waf_status_401"`                   // Number of responses received with status code 401 (Unauthorized) received for origin requests made by the Fastly WAF.
	WafStatus403               uint64 `json:"waf_status_403"`                   // Number of responses received with status code 403 (Forbidden) received for origin requests made by the Fastly WAF.
	WafStatus404               uint64 `json:"waf_status_404"`                   // Number of responses received with status code 404 (Not Found) received for origin requests made by the Fastly WAF.
	WafStatus416               uint64 `json:"waf_status_416"`                   // Number of responses received with status code 416 (Range Not Satisfiable) received for origin requests made by the Fastly WAF.
	WafStatus429               uint64 `json:"waf_status_429"`                   // Number of responses received with status code 429 (Too Many Requests) received for origin requests made by the Fastly WAF.
	WafStatus500               uint64 `json:"waf_status_500"`                   // Number of responses received with status code 500 (Internal Server Error) received for origin requests made by the Fastly WAF.
	WafStatus501               uint64 `json:"waf_status_501"`                   // Number of responses received with status code 501 (Not Implemented) received for origin requests made by the Fastly WAF.
	WafStatus502               uint64 `json:"waf_status_502"`                   // Number of responses received with status code 502 (Bad Gateway) received for origin requests made by the Fastly WAF.
	WafStatus503               uint64 `json:"waf_status_503"`                   // Number of responses received with status code 503 (Service Unavailable) received for origin requests made by the Fastly WAF.
	WafStatus504               uint64 `json:"waf_status_504"`                   // Number of responses received with status code 504 (Gateway Timeout) received for origin requests made by the Fastly WAF.
	WafStatus505               uint64 `json:"waf_status_505"`                   // Number of responses received with status code 505 (HTTP Version Not Supported) received for origin requests made by the Fastly WAF.
	WafStatus530               uint64 `json:"waf_status_530"`                   // Number of responses received with status code 530 received for origin requests made by the Fastly WAF.
	WafLatency0to1             uint64 `json:"waf_latency_0_to_1ms"`             // Number of responses with latency between 0 and 1 millisecond received for origin requests made by the Fastly WAF.
	WafLatency1to5             uint64 `json:"waf_latency_1_to_5ms"`             // Number of responses with latency between 1 and 5 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency5to10            uint64 `json:"waf_latency_5_to_10ms"`            // Number of responses with latency between 5 and 10 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency10to50           uint64 `json:"waf_latency_10_to_50ms"`           // Number of responses with latency between 10 and 50 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency50to100          uint64 `json:"waf_latency_50_to_100ms"`          // Number of responses with latency between 50 and 100 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency100to250         uint64 `json:"waf_latency_100_to_250ms"`         // Number of responses with latency between 100 and 250 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency250to500         uint64 `json:"waf_latency_250_to_500ms"`         // Number of responses with latency between 250 and 500 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency500to1000        uint64 `json:"waf_latency_500_to_1000ms"`        // Number of responses with latency between 500 and 1,000 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency1000to5000       uint64 `json:"waf_latency_1000_to_5000ms"`       // Number of responses with latency between 1,000 and 5,000 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency5000to10000      uint64 `json:"waf_latency_5000_to_10000ms"`      // Number of responses with latency between 5,000 and 10,000 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency10000to60000     uint64 `json:"waf_latency_10000_to_60000ms"`     // Number of responses with latency between 10,000 and 60,000 milliseconds received for origin requests made by the Fastly WAF.
	WafLatency60000plus        uint64 `json:"waf_latency_60000ms"`              // Number of responses with latency of 60,000 milliseconds and above received for origin requests made by the Fastly WAF.
	ComputeResponses           uint64 `json:"compute_responses"`                // Number of responses for origin received by Compute@Edge.
	ComputeRespHeaderBytes     uint64 `json:"compute_resp_header_bytes"`        // Number of header bytes for origin received by Compute@Edge.
	ComputeRespBodyBytes       uint64 `json:"compute_resp_body_bytes"`          // Number of body bytes for origin received by Compute@Edge.
	ComputeStatus1xx           uint64 `json:"compute_status_1xx"`               // Number of 1xx "Informational" status codes for origin received by Compute@Edge.
	ComputeStatus2xx           uint64 `json:"compute_status_2xx"`               // Number of 2xx "Success" status codes for origin received by Compute@Edge.
	ComputeStatus3xx           uint64 `json:"compute_status_3xx"`               // Number of 3xx "Redirection" codes for origin received by Compute@Edge.
	ComputeStatus4xx           uint64 `json:"compute_status_4xx"`               // Number of 4xx "Client Error" codes for origin received by Compute@Edge.
	ComputeStatus5xx           uint64 `json:"compute_status_5xx"`               // Number of 5xx "Server Error" codes for origin received by Compute@Edge.
	ComputeStatus200           uint64 `json:"compute_status_200"`               // Number of responses received with status code 200 (Success) for origin received by Compute@Edge.
	ComputeStatus204           uint64 `json:"compute_status_204"`               // Number of responses received with status code 204 (No Content) for origin received by Compute@Edge.
	ComputeStatus206           uint64 `json:"compute_status_206"`               // Number of responses received with status code 206 (Partial Content) for origin received by Compute@Edge.
	ComputeStatus301           uint64 `json:"compute_status_301"`               // Number of responses received with status code 301 (Moved Permanently) for origin received by Compute@Edge.
	ComputeStatus302           uint64 `json:"compute_status_302"`               // Number of responses received with status code 302 (Found) for origin received by Compute@Edge.
	ComputeStatus304           uint64 `json:"compute_status_304"`               // Number of responses received with status code 304 (Not Modified) for origin received by Compute@Edge.
	ComputeStatus400           uint64 `json:"compute_status_400"`               // Number of responses received with status code 400 (Bad Request) for origin received by Compute@Edge.
	ComputeStatus401           uint64 `json:"compute_status_401"`               // Number of responses received with status code 401 (Unauthorized) for origin received by Compute@Edge.
	ComputeStatus403           uint64 `json:"compute_status_403"`               // Number of responses received with status code 403 (Forbidden) for origin received by Compute@Edge.
	ComputeStatus404           uint64 `json:"compute_status_404"`               // Number of responses received with status code 404 (Not Found) for origin received by Compute@Edge.
	ComputeStatus416           uint64 `json:"compute_status_416"`               // Number of responses received with status code 416 (Range Not Satisfiable) for origin received by Compute@Edge.
	ComputeStatus429           uint64 `json:"compute_status_429"`               // Number of responses received with status code 429 (Too Many Requests) for origin received by Compute@Edge.
	ComputeStatus500           uint64 `json:"compute_status_500"`               // Number of responses received with status code 500 (Internal Server Error) for origin received by Compute@Edge.
	ComputeStatus501           uint64 `json:"compute_status_501"`               // Number of responses received with status code 501 (Not Implemented) for origin received by Compute@Edge.
	ComputeStatus502           uint64 `json:"compute_status_502"`               // Number of responses received with status code 502 (Bad Gateway) for origin received by Compute@Edge.
	ComputeStatus503           uint64 `json:"compute_status_503"`               // Number of responses received with status code 503 (Service Unavailable) for origin received by Compute@Edge.
	ComputeStatus504           uint64 `json:"compute_status_504"`               // Number of responses received with status code 504 (Gateway Timeout) for origin received by Compute@Edge.
	ComputeStatus505           uint64 `json:"compute_status_505"`               // Number of responses received with status code 505 (HTTP Version Not Supported) for origin received by Compute@Edge.
	ComputeStatus530           uint64 `json:"compute_status_530"`               // Number of responses received with status code 530 for origin received by the Compute platform.
	ComputeLatency0to1         uint64 `json:"compute_latency_0_to_1ms"`         // Number of responses with latency between 0 and 1 millisecond for origin received by Compute@Edge.
	ComputeLatency1to5         uint64 `json:"compute_latency_1_to_5ms"`         // Number of responses with latency between 1 and 5 milliseconds for origin received by Compute@Edge.
	ComputeLatency5to10        uint64 `json:"compute_latency_5_to_10ms"`        // Number of responses with latency between 5 and 10 milliseconds for origin received by Compute@Edge.
	ComputeLatency10to50       uint64 `json:"compute_latency_10_to_50ms"`       // Number of responses with latency between 10 and 50 milliseconds for origin received by Compute@Edge.
	ComputeLatency50to100      uint64 `json:"compute_latency_50_to_100ms"`      // Number of responses with latency between 50 and 100 milliseconds for origin received by Compute@Edge.
	ComputeLatency100to250     uint64 `json:"compute_latency_100_to_250ms"`     // Number of responses with latency between 100 and 250 milliseconds for origin received by Compute@Edge.
	ComputeLatency250to500     uint64 `json:"compute_latency_250_to_500ms"`     // Number of responses with latency between 250 and 500 milliseconds for origin received by Compute@Edge.
	ComputeLatency500to1000    uint64 `json:"compute_latency_500_to_1000ms"`    // Number of responses with latency between 500 and 1,000 milliseconds for origin received by Compute@Edge.
	ComputeLatency1000to5000   uint64 `json:"compute_latency_1000_to_5000ms"`   // Number of responses with latency between 1,000 and 5,000 milliseconds for origin received by Compute@Edge.
	ComputeLatency5000to10000  uint64 `json:"compute_latency_5000_to_10000ms"`  // Number of responses with latency between 5,000 and 10,000 milliseconds for origin received by Compute@Edge.
	ComputeLatency10000to60000 uint64 `json:"compute_latency_10000_to_60000ms"` // Number of responses with latency between 10,000 and 60,000 milliseconds for origin received by Compute@Edge.
	ComputeLatency60000plus    uint64 `json:"compute_latency_60000ms"`          // Number of responses with latency of 60,000 milliseconds and above for origin received by Compute@Edge.
}
