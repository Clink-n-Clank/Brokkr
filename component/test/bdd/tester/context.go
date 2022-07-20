package tester

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/Clink-n-Clank/Brokkr/component/test/bdd/tester/metrics"
)

// ActorAPI container context
type ActorAPI struct {
	HTTPClient   *http.Client
	HTTPLastResp HTTPResponse

	HTTPBaseHost url.URL
	HTTPHeaders  map[string]string
	HTTPQuery    map[string]string
	HTTPPath     []string
	Payload      string

	Storage map[string]string

	StressConcurrentRequests      uint32
	StressConcurrentRequestsDelay time.Duration

	Metrics *metrics.Collector
}

// HTTPResponse between actions
type HTTPResponse struct {
	Body        []byte
	ResponseObj http.Response
}

// NewActorAPI context
func NewActorAPI() *ActorAPI {
	c := new(http.Client)
	c.Timeout = time.Minute

	return &ActorAPI{
		HTTPClient:                    c,
		HTTPHeaders:                   make(map[string]string),
		HTTPPath:                      make([]string, 0),
		HTTPQuery:                     make(map[string]string),
		Storage:                       make(map[string]string),
		StressConcurrentRequests:      1,
		StressConcurrentRequestsDelay: time.Nanosecond,
		Metrics:                       metrics.NewMetrics(),
	}
}

// SaveToMemory key and value in ActorAPI
func (a *ActorAPI) SaveToMemory(key, value string) {
	a.Storage[key] = value
}

// HandleHTTPRequest ...
func (a *ActorAPI) HandleHTTPRequest(method, endpointPath string, payload io.Reader) (failedReason error) {
	var wg sync.WaitGroup

	for i := uint32(0); i < a.StressConcurrentRequests; i++ {
		wg.Add(1)

		time.Sleep(a.StressConcurrentRequestsDelay)

		go func() {
			var runtimeErr error
			var req *http.Request
			roundTrip := metrics.RoundTrip{}

			reqURL := a.HTTPBaseHost
			reqURL.Path = path.Join(reqURL.Path, endpointPath)

			req, runtimeErr = http.NewRequest(method, reqURL.String(), payload)
			if runtimeErr != nil {
				failedReason = runtimeErr

				return
			}

			q := req.URL.Query()
			for name, value := range a.HTTPQuery {
				q.Add(name, value)
			}
			req.URL.RawQuery = q.Encode()

			for _, p := range a.HTTPPath {
				req.URL.Path = path.Join(req.URL.Path, p)
			}

			for name, value := range a.HTTPHeaders {
				req.Header.Set(name, value)
			}

			defer func() {
				roundTrip.End = time.Now()
				a.Metrics.Collect(*req, roundTrip, runtimeErr == nil)
				wg.Done()
			}()

			roundTrip.Start = time.Now()
			resp, respErr := a.HTTPClient.Do(req)
			if respErr != nil {
				runtimeErr = fmt.Errorf(
					"failed to send request to %s, error: %w",
					req.URL.String(),
					respErr,
				)

				failedReason = runtimeErr

				return
			}

			body, bodyErr := ioutil.ReadAll(resp.Body)
			if bodyErr != nil {
				runtimeErr = fmt.Errorf(
					"failed to read body of request (%s), error: %w",
					req.URL.String(),
					bodyErr,
				)

				failedReason = runtimeErr

				return
			}

			a.HTTPLastResp = HTTPResponse{
				Body:        body,
				ResponseObj: *resp,
			}
		}()
	}

	wg.Wait()
	a.HTTPPath = make([]string, 0)
	a.HTTPQuery = make(map[string]string)

	// If that was concurrent requests ignore errors (stress tests)
	if a.StressConcurrentRequests > 1 {
		return nil
	}

	return failedReason
}
