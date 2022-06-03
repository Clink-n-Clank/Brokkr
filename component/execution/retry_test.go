package execution

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunWithRetry(t *testing.T) {
	execTimes := 0
	unitFunc := func() error {
		execTimes++
		if execTimes == 2 {
			return nil
		}

		return fmt.Errorf("some error")
	}

	assert.NoError(
		t,
		RunWithRetry(2, time.Nanosecond, unitFunc),
		"expected no error from function on 2 attempt",
	)
}
