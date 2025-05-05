package async

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Future[T any] interface {
	Await(ctx context.Context) (T, error)
	AwaitWithCallback(ctx context.Context, callback func(result T, err error))
}

type future[T any] struct {
	result    T
	err       error
	mu        sync.Mutex
	once      sync.Once
	awaitFunc func(ctx context.Context) (T, error)
}

func (f *future[T]) Await(ctx context.Context) (T, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.once.Do(func() {
		f.result, f.err = f.awaitFunc(ctx)
	})
	return f.result, f.err
}

func (f *future[T]) AwaitWithCallback(ctx context.Context, callback func(result T, err error)) {
	result, err := f.Await(ctx)
	callback(result, err)
}

func Promise[T any](f func() (T, error)) Future[T] {
	resultChan := make(chan struct {
		result T
		err    error
	}, 1)

	go func() {
		defer close(resultChan)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v\nStack trace: %s", r, string(debug.Stack()))
				var defaultResult T
				resultChan <- struct {
					result T
					err    error
				}{defaultResult, fmt.Errorf("panic recovered: %v", r)}
			}
		}()

		result, err := f()
		resultChan <- struct {
			result T
			err    error
		}{result, err}
	}()

	return &future[T]{
		awaitFunc: func(ctx context.Context) (T, error) {
			var defaultResult T
			select {
			case <-ctx.Done():
				return defaultResult, ctx.Err()
			case res := <-resultChan:
				return res.result, res.err
			}
		},
	}
}

func OfAll[T any](ctx context.Context, futures ...Future[T]) ([]T, []error) {
	var wg sync.WaitGroup
	results := make([]T, len(futures))
	errs := make([]error, len(futures))
	for i, f := range futures {
		wg.Add(1)
		go func(i int, f Future[T]) {
			defer wg.Done()
			results[i], errs[i] = f.Await(ctx)
		}(i, f)
	}
	wg.Wait()
	return results, errs
}
