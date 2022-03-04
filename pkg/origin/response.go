package origin

type Response struct {
	Data           []Data `json:"Data"`
	Timestamp      uint64 `json:"Timestamp"`
	AggregateDelay int    `json:"AggregateDelay"`
	Error          string `json:"Error"`
}

type Data struct {
	Datacenter ByDatacenter `json:"datacenter"`
	Aggregated ByOrigin     `json:"aggregated"`
}

type ByDatacenter map[string]ByOrigin

type ByOrigin map[string]Stats

type Stats struct {
	RespBodyBytes   int `json:"resp_body_bytes"`   // Number of body bytes from origin.
	RespHeaderBytes int `json:"resp_header_bytes"` // Number of header bytes from origin.
	Responses       int `json:"responses"`         // Number of responses from origin.
	Status1xx       int `json:"status_1xx"`        // Number of 1xx "Informational" category status codes delivered from origin.
	Status200       int `json:"status_200"`        // Number of responses received with status code 200 (Success) from origin.
	Status204       int `json:"status_204"`        // Number of responses received with status code 204 (No Content) from origin.
	Status2xx       int `json:"status_2xx"`        // Number of 2xx "Success" status codes delivered from origin.
	Status301       int `json:"status_301"`        // Number of responses received with status code 301 (Moved Permanently) from origin.
	Status302       int `json:"status_302"`        // Number of responses received with status code 302 (Found) from origin.
	Status304       int `json:"status_304"`        // Number of responses received with status code 304 (Not Modified) from origin.
	Status3xx       int `json:"status_3xx"`        // Number of 3xx "Redirection" codes delivered from origin.
	Status400       int `json:"status_400"`        // Number of responses received with status code 400 (Bad Request) from origin.
	Status401       int `json:"status_401"`        // Number of responses received with status code 401 (Unauthorized) from origin.
	Status403       int `json:"status_403"`        // Number of responses received with status code 403 (Forbidden) from origin.
	Status404       int `json:"status_404"`        // Number of responses received with status code 404 (Not Found) from origin.
	Status416       int `json:"status_416"`        // Number of responses received with status code 416 (Range Not Satisfiable) from origin.
	Status4xx       int `json:"status_4xx"`        // Number of 4xx "Client Error" codes delivered from origin.
	Status500       int `json:"status_500"`        // Number of responses received with status code 500 (Internal Server Error) from origin.
	Status501       int `json:"status_501"`        // Number of responses received with status code 501 (Not Implemented) from origin.
	Status502       int `json:"status_502"`        // Number of responses received with status code 502 (Bad Gateway) from origin.
	Status503       int `json:"status_503"`        // Number of responses received with status code 503 (Service Unavailable) from origin.
	Status504       int `json:"status_504"`        // Number of responses received with status code 504 (Gateway Timeout) from origin.
	Status505       int `json:"status_505"`        // Number of responses received with status code 505 (HTTP Version Not Supported) from origin.
	Status5xx       int `json:"status_5xx"`        // Number of 5xx "Server Error" codes delivered from origin.
}
