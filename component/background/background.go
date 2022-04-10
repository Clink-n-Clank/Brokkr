package background

import "context"

// ProcessSeverity identify how is imported background task is to execute
type ProcessSeverity byte

const (
	TaskSeverityMajor ProcessSeverity = iota
	TaskSeverityMinor
)

// Process that must be executed inside microservice, could be a server, events, parsers, aggregators etc...
type Process interface {
	// GetName of the task
	GetName() string
	// GetSeverity of the task
	GetSeverity() ProcessSeverity
	// OnStart event to be called when main loop will be started
	OnStart(ctx context.Context) error
	// OnStop event to be called when main loop will be started
	OnStop(ctx context.Context) error
}

// IsCriticalToStop verifying if task critical to execute
func IsCriticalToStop(t Process) bool {
	return t.GetSeverity() == TaskSeverityMajor
}
