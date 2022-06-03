package execution

import (
	"fmt"
	"time"
)

// RunWithRetry of function with  time
func RunWithRetry(attempts uint, backoff time.Duration, exec func() error) (execErr error) {
	for i := uint(0); i < attempts; i++ {
		if execErr = exec(); execErr == nil {
			return nil
		}

		time.Sleep(backoff)
		backoff <<= 2
	}

	return fmt.Errorf("failed to execute function in (%d) attempts, last error: %s", attempts, execErr)
}
