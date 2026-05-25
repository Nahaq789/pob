package sync

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

func P[T any](ctx context.Context, ids []int, fn func(context.Context, int) (T, error)) ([]T, error) {
	eg, egCtx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	var results []T

	for _, id := range ids {
		id := id
		eg.Go(func() error {
			res, err := fn(egCtx, id)
			if err != nil {
				return err
			}
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}
