package retry

import (
	"context"
	"time"
)

// Option for retry
type Option func(*retry)

// WithBackoff adds custom backoff
func WithBackoff(backoff Backoff) Option {
	return func(r *retry) {
		r.backoff = backoff
	}
}

// WithConstantBackoff adds constant backoff with given sleep between calls
func WithConstantBackoff(sleep time.Duration) Option {
	return WithBackoff(ConstantBackoff(sleep))
}

// WithExponentialBackoff add exponential backoff with given start duration
func WithExponentialBackoff(start time.Duration) Option {
	return WithBackoff(ExponentialBackoff(start))
}

// WithMaxRetries sets max retries, if max retries is set to 0, there is no limit
func WithMaxRetries(maxRetries uint) Option {
	return func(r *retry) {
		r.maxRetries = maxRetries
	}
}

// WithContext adds custom context
func WithContext(ctx context.Context) Option {
	return func(r *retry) {
		r.context = ctx
	}
}
