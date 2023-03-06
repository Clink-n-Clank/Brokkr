package task

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Clink-n-Clank/Brokkr/component/execution"
)

func TestCronWorker_StartProcessStop(t *testing.T) {
	taskName := "UnitTestTask"
	isTaskHandled := false
	taskHandlerCtx := context.Background()

	testTimeout := 5 * time.Second
	testFunction := func() error {
		var wg sync.WaitGroup

		taskStopCallback := func() {
			taskHandlerCtx.Done()
			wg.Done()
		}
		taskHandlerFunc := func() error {
			isTaskHandled = true
			time.Sleep(500 * time.Millisecond)
			return nil
		}

		var opts []Option
		opts = append(opts, SetExecInterval(time.Microsecond))
		opts = append(opts, SetProcessingTimeout(time.Second))
		opts = append(opts, SetHandler(taskHandlerFunc))

		c := NewBackgroundTask("UnitTestCron", taskStopCallback, opts...)

		cronWorkerCtx := context.Background()

		wg.Add(1)
		go func() {
			_ = c.OnStart(cronWorkerCtx)
		}()

		_ = c.OnStop(cronWorkerCtx)
		wg.Wait()

		return nil
	}

	testExecErr := execution.RunWithTimeout(context.Background(), testTimeout, testFunction)
	assert.NoError(
		t,
		testExecErr,
		"Test didn't finish in time, possible dead lock in Cron loop",
	)

	if isTaskHandled {
		t.Logf("Expected to be handled at least one job in worker %s", taskName)
		t.FailNow()
	}

	assert.NoError(t, taskHandlerCtx.Err(), "job context must have no error")
}

func TestCronWorker_StartAndStop(t *testing.T) {
	testTimeout := time.Second
	testFunction := func() error {
		taskHandlerCtx := context.Background()
		taskStopCallback := func() { taskHandlerCtx.Done() }
		taskHandlerFunc := func() error { return nil }

		var opts []Option
		opts = append(opts, SetExecInterval(time.Nanosecond))
		opts = append(opts, SetProcessingTimeout(time.Second))
		opts = append(opts, SetHandler(taskHandlerFunc))

		c := NewBackgroundTask("UnitTestCron", taskStopCallback, opts...)

		cronWorkerCtx := context.Background()

		go func() { _ = c.OnStart(cronWorkerCtx) }()
		go func() { _ = c.OnStop(cronWorkerCtx) }()

		return nil
	}

	testExecErr := execution.RunWithTimeout(context.Background(), testTimeout, testFunction)
	assert.NoError(
		t,
		testExecErr,
		"Test didn't finish in time, possible dead lock in Cron loop",
	)
}
