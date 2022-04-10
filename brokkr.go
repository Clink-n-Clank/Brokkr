package brokkr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/Clink-n-Clank/Brokkr/component/background"
)

type (
	// Brokkr core loop system of the app
	Brokkr struct {
		// signals to listen and stop Brokkr
		signals []os.Signal
		// stopTimeout for force stop if exceeds
		stopTimeout time.Duration
		// backgroundTasks for Brokkr, it will launch them in the background and executes
		backgroundTasks []background.Process

		mainContext       context.Context
		mainContextCancel func()
	}

	// Options sets of configurations for Brokkr
	Options func(o *Brokkr)

	// coreContextKey child context key
	contextOfBrokkr interface{}
)

// SetForceStopTimeout redefines force shutdown timeout
func SetForceStopTimeout(t time.Duration) Options {
	return func(c *Brokkr) { c.stopTimeout = t }
}

// AddBackgroundTasks that will be executed in background of main loop
func AddBackgroundTasks(bt ...background.Process) Options {
	return func(c *Brokkr) {
		c.backgroundTasks = append(c.backgroundTasks, bt...)
	}
}

// NewBrokkr framework instance
func NewBrokkr(opts ...Options) (b *Brokkr) {
	b = &Brokkr{
		signals:     []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: 60 * time.Second,
	}

	b.mainContext, b.mainContextCancel = context.WithCancel(context.Background())
	for _, o := range opts {
		o(b)
	}

	return
}

// Start main loop and call callback in the end
func (c *Brokkr) Start() error {
	var bgTasksWG sync.WaitGroup                                             // background tasks group to control it synchronization
	interruptSignal := make(chan os.Signal, 1)                               // listen for interrupt
	TaskErrorGroup, TaskErrorGroupCtx := errgroup.WithContext(c.mainContext) // sub-task process context and it's error group

	// Init background tasks
	for _, t := range c.backgroundTasks {
		task := t

		// Setup termination workflow for background task
		TaskErrorGroup.Go(func() error {
			<-TaskErrorGroupCtx.Done()

			newUUID, errNewUUID := uuid.NewUUID()
			if errNewUUID != nil {
				return fmt.Errorf("unable to generate background task UUID, err: %v", errNewUUID)
			}

			taskStopCtx, taskStopCtxCancel := context.WithTimeout(
				c.createChildContext(newUUID.String(), t.GetName()),
				c.stopTimeout,
			)
			defer taskStopCtxCancel()

			return task.OnStop(taskStopCtx)
		})

		bgTasksWG.Add(1)

		// Setup execution for background task
		TaskErrorGroup.Go(func() error {
			defer bgTasksWG.Done()

			taskErr := task.OnStart(TaskErrorGroupCtx)
			if taskErr != nil && background.IsCriticalToStop(task) {
				return taskErr
			}

			return nil
		})
	}

	// Listen and Replay
	signal.Notify(interruptSignal, c.signals...)

	// Main loop
	TaskErrorGroup.Go(func() error {
		for {
			select {
			case <-TaskErrorGroupCtx.Done():
				return TaskErrorGroupCtx.Err()
			case <-interruptSignal:
				return c.Stop()
			}
		}
	})

	if err := TaskErrorGroup.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Stop in graceful mode
func (c *Brokkr) Stop() error {
	if c.mainContextCancel != nil {
		c.mainContextCancel()
	}

	return nil
}

// createChildContext from parent
func (c Brokkr) createChildContext(k contextOfBrokkr, v string) context.Context {
	return context.WithValue(c.mainContext, k, v)
}
