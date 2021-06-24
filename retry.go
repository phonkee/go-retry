package retry

import (
	"context"
	"errors"
	"fmt"
)

var (
	// DefaultBackoff is set when no backoff is set
	DefaultBackoff = ConstantBackoff(0)

	// ErrMaxRetriesExceeded is returned when max retries is set and max retries exceed
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")

	// ErrPanic when callback panics
	ErrPanic = errors.New("panic")

	// ErrRetryAgain should be returned when we want to retry call again
	ErrRetryAgain = errors.New("please retry again")

	// ErrInvalidCallable is returned when callable was not provided
	ErrInvalidCallable = errors.New("callable not provided")
)

// contextKey for custom context key
type contextKey int

const (
	// key that we set values to context
	key contextKey = 1
)

// Retry interface
type Retry interface {

	// Do performs retried call
	// calls given function until it succeeds or until it exceeds max retries, or context timeouts
	Do(ctx context.Context) error

	// Retry stores callable to be called in Do
	Retry(callable func(ctx context.Context) error) Retry
}

// New instantiates new Retry interface
func New(options ...Option) Retry {
	result := &retry{
	}

	// apply all options
	for _, opt := range options {
		opt(result)
	}

	// if there is no backoff, set default value
	if result.backoff == nil {
		result.backoff = DefaultBackoff
	}

	return result
}

// retry implements Retry interface
type retry struct {
	backoff    Backoff
	maxRetries uint
	callable   func(ctx context.Context) error
}

// Do calls callable multiple times, when callable returns error, we continue with next iteration
// Do ends when either:
//   * callable succeeds
//   * callable returns error (other than ErrRetryAgain)
//   * context timeouts
//   * max retries exceed
func (r *retry) Do(ctx context.Context) error {
	var (
		err error
	)

	if r.callable == nil {
		return ErrInvalidCallable
	}

	// now run main loop
	for current := uint(0); ; current++ {

		// check for max retries (if available)
		if r.maxRetries > 0 && current > r.maxRetries {
			return ErrMaxRetriesExceeded
		}

		// check if context was closed
		if err = ctx.Err(); err != nil {
			return err
		}

		// create context with values
		ctx := context.WithValue(ctx, key, Info{
			CurrentRetry: current,
			MaxRetries:   r.maxRetries,
		})

		// wait for backoff only when retry (not first call - technically not retry)
		if current > 0 {
			if err = r.backoff.Wait(ctx); err != nil {
				return err
			}
		}

		// call callback with panic protection
		func() {
			// recover from panic
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%w: %v", ErrPanic, r)
				}
			}()

			// now call callable with custom context
			err = r.callable(ctx)
		}()

		// now check for errors
		// if error is ErrRetryAgain we retry call again
		// if error is other error, or no error, we return back to caller
		if err != nil {

			// we need to retry again
			if err == ErrRetryAgain {

				// clear error first
				err = nil

				// continue to next call
				continue
			}

			// error from callback that needs to be returned to caller
			return err
		}

		// callback returned no error, we gladly return it back to caller
		return nil
	}

}

// Retry stores callable
func (r *retry) Retry(callable func(ctx context.Context) error) Retry {
	r.callable = callable
	return r
}
