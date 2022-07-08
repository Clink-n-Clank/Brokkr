package metrics

import (
	"fmt"
	"github.com/cucumber/godog/colors"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Collector contains collected statistic while API server works
type Collector struct {
	sync.Mutex
	// startTime when metrics was created
	startTime time.Time
	// totalCreatedRequests contains number of requests per second
	totalCreatedRequests uint32
	// totalSentRequests contains url and count
	totalSentRequests map[string]*RequestState
	// timer when need to log metrics
	timer time.Duration
}

// RequestState stores result of it
type RequestState struct {
	URL       string
	Method    string
	Succeeded uint32
	Failed    uint32
	Trip      RoundTrip
}

// RoundTrip of request to response
type RoundTrip struct {
	Start time.Time
	End   time.Time
}

// NewMetrics create new metrics collector instance
func NewMetrics() *Collector {
	return &Collector{
		timer:             time.Second,
		startTime:         time.Now(),
		totalSentRequests: map[string]*RequestState{},
	}
}

// Collect requests result
func (c *Collector) Collect(req http.Request, rt RoundTrip, isOk bool) {
	c.Lock()
	defer c.Unlock()

	c.totalCreatedRequests++
	key := strings.ToLower(fmt.Sprintf("%s_%s", req.URL.String(), req.Method))

	if _, isExist := c.totalSentRequests[key]; !isExist {
		c.totalSentRequests[key] = &RequestState{
			URL:    req.URL.String(),
			Method: req.Method,
		}
	}

	c.totalSentRequests[key].Trip = rt
	if !isOk {
		c.totalSentRequests[key].Failed++
	} else {
		c.totalSentRequests[key].Succeeded++
	}
}

// RequestDropPercentage calc % of drops
func (c *Collector) RequestDropPercentage() (float64, error) {
	var ok, bad uint32
	for _, total := range c.totalSentRequests {
		ok += total.Succeeded
		bad += total.Failed
	}

	return float64(bad) / float64(c.totalCreatedRequests) * 100, nil
}

// ToString makes pretty print format of collected metrics
func (c *Collector) ToString() (out string) {
	c.Lock()
	defer c.Unlock()

	// Header
	out = colors.Green("------------------------------------- Metrics -------------------------------------")

	// Time
	out = fmt.Sprintf("%s\n%s", out, colors.Yellow(fmt.Sprintf("│ Run time: %v │", time.Since(c.startTime))))
	// Requests count
	out = fmt.Sprintf("%s%s", out, colors.Yellow(fmt.Sprintf("│ Requests count: %d │", c.totalCreatedRequests)))
	// Unique endpoints count
	out = fmt.Sprintf("%s%s", out, colors.Yellow(fmt.Sprintf("│ Unique endpoints count: %d │", len(c.totalSentRequests))))

	var stats string
	for _, unique := range c.totalSentRequests {
		l := fmt.Sprintf(
			"├─ [%s] Endpoint: %s\n └── Total Requests hits: %d\n └── Result \n    └── Last round duration: %v \n    └── OK: %d\n    └── BAD: %d\n",
			unique.Method,
			unique.URL,
			unique.Succeeded+unique.Failed,
			unique.Trip.End.Sub(unique.Trip.Start).String(),
			unique.Succeeded,
			unique.Failed,
		)

		stats = fmt.Sprintf("%s%s", stats, colors.Yellow(l))
	}

	out = fmt.Sprintf("%s\n%s", out, stats)

	// Footer
	out = fmt.Sprintf("%s%s\n", out, colors.Green("------------------------------------- Metrics -------------------------------------"))

	return
}
