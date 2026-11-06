package veille

import (
	"context"
	"sync"
)

// BatchResult is the outcome of one Batch item.
type BatchResult[T any] struct {
	Index int
	Value T
	Err   error
}

// Batch runs task for every item with at most concurrency calls in flight (DefaultConcurrency when
// concurrency <= 0) and returns one result per item, in input order. A failing item never stops the others;
// items not started when ctx is done get ctx.Err().
func Batch[I, T any](ctx context.Context, items []I, concurrency int, task func(ctx context.Context, item I) (T, error)) []BatchResult[T] {
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}
	results := make([]BatchResult[T], len(items))
	slots := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, item := range items {
		results[i].Index = i
		select {
		case <-ctx.Done():
			results[i].Err = ctx.Err()
			continue
		case slots <- struct{}{}:
		}
		wg.Add(1)
		go func(i int, item I) {
			defer wg.Done()
			defer func() { <-slots }()
			results[i].Value, results[i].Err = task(ctx, item)
		}(i, item)
	}
	wg.Wait()
	return results
}
