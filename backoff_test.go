package retry

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {

	wg := sync.WaitGroup{}

	type dataItem struct {
		start          time.Duration
		retry          uint
		expectDuration int
	}

	data := []dataItem{
		{0, 0, 0},
		{0, 1, 1},
		{0, 2, 2},
		{0, 3, 4},
		{0, 4, 8},
		{0, 5, 16},
		{0, 6, 32},
	}

	wg.Add(len(data))

	for _, item := range data {
		go func(item dataItem) {
			defer wg.Done()
			b := ExponentialBackoff(item.start)

			step1 := time.Now()

			// create context with values
			ctx := context.WithValue(context.Background(), key, Info{
				CurrentRetry: item.retry,
				MaxRetries:   1000,
			})
			_ = b.Wait(ctx)
			step2 := time.Now()
			sub := step2.Sub(step1)

			if item.expectDuration != int(sub.Seconds()) {
				t.Fatalf("expected: %v, got: %v", item.expectDuration, int(sub.Seconds()))
			}
		}(item)
	}

	// wait for all to finish
	wg.Wait()

}
