package circuitbreaker

import (
	"fmt"
)

// State is a type that represents a state of CircuitBreaker.
type State byte

// These constants are states of CircuitBreaker.
const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

// String implements stringer interface.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateHalfOpen:
		return "Half-Open"
	case StateOpen:
		return "Open"
	default:
		return fmt.Sprintf("unknown state: %d", s)
	}
}
