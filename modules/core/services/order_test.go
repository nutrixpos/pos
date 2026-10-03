package services

import (
	"sync"
	"testing"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetOrderDisplayId_ConcurrentUnique(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	const workers = 25
	type result struct {
		id  string
		err error
	}
	results := make(chan result, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc := &OrderService{Logger: env.Logger, Config: env.Config}
			id, err := svc.GetOrderDisplayId()
			results <- result{id: id, err: err}
		}()
	}
	wg.Wait()
	close(results)

	seen := map[string]bool{}
	for r := range results {
		require.NoError(t, r.err)
		require.NotEmpty(t, r.id)
		require.False(t, seen[r.id], "duplicate display id %s", r.id)
		seen[r.id] = true
	}
	require.Len(t, seen, workers)
}
