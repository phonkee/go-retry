package retry

import (
	"context"
	"time"
)

// Backoff interface
type Backoff interface {

	// Wait waits until next iteration can be done, or when context is cancelled
	Wait(ctx context.Context) error
}

// BackoffFunc is function that implements Backoff
type BackoffFunc func(ctx context.Context) error

// satisfy Backoff interface
func (s BackoffFunc) Wait(ctx context.Context) error {
	return s(ctx)
}

// ConstantBackoff constantly waits with the same duration
func ConstantBackoff(duration time.Duration) Backoff {
	return BackoffFunc(func(ctx context.Context) error {
		tmr := time.NewTimer(duration)
		for {
			select {
			case <-tmr.C:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

// ExponentialBackoff exponentially adds duration
func ExponentialBackoff(start time.Duration) Backoff {
	return BackoffFunc(func(ctx context.Context) error {
		info := GetInfo(ctx)
		tmr := time.NewTimer(start + time.Duration(1<<(info.CurrentRetry-1)))
		for {
			select {
			case <-tmr.C:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}
