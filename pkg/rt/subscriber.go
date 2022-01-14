package rt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/fastly-exporter/pkg/gen"
)

// Subscriber continuously polls rt.fastly.com for a single service ID. Response
// data is processed and used to update a set of Prometheus metrics.
type Subscriber struct {
	SubscriberConfig
}

// SubscriberConfig collects the parameters required to construct a subscriber.
type SubscriberConfig struct {
	// Client used to query rt.fastly.com. If not provided, http.DefaultClient
	// is used, which may not include the desired User-Agent, among other
	// things.
	Client HTTPClient

	// Token provided as the Fastly-Key when querying rt.fastly.com.
	Token string

	// ServiceID managed by the subscriber. Required.
	ServiceID string

	// Metrics that will be updated by the subscriber. Required.
	Metrics *gen.Metrics

	// Metadata yields per-service metadata like service name, which ends up in
	// Prometheus labels. If not provided, relevant labels will have "unknown"
	// or zero values.
	Metadata MetadataProvider

	// Delay is called within Run when there needs to be a delay between calls
	// to the real-time stats API. If not provided, a suitable default based on
	// time.After is used. Only meant for tests.
	Delay func(context.Context, time.Duration)

	// Postprocess is called within Run immediately after gen.Metrics.Process is
	// called. If not provided, a no-op default is used. Only meant for tests.
	Postprocess func()

	// Logger is used for runtime diagnostic information.
	// If not provided, a no-op logger is used.
	Logger log.Logger
}

func (c *SubscriberConfig) validate() error {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	if c.ServiceID == "" {
		return fmt.Errorf("subscriber: service ID is required")
	}

	if c.Metrics == nil {
		return fmt.Errorf("subscriber: metrics are required")
	}

	if c.Metadata == nil {
		c.Metadata = nopMetadataProvider{}
	}

	if c.Delay == nil {
		c.Delay = func(ctx context.Context, d time.Duration) {
			select {
			case <-ctx.Done():
			case <-time.After(d):
			}
		}
	}

	if c.Postprocess == nil {
		c.Postprocess = func() {}
	}

	if c.Logger == nil {
		c.Logger = nopLogger
	}

	return nil
}

// NewSubscriber returns a ready-to-use subscriber.
// Run must be called to update the metrics.
func NewSubscriber(c SubscriberConfig) (*Subscriber, error) {
	err := c.validate()
	return &Subscriber{SubscriberConfig: c}, err
}

// Run polls rt.fastly.com in a hot loop, collecting real-time stats information
// and emitting it to the Prometheus metrics provided to the constructor. The
// method returns when the context is canceled, or a non-recoverable error
// occurs.
func (s *Subscriber) Run(ctx context.Context) error {
	var ts uint64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			name, result, delay, newts, fatal := s.query(ctx, ts)
			s.Metrics.RealtimeAPIRequestsTotal.WithLabelValues(s.ServiceID, name, string(result)).Inc()
			if fatal != nil {
				return fatal
			}
			s.Metrics.LastSuccessfulResponse.WithLabelValues(s.ServiceID, name).Set(float64(time.Now().Unix()))
			if delay > 0 {
				s.Delay(ctx, delay)
			}
			ts = newts
		}
	}
}

// query rt.fastly.com to get a batch of real-time stats for the service
// represented by the subscriber, and at the given timestamp. The first call to
// query for a subscriber should pass a timestamp value of zero. Subsequent
// calls should pass the newts value received from the previous call.
//
// The method may block for several seconds; cancel the context to provoke early
// termination. On success, the received real-time data is processed, and the
// Prometheus metrics related to the Fastly service are updated.
//
// Returns the current name of the service, the broad class of result of the API
// request, any delay that should pass before query is invoked again, the new
// timestamp that should be provided to the next call to query, and an error.
// Recoverable errors are logged internally and not returned, so any non-nil
// error returned by this method should be considered fatal to the subscriber.
func (s *Subscriber) query(ctx context.Context, ts uint64) (currentName string, result apiResult, delay time.Duration, newts uint64, fatal error) {
	name, ver, found := s.Metadata.Metadata(s.ServiceID)
	version := strconv.Itoa(ver)
	if !found {
		name, version = s.ServiceID, "unknown"
	}
	s.Metrics.ServiceInfo.WithLabelValues(s.ServiceID, name, version).Set(1)

	// rt.fastly.com blocks until it has data to return.
	// It's safe to call in a (single-threaded!) hot loop.
	u := fmt.Sprintf("https://rt.fastly.com/v1/channel/%s/ts/%d", url.QueryEscape(s.ServiceID), ts)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return name, apiResultError, 0, ts, fmt.Errorf("error constructing real-time stats API request: %w", err)
	}

	req.Header.Set("Fastly-Key", s.Token)
	req.Header.Set("Accept", "application/json")
	resp, err := s.Client.Do(req.WithContext(ctx))
	if err != nil {
		levelForError(s.Logger, err).Log("during", "execute request", "err", err)
		return name, apiResultError, time.Second, ts, nil
	}

	var response gen.APIResponse
	if err := jsoniterAPI.NewDecoder(resp.Body).Decode(&response); err != nil {
		resp.Body.Close()
		level.Error(s.Logger).Log("during", "decode response", "err", err)
		return name, apiResultError, time.Second, ts, nil
	}
	resp.Body.Close()

	apiErr := response.Error
	if apiErr == "" {
		apiErr = "<none>"
	}

	switch resp.StatusCode {
	case http.StatusOK:
		level.Debug(s.Logger).Log("status_code", resp.StatusCode, "response_ts", response.Timestamp, "err", apiErr)
		if strings.Contains(apiErr, "No data available") {
			result = apiResultNoData
		} else {
			result = apiResultSuccess
		}
		gen.Process(&response, s.ServiceID, name, version, s.Metrics)
		s.Postprocess()

	case http.StatusUnauthorized, http.StatusForbidden:
		result = apiResultError
		level.Error(s.Logger).Log("status_code", resp.StatusCode, "response_ts", response.Timestamp, "err", apiErr, "msg", "token may be invalid")
		delay = 15 * time.Second

	default:
		result = apiResultUnknown
		level.Error(s.Logger).Log("status_code", resp.StatusCode, "response_ts", response.Timestamp, "err", apiErr)
		delay = 5 * time.Second
	}

	return name, result, delay, response.Timestamp, nil
}

//
//
//

var jsoniterAPI = jsoniter.ConfigFastest

type apiResult string

const (
	apiResultUnknown apiResult = "unknown"
	apiResultError   apiResult = "error"
	apiResultNoData  apiResult = "no data"
	apiResultSuccess apiResult = "success"
)

type nopMetadataProvider struct{}

func (nopMetadataProvider) Metadata(string) (string, int, bool) { return "", 0, false }

//
//
//

var nopLogger = log.NewNopLogger()

func levelForError(base log.Logger, err error) log.Logger {
	switch {
	case errors.Is(err, context.Canceled):
		return level.Debug(base)
	case err != nil:
		return level.Error(base)
	default:
		return nopLogger
	}
}
