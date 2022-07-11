package tester

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

func (a *ActorAPI) WhenCurrentTimeAfter(clock string) error {
	tGiven, tNow, err := parseTime(clock)
	if err != nil {
		return err
	}

	if tNow.After(tGiven) {
		return nil
	}

	return godog.ErrPending
}

func (a *ActorAPI) WhenCurrentTimeBefore(clock string) error {
	tGiven, tNow, err := parseTime(clock)
	if err != nil {
		return err
	}

	if tNow.Before(tGiven) {
		return nil
	}

	return godog.ErrPending
}

func parseTime(clock string) (time.Time, time.Time, error) {
	t, tErr := time.Parse(time.Kitchen, clock)
	if tErr != nil {
		return time.Time{}, time.Time{}, tErr
	}

	tNow, _ := time.Parse(time.Kitchen, time.Now().Format(time.Kitchen))
	return t, tNow, nil
}

func (a *ActorAPI) IWaitSeconds(amount string) error {
	t, tErr := strconv.Atoi(amount)
	if tErr != nil {
		return tErr
	}

	time.Sleep(time.Duration(t) * time.Second)

	return nil
}

func (a *ActorAPI) SetConcurrentRequestsAmount(reqAmount, reqDelayMilliseconds string) error {
	ra, raErr := strconv.Atoi(reqAmount)
	if raErr != nil {
		return raErr
	}

	a.StressConcurrentRequests = uint32(ra)


	rd, rdErr := strconv.Atoi(reqDelayMilliseconds)
	if rdErr != nil {
		return rdErr
	}

	a.StressConcurrentRequestsDelay = time.Duration(rd) * time.Millisecond

	return nil
}

func (a *ActorAPI) PrintMetricsResult() error {
	fmt.Println(a.Metrics.ToString())

	return nil
}

func (a *ActorAPI) MetricsRequestDropIsLowerThan(p string) error {
	r, rErr := a.Metrics.RequestDropPercentage()
	if rErr != nil {
		return rErr
	}

	pf, err := strconv.ParseFloat(p, 32)
	if  err != nil {
		return err
	}

	if r > pf {
		return fmt.Errorf("requests has more drops than expected (%f) but actuall drops are (%f)", pf, r)
	}

	return nil
}

func (a *ActorAPI) PrintLastResponse() error {
	fmt.Println("------------- DEBUG OUTPUT -------------")
	fmt.Printf("Request URL: %s\n", a.HTTPLastResp.ResponseObj.Request.URL.String())
	fmt.Printf("Request Payload: %s\n", a.Payload)
	fmt.Printf("Response:\n%s\n", createHTTPResponseDumpFrom(a.HTTPLastResp))
	fmt.Println("------------- DEBUG OUTPUT -------------")

	return nil
}
