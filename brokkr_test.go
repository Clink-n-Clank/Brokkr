package brokkr

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Clink-n-Clank/Brokkr/component/background"
	"github.com/Clink-n-Clank/Brokkr/component/execution"
)

func TestBrokkr_StartStop(t *testing.T) {
	testTimeout := 5 * time.Second
	testFunction := func() error {
		c := NewBrokkr([]Options{SetForceStopTimeout(time.Second)}...)

		go func() {
			time.Sleep(500 * time.Millisecond)
			assert.NoError(t, c.Stop())
		}()

		assert.NoError(t, c.Start())

		return nil
	}

	testExecErr := execution.RunWithTimeout(context.Background(), testTimeout, testFunction)
	assert.NoError(
		t,
		testExecErr,
		"Test didn't finish in time, possible dead lock with Core process lifetime management",
	)
}

func TestBrokkr_BackgroundTaskExecution(t *testing.T) {
	testCases := []struct {
		caseName               string
		backgroundTaskSeverity background.ProcessSeverity
	}{
		{
			caseName:               "Major task severity - must stop main loop",
			backgroundTaskSeverity: background.TaskSeverityMajor,
		},
		{
			caseName:               "Minor task severity - main loop will remain work",
			backgroundTaskSeverity: background.TaskSeverityMinor,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.caseName, func(t *testing.T) {
			tBgTask := testBackgroundTask{sv: tCase.backgroundTaskSeverity}

			c := NewBrokkr(
				SetForceStopTimeout(time.Second),
				AddBackgroundTasks(tBgTask),
			)

			if background.IsCriticalToStop(tBgTask) {
				assert.NotNil(t, c.Start())
			} else {
				okToStart := func() error {
					startErr := c.Start()
					assert.NoError(t, startErr, "Unexpected error during the start, minor background task can fail")
					return startErr
				}

				_ = execution.RunWithTimeout(context.Background(), time.Millisecond, okToStart)
			}
		})
	}
}

type testBackgroundTask struct {
	sv background.ProcessSeverity
}

func (t testBackgroundTask) GetName() string {
	return "test"
}

func (t testBackgroundTask) GetSeverity() background.ProcessSeverity {
	return t.sv
}

func (t testBackgroundTask) OnStart(context.Context) error {
	return fmt.Errorf("unit test error")
}

func (t testBackgroundTask) OnStop(context.Context) error {
	return nil
}
