package domain

// Response models the domain inspector real-time data from rt.fastly.com.
type Response struct {
	Data           []Data `json:"Data"`
	Timestamp      uint64 `json:"Timestamp"`
	AggregateDelay int    `json:"AggregateDelay"`
	Error          string `json:"Error"`
}

// Data is the top-level grouping of real-time domain inspector stats.
type Data struct {
	Datacenter ByDatacenter `json:"datacenter"`
	Aggregated ByDomain     `json:"aggregated"`
}

// ByDatacenter groups domain inspector stats by datacenter.
type ByDatacenter map[string]ByDomain

// ByDomain groups domain inspector stats by domain.
type ByDomain map[string]Stats

// Stats for a specific datacenter and domain.
type Stats struct {
	Bandwidth                  uint64  `json:"bandwidth"`                      //	integer	Total bytes delivered (resp_header_bytes + resp_body_bytes + bereq_header_bytes + bereq_body_bytes).
	BereqBodyBytes             uint64  `json:"bereq_body_bytes"`               //	integer	Total body bytes sent to origin.
	BereqHeaderBytes           uint64  `json:"bereq_header_bytes"`             //	integer	Total header bytes sent to origin.
	EdgeHitRatio               float64 `json:"edge_hit_ratio"`                 //	float	Ratio of cache hits to cache misses at the edge, between 0 and 1 (edge_hit_requests / (edge_hit_requests + edge_miss_requests)).
	EdgeHitRequests            uint64  `json:"edge_hit_requests"`              //	integer	Number of requests sent by end users to Fastly that resulted in a hit at the edge.
	EdgeMissRequests           uint64  `json:"edge_miss_requests"`             //	integer	Number of requests sent by end users to Fastly that resulted in a miss at the edge.
	EdgeRequests               uint64  `json:"edge_requests"`                  //	integer	Number of requests sent by end users to Fastly.
	EdgeRespBodyBytes          uint64  `json:"edge_resp_body_bytes"`           //	integer	Total body bytes delivered from Fastly to the end user.
	EdgeRespHeaderBytes        uint64  `json:"edge_resp_header_bytes"`         //	integer	Total header bytes delivered from Fastly to the end user.
	OriginFetchRespBodyBytes   uint64  `json:"origin_fetch_resp_body_bytes"`   //	integer	Total body bytes received from origin.
	OriginFetchRespHeaderBytes uint64  `json:"origin_fetch_resp_header_bytes"` //	integer	Total header bytes received from origin.
	OriginFetches              uint64  `json:"origin_fetches"`                 //	integer	Number of requests sent to origin.
	OriginOffload              float64 `json:"origin_offload"`                 //	float	Ratio of response bytes delivered from the edge compared to what is delivered from origin, between 0 and 1. (edge_resp_body_bytes + edge_resp_header_bytes) / (origin_fetch_resp_body_bytes + origin_fetch_resp_header_bytes + edge_resp_body_bytes + edge_resp_header_bytes).
	OriginStatus1xx            uint64  `json:"origin_status_1xx"`              //	integer	Number of "Informational" category status codes received from origin.
	OriginStatus200            uint64  `json:"origin_status_200"`              //	integer	Number of responses received from origin with status code 200 (Success).
	OriginStatus204            uint64  `json:"origin_status_204"`              //	integer	Number of responses received from origin with status code 204 (No Content).
	OriginStatus206            uint64  `json:"origin_status_206"`              //	integer	Number of responses received from origin with status code 206 (Partial Content).
	OriginStatus2xx            uint64  `json:"origin_status_2xx"`              //	integer	Number of "Success" status codes received from origin.
	OriginStatus301            uint64  `json:"origin_status_301"`              //	integer	Number of responses received from origin with status code 301 (Moved Permanently).
	OriginStatus302            uint64  `json:"origin_status_302"`              //	integer	Number of responses received from origin with status code 302 (Found).
	OriginStatus304            uint64  `json:"origin_status_304"`              //	integer	Number of responses received from origin with status code 304 (Not Modified).
	OriginStatus3xx            uint64  `json:"origin_status_3xx"`              //	integer	Number of "Redirection" codes received from origin.
	OriginStatus400            uint64  `json:"origin_status_400"`              //	integer	Number of responses received from origin with status code 400 (Bad Request).
	OriginStatus401            uint64  `json:"origin_status_401"`              //	integer	Number of responses received from origin with status code 401 (Unauthorized).
	OriginStatus403            uint64  `json:"origin_status_403"`              //	integer	Number of responses received from origin with status code 403 (Forbidden).
	OriginStatus404            uint64  `json:"origin_status_404"`              //	integer	Number of responses received from origin with status code 404 (Not Found).
	OriginStatus416            uint64  `json:"origin_status_416"`              //	integer	Number of responses received from origin with status code 416 (Range Not Satisfiable).
	OriginStatus429            uint64  `json:"origin_status_429"`              //	integer	Number of responses received from origin with status code 429 (Too Many Requests).
	OriginStatus4xx            uint64  `json:"origin_status_4xx"`              //	integer	Number of "Client Error" codes received from origin.
	OriginStatus500            uint64  `json:"origin_status_500"`              //	integer	Number of responses received from origin with status code 500 (Internal Server Error).
	OriginStatus501            uint64  `json:"origin_status_501"`              //	integer	Number of responses received from origin with status code 501 (Not Implemented).
	OriginStatus502            uint64  `json:"origin_status_502"`              //	integer	Number of responses received from origin with status code 502 (Bad Gateway).
	OriginStatus503            uint64  `json:"origin_status_503"`              //	integer	Number of responses received from origin with status code 503 (Service Unavailable).
	OriginStatus504            uint64  `json:"origin_status_504"`              //	integer	Number of responses received from origin with status code 504 (Gateway Timeout).
	OriginStatus505            uint64  `json:"origin_status_505"`              //	integer	Number of responses received from origin with status code 505 (HTTP Version Not Supported).
	OriginStatus5xx            uint64  `json:"origin_status_5xx"`              //	integer	Number of "Server Error" codes received from origin.
	Requests                   uint64  `json:"requests"`                       //	integer	Number of requests processed.
	RespBodyBytes              uint64  `json:"resp_body_bytes"`                //	integer	Total body bytes delivered.
	RespHeaderBytes            uint64  `json:"resp_header_bytes"`              //	integer	Total header bytes delivered.
	Status1xx                  uint64  `json:"status_1xx"`                     //	integer	Number of 1xx "Informational" category status codes delivered.
	Status200                  uint64  `json:"status_200"`                     //	integer	Number of responses received with status code 200 (Success).
	Status204                  uint64  `json:"status_204"`                     //	integer	Number of responses received with status code 204 (No Content).
	Status206                  uint64  `json:"status_206"`                     //	integer	Number of responses received with status code 206 (Partial Content).
	Status2xx                  uint64  `json:"status_2xx"`                     //	integer	Number of 2xx "Success" status codes delivered.
	Status301                  uint64  `json:"status_301"`                     //	integer	Number of responses received with status code 301 (Moved Permanently).
	Status302                  uint64  `json:"status_302"`                     //	integer	Number of responses received with status code 302 (Found).
	Status304                  uint64  `json:"status_304"`                     //	integer	Number of responses received with status code 304 (Not Modified).
	Status3xx                  uint64  `json:"status_3xx"`                     //	integer	Number of 3xx "Redirection" codes delivered.
	Status400                  uint64  `json:"status_400"`                     //	integer	Number of responses received with status code 400 (Bad Request).
	Status401                  uint64  `json:"status_401"`                     //	integer	Number of responses received with status code 401 (Unauthorized).
	Status403                  uint64  `json:"status_403"`                     //	integer	Number of responses received with status code 403 (Forbidden).
	Status404                  uint64  `json:"status_404"`                     //	integer	Number of responses received with status code 404 (Not Found).
	Status416                  uint64  `json:"status_416"`                     //	integer	Number of responses received with status code 416 (Range Not Satisfiable).
	Status429                  uint64  `json:"status_429"`                     //	integer	Number of responses received with status code 429 (Too Many Requests).
	Status4xx                  uint64  `json:"status_4xx"`                     //	integer	Number of 4xx "Client Error" codes delivered.
	Status500                  uint64  `json:"status_500"`                     //	integer	Number of responses received with status code 500 (Internal Server Error).
	Status501                  uint64  `json:"status_501"`                     //	integer	Number of responses received with status code 501 (Not Implemented).
	Status502                  uint64  `json:"status_502"`                     //	integer	Number of responses received with status code 502 (Bad Gateway).
	Status503                  uint64  `json:"status_503"`                     //	integer	Number of responses received with status code 503 (Service Unavailable).
	Status504                  uint64  `json:"status_504"`                     //	integer	Number of responses received with status code 504 (Gateway Timeout).
	Status505                  uint64  `json:"status_505"`                     //	integer	Number of responses received with status code 505 (HTTP Version Not Supported).
	Status5xx                  uint64  `json:"status_5xx"`                     //	integer	Number of 5xx "Server Error" codes delivered.
}
