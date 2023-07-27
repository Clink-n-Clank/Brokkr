package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {
	cfg := Configuration{
		MaxFailuresThreshold: "3",
		ResetTimeout:         "5",
	}

	_, err := NewCircuitBreaker(cfg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCircuitBreakerOpenAfterFailures(t *testing.T) {
	onSuccessChanged := false
	onFailureChanged := false

	cfg := Configuration{
		MaxFailuresThreshold: "3",
		ResetTimeout:         "1",
	}

	var exec Action
	exec = func() (any, error) {
		return nil, errors.New("error")
	}

	cb, cbErr := NewCircuitBreaker(cfg)
	assert.Nil(t, cbErr)

	cb.timeout = time.Millisecond
	cb.OnSuccess = func() {
		onSuccessChanged = true
	}
	cb.OnFailure = func() {
		onFailureChanged = true
	}

	for i := 0; i < 4; i++ {
		time.Sleep(time.Millisecond)
		_, err := cb.Proceed(exec)

		if errors.Is(err, ErrCircuitOpen) {
			t.Errorf("CircuitBreaker should not have opened at iteration %d", i)
		}
	}

	assert.False(t, onSuccessChanged)
	assert.True(t, onFailureChanged)
}

func TestCircuitBreakerResetsAfterSuccess(t *testing.T) {
	onSuccessChanged := false
	onFailureChanged := false

	cfg := Configuration{
		MaxFailuresThreshold: "3",
		ResetTimeout:         "1",
	}

	cb, cbErr := NewCircuitBreaker(cfg)
	assert.Nil(t, cbErr)
	cb.timeout = time.Millisecond
	cb.OnSuccess = func() {
		onSuccessChanged = true
	}
	cb.OnFailure = func() {
		onFailureChanged = true
	}

	for i := 0; i < 3; i++ {
		exec := func() (any, error) {
			return "", nil
		}
		_, err := cb.Proceed(exec)
		if errors.Is(ErrCircuitOpen, err) {
			t.Errorf("CircuitBreaker should have reset after a successful call")
		}
	}

	val, err := cb.Proceed(func() (any, error) {
		return "unit", nil
	})

	// Sleep for slightly more than the timeout period
	time.Sleep(500 * time.Millisecond)

	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "unit", val)
	assert.True(t, onSuccessChanged)
	assert.False(t, onFailureChanged)
}

func TestCircuitBreakerResetsAfterTimeout(t *testing.T) {
	onSuccessChanged := false
	onFailureChanged := false

	cfg := Configuration{
		MaxFailuresThreshold: "3",
		ResetTimeout:         "1",
	}

	var exec Action
	exec = func() (any, error) {
		return nil, errors.New("error")
	}

	cb, cbErr := NewCircuitBreaker(cfg)
	assert.Nil(t, cbErr)
	cb.timeout = time.Millisecond
	cb.OnSuccess = func() {
		onSuccessChanged = true
	}
	cb.OnFailure = func() {
		onFailureChanged = true
	}

	for i := 0; i < 3; i++ {
		_, _ = cb.Proceed(exec)
	}

	// Sleep for slightly more than the timeout period
	time.Sleep(500 * time.Millisecond)

	_, err := cb.Proceed(func() (interface{}, error) {
		return nil, nil
	})

	if errors.Is(err, ErrCircuitOpen) {
		t.Errorf("CircuitBreaker should have reset after the timeout period")
	}

	assert.True(t, onSuccessChanged)
	assert.True(t, onFailureChanged)
}

func TestProceedAfterTimeout(t *testing.T) {
	cfg := Configuration{
		MaxFailuresThreshold: "1",
		ResetTimeout:         "1",
	}

	cb, _ := NewCircuitBreaker(cfg)
	cb.timeout = time.Nanosecond

	action := func() (any, error) {
		return nil, errors.New("failed")
	}

	for i := 0; i < 4; i++ {
		_, _ = cb.Proceed(action)
	}

	if cb.GetState() != StateOpen {
		t.Errorf("Expected state to be 'StateOpen', got %v", cb.GetState())
	}

	_, _ = cb.Proceed(func() (any, error) { return nil, nil })

	if cb.GetState() != StateHalfOpen {
		t.Errorf("Expected state to be 'StateHalfOpen', got %v", cb.GetState())
	}
}
