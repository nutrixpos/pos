package services

import (
	"context"
	"testing"
	"time"

	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/stretchr/testify/require"
)

// newServiceEnv creates an isolated database for a service test.
func newServiceEnv(t *testing.T) *testutil.TestEnv {
	t.Helper()
	return testutil.NewTestEnv(t, testutil.BackendFromEnv())
}

func nowTime() time.Time { return time.Now() }

func seedServiceDoc(t *testing.T, env *testutil.TestEnv, collection string, doc interface{}) {
	t.Helper()
	_, err := env.Client.Database(env.Config.Databases[0].Database).Collection(collection).
		InsertOne(context.Background(), doc)
	require.NoError(t, err)
}
