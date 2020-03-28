# go-retry

Simple package to retry function calls.

# usage

This package provides simple interface that u can use to retry calls.
Let's show example:

```go
// Call retry function
err := retry.New(
    retry.WithConstantBackoff(time.Second * 2), 
    retry.WithMaxRetries(10), 
    retry.WithContext(context.Background()),
).RetryCall(func(ctx context.Context) error {
	
    // just skip first call and retry this call 
    if GetInfo(ctx).CurrentRetry == 0 {
        return ErrRetryAgain
    }
    
    // we can get information about current retry
    info := retry.GetInfo(ctx)
    fmt.Printf("currently we are %v from %v", info.CurrentRetry, info.MaxRetries)
    
    // return error
    return fmt.Errorf("nope")
})
```

Even in scenarios where you need to indefinitely retry calls:

```go
err := retry.New(
    retry.WithExponentialBackoff(time.Second),
).RetryCall(func(ctx context.Context) error {
    // Do some work
    return nil
})
```

You can rely on context to be canceled

```go
func() {
    ctx, cf := context.WithTimeout(context.Background(), time.Second * 60)
    defer cf

    // No need to set max retries, when context is canceled retry will end (assuming callback is returned when context is canceled)
    err := retry.New(
        retry.WithExponentialBackoff(time.Second), 
        retry.WithContext(ctx),
    ).RetryCall(func(ctx context.Context) error {
        // Do some work
        return nil
    })
}()
```

# backoff

Go-retry provides constant and exponential retry. You can even implement your own
retry if you need and pass it to `WithBackoff` option.

Example:

```go
// Custom backoff implementation
err := retry.New(
    retry.WithBackoff(BackoffFunc(func(ctx context.Context) error {
        timer := time.NewTimer(time.Second * (GetInfo(ctx).CurrentRetry % 2) + time.Second)
        for {
            select {
            case <-timer.C:
                return nil
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    })).RetryCall(func(ctx context.Context) error {
    	// Do some work
    	return nil
    })
)
```

# author

Peter Vrba <phonkee@phonkee.eu>
