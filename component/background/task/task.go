package task

import (
	"context"
	"sync"
	"time"

	"github.com/Clink-n-Clank/Brokkr/component/background"
)

type (
	// BackgroundTask a process that works in configured iteration to execute job handling in the background
	BackgroundTask struct {
		ticker   *time.Ticker
		state    processState
		name     string
		severity background.ProcessSeverity

		handler           func() error
		execInterval      time.Duration
		processingTimeout time.Duration
	}

	// processState of the worker
	processState struct {
		isRunningTask     bool
		pendingToShutdown bool

		gracefulShutdown         chan struct{}
		gracefulShutdownCallback func()

		sync.Mutex
	}

	// Option for task execution, timeouts, handler, tick interval etc.
	Option func(c *BackgroundTask)
)

// SetExecInterval how often task must be executed
func SetExecInterval(interval time.Duration) Option {
	return func(c *BackgroundTask) {
		c.execInterval = interval
	}
}

// SetProcessingTimeout for task handling
func SetProcessingTimeout(interval time.Duration) Option {
	return func(c *BackgroundTask) {
		c.processingTimeout = interval
	}
}

// SetHandler of the task (logic that will be performed)
func SetHandler(handler func() error) Option {
	return func(c *BackgroundTask) {
		c.handler = handler
	}
}

// SetSeverity of how important for the application to run this task
func SetSeverity(s background.ProcessSeverity) Option {
	return func(c *BackgroundTask) {
		c.severity = s
	}
}

// NewBackgroundTask a new instance
func NewBackgroundTask(TaskName string, GracefulShutdownCallback func(), opts ...Option) *BackgroundTask {
	cw := &BackgroundTask{
		name:     TaskName,
		severity: background.TaskSeverityMajor,
	}

	for _, o := range opts {
		o(cw)
	}

	cw.ticker = time.NewTicker(cw.execInterval)
	cw.state = processState{
		gracefulShutdown:         make(chan struct{}),
		gracefulShutdownCallback: GracefulShutdownCallback,
	}

	return cw
}

// GetName of the task
func (t *BackgroundTask) GetName() string {
	return t.name
}

// GetSeverity of the task
func (t *BackgroundTask) GetSeverity() background.ProcessSeverity {
	return t.severity
}

// OnStart event to be called when main loop will be started
func (t *BackgroundTask) OnStart(ctx context.Context) error {
	if err := t.processJob(); err != nil {
		return err
	}

	for {
		select {
		case <-t.ticker.C:
			_ = t.processJob()
		case <-t.state.gracefulShutdown:
			close(t.state.gracefulShutdown)
			t.state.gracefulShutdownCallback()
			ctx.Done()

			return nil
		}
	}
}

// OnStop event to be called when main loop will be started
func (t *BackgroundTask) OnStop(ctx context.Context) error {
	defer ctx.Done()

	t.togglePendingToShutdown()

	// If there is no running jobs, stop it now
	if !t.IsProcessingJob() {
		t.state.gracefulShutdown <- struct{}{}
	}

	return nil
}

func (t *BackgroundTask) processJob() error {
	if t.IsPendingToShutdown() || t.IsProcessingJob() {
		return nil
	}

	t.toggleIsProcessingJob()
	defer t.toggleIsProcessingJob()

	return t.handler()
}

// IsPendingToShutdown a worker
func (t *BackgroundTask) IsPendingToShutdown() bool {
	t.state.Lock()
	defer t.state.Unlock()

	return t.state.pendingToShutdown
}

func (t *BackgroundTask) togglePendingToShutdown() {
	t.state.Lock()
	defer t.state.Unlock()

	t.state.pendingToShutdown = !t.state.pendingToShutdown
}

// IsProcessingJob in worker cycle handling
func (t *BackgroundTask) IsProcessingJob() bool {
	t.state.Lock()
	defer t.state.Unlock()

	return t.state.isRunningTask
}

func (t *BackgroundTask) toggleIsProcessingJob() {
	t.state.Lock()
	defer t.state.Unlock()

	t.state.isRunningTask = !t.state.isRunningTask
}
