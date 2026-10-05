// Package testutil provides helpers to spin up isolated databases for tests.
//
// It is only imported from test files, so it never ships in the production
// binary. Tests run against the embedded FerretDB backend by default, or
// against an external MongoDB when TEST_MONGO_URI is set.
package testutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Backend selects which database backend a test should run against.
type Backend string

const (
	// BackendFerret runs the embedded FerretDB (SQLite) backend.
	BackendFerret Backend = "ferret"
	// BackendMongo runs against the external MongoDB defined by TEST_MONGO_URI.
	BackendMongo Backend = "mongo"
)

// BackendFromEnv returns the backend selected by TEST_MONGO_URI: mongo when
// the variable is set, ferret otherwise.
func BackendFromEnv() Backend {
	if os.Getenv("TEST_MONGO_URI") != "" {
		return BackendMongo
	}
	return BackendFerret
}

type discardLogger struct{}

func (discardLogger) Info(string, ...interface{})    {}
func (discardLogger) Warning(string, ...interface{}) {}
func (discardLogger) Error(string, ...interface{})   {}

// NewLogger returns a logger that discards all output.
func NewLogger() logger.ILogger { return discardLogger{} }

// TestEnv bundles everything a test needs to talk to an isolated database.
type TestEnv struct {
	Config config.Config
	Client *mongo.Client
	Logger logger.ILogger
}

// NewTestEnv creates an isolated database of the given backend and seeds the
// settings collection. Mongo tests are skipped unless TEST_MONGO_URI is set.
// The database is torn down automatically when the test finishes.
//
// Note: the database client is a process-wide singleton, so a single test must
// create at most one env and DB tests must not run in parallel.
func NewTestEnv(t testing.TB, backend Backend) *TestEnv {
	t.Helper()

	cfg, ok := ConfigForBackend(t, backend)
	if !ok {
		t.Skipf("TEST_MONGO_URI not set; skipping %s backend test", backend)
	}

	env := &TestEnv{Config: cfg, Logger: NewLogger()}
	env.Client = NewClient(t, &cfg)

	SeedSettings(t, &cfg)

	return env
}

// ConfigForBackend returns an isolated config for the given backend. For mongo
// it returns ok=false when TEST_MONGO_URI is unset.
func ConfigForBackend(t testing.TB, backend Backend) (config.Config, bool) {
	t.Helper()

	if backend == BackendMongo {
		return mongoConfigFromEnv(t)
	}

	return config.Config{
		Env: "dev",
		Databases: []config.Database{{
			Type:     "ferret",
			FilePath: t.TempDir(),
			Database: "nutrix",
			Tables:   map[string]string{"sales": "sales"},
		}},
	}, true
}

// NewClient creates (or reuses) the process-wide database client for the given
// config and registers a cleanup that drops the test database and closes the
// client.
func NewClient(t testing.TB, cfg *config.Config) *mongo.Client {
	t.Helper()

	client, err := common.GetDatabaseClient(NewLogger(), cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cfg.Databases[0].Type != "ferret" {
			_ = client.Database(cfg.Databases[0].Database).Drop(context.Background())
		}
		common.CloseDatabase()
	})

	return client
}

// SeedSettings writes the bootstrapped settings document required by the core
// services. It mirrors the production seeder defaults so services depending on
// queues/stock thresholds behave as they do in the running app.
func SeedSettings(t testing.TB, cfg *config.Config) {
	t.Helper()

	client, err := common.GetDatabaseClient(NewLogger(), cfg)
	require.NoError(t, err)

	ctx := context.Background()
	db := client.Database(cfg.Databases[0].Database)

	names, err := db.ListCollectionNames(ctx, bson.M{"name": "settings"})
	require.NoError(t, err)

	if len(names) == 0 {
		require.NoError(t, db.CreateCollection(ctx, "settings"))
	}

	settings := models.Settings{
		Id:        primitive.NewObjectID().Hex(),
		Inventory: models.MaterialSettings{StockAlertTreshold: 1000},
		Orders: models.OrderSettings{
			Queues:                       []models.OrderQueueSettings{{Prefix: "A", Next: 1}},
			DefaultCostCalculationMethod: "average",
		},
		Language:           models.LanguageSettings{Code: "en", Language: "English"},
		AutoOpenCashDrawer: true,
		PaymentSources:     []models.PaymentSource{{Name: "Cash"}, {Name: "Card"}},
	}

	_, err = db.Collection("settings").InsertOne(ctx, settings)
	require.NoError(t, err)
}

func mongoConfigFromEnv(t testing.TB) (config.Config, bool) {
	uri := os.Getenv("TEST_MONGO_URI")
	if uri == "" {
		return config.Config{}, false
	}

	u, err := url.Parse(uri)
	if err != nil {
		t.Fatalf("invalid TEST_MONGO_URI: %v", err)
	}

	// The test config only carries host/port, so fail fast instead of silently
	// dropping credentials, TLS options or additional hosts.
	if u.Scheme != "mongodb" || u.User != nil || strings.Contains(u.Host, ",") || u.RawQuery != "" {
		t.Fatalf("TEST_MONGO_URI %q uses features (credentials, mongodb+srv, multiple hosts or options) that are not supported; use a plain mongodb://host:port[/db] URI", uri)
	}

	port := 27017
	if p := u.Port(); p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			t.Fatalf("invalid port in TEST_MONGO_URI: %v", err)
		}
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		dbName = "nutrix_test_" + randomSuffix()
	} else {
		dbName = dbName + "_" + randomSuffix()
	}

	return config.Config{
		Env: "dev",
		Databases: []config.Database{{
			Type:     "mongo",
			Host:     u.Hostname(),
			Port:     port,
			Database: dbName,
			Tables:   map[string]string{"sales": "sales"},
		}},
	}, true
}

func randomSuffix() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
