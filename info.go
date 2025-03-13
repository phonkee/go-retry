package retry

import "context"

// Info about current retry
type Info struct {
	CurrentRetry uint
	MaxRetries   uint
}

// GetInfo returns info about current retry, you can call this in callback function, or even in your custom backoff
// implementations
func GetInfo(ctx context.Context) (info Info) {
	info, _ = ctx.Value(key).(Info)
	return info
}