package execution

import (
	"context"
	"time"
)

// RunWithTimeout function with timeout handling
func RunWithTimeout(parentCtx context.Context, timeout time.Duration, exec func() error) error {
	ctx, ctxCancel := context.WithTimeout(parentCtx, timeout)
	defer ctxCancel()

	execDone := make(chan error, 1)
	panicChan := make(chan interface{}, 1)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		execDone <- exec()
	}()

	select {
	case execErr := <-execDone:
		return execErr
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done():
		return ctx.Err()
	}
}
