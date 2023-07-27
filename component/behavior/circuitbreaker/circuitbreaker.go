package circuitbreaker

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// ErrCircuitOpen is returned when the state of Circuit Breaker is open.
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Action function in CircuitBreaker.
type Action func() (any, error)

// Configuration for CircuitBreaker:
//
// MaxFailuresThreshold Maximum number of failures allowed.
//
// ResetTimeout in seconds, is the period of the open state. After which the state of the CircuitBreaker becomes half-open.
// fields can be used for example as ENV variables in your project like MY_APP_CB_MAX_FAILURES_THRESHOLD:"3"
type Configuration struct {
	MaxFailuresThreshold string `json:"cb_max_failures_threshold,omitempty"`
	ResetTimeout         string `json:"cb_reset_timeout,omitempty"`
}

// CircuitBreaker component that is designed  to prevent sending execution that are likely to fail.
type CircuitBreaker struct {
	OnSuccess func() // OnSuccess will be called when Action can be processed, when state is StateClosed or StateHalfOpen
	OnFailure func() // OnFailure will be triggered when CircuitBreaker moved to StateOpen and Action failed to execute and will return error.

	mu           sync.Mutex // A mu is a mutual exclusion lock
	currentState State

	timeout      time.Duration // Duration when state must be closed
	lastAttempt  time.Time     // Timestamp of the last attempt to execution
	failureCount uint64        // Current count of consecutive failures
	failureLimit uint64        // Number of failures that will switch the state from closed to open
}

// NewCircuitBreaker creates a new CircuitBreaker instance with the specified configuration.
func NewCircuitBreaker(cfg Configuration) (*CircuitBreaker, error) {
	errTmpl := "failed to parse parameter for %s"

	iRestTimout, iRestTimoutErr := strconv.Atoi(cfg.ResetTimeout)
	if iRestTimoutErr != nil {
		return nil, errors.Join(fmt.Errorf(errTmpl, "ResetTimout in CircuitBreaker"), iRestTimoutErr)
	}
	iMaxFailuresThreshold, iMaxFailuresThresholdErr := strconv.Atoi(cfg.MaxFailuresThreshold)
	if iMaxFailuresThresholdErr != nil {
		return nil, errors.Join(fmt.Errorf(errTmpl, "ResetTimout in MaxFailuresThreshold"), iMaxFailuresThresholdErr)
	}

	cb := &CircuitBreaker{
		currentState: StateClosed,
		timeout:      time.Second * time.Duration(iRestTimout),
		OnSuccess:    func() {},
		OnFailure:    func() {},
		failureLimit: uint64(iMaxFailuresThreshold),
	}

	return cb, nil
}

// Proceed function in CircuitBreaker
func (cb *CircuitBreaker) Proceed(action Action) (any, error) {
	switch cb.GetState() {
	case StateOpen:
		if cb.isTimeout() {
			cb.setState(StateHalfOpen)
			return cb.Proceed(action)
		} else {
			return nil, ErrCircuitOpen
		}
	case StateHalfOpen, StateClosed:
		result, err := action()
		if err != nil {
			cb.recordFailure()
			cb.OnFailure()
			return nil, err
		}

		cb.OnSuccess()

		return result, nil
	default:
		return nil, ErrCircuitOpen
	}
}

// GetState returns current state of the Circuit Breaker.
func (cb *CircuitBreaker) GetState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.currentState
}

func (cb *CircuitBreaker) setState(state State) {
	cb.currentState = state
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastAttempt = time.Now()
	if cb.failureCount > cb.failureLimit {
		cb.setState(StateOpen)
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0
	cb.setState(StateClosed)
}

func (cb *CircuitBreaker) isTimeout() bool {
	return time.Since(cb.lastAttempt) > cb.timeout
}
