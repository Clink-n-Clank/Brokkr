package execution

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunWithTimeout(t *testing.T) {
	execTime := time.Second
	unitFunc := func() error {
		time.Sleep(execTime)

		return nil
	}

	execErr := RunWithTimeout(context.Background(), time.Nanosecond, unitFunc)
	assert.Error(t, execErr, "expected error, exec function timeout")
}
