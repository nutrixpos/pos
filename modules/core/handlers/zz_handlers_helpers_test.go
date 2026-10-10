package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/auth/middlewares"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/require"
)

// newHandlerEnv creates an isolated database for a handler test.
func newHandlerEnv(t *testing.T) *testutil.TestEnv {
	t.Helper()
	return testutil.NewTestEnv(t, testutil.BackendFromEnv())
}

// testHandlerSettings returns settings that satisfy the services used by the
// handlers (queues, cost method, language and payment sources).
func testHandlerSettings() models.Settings {
	return models.Settings{
		Inventory: models.MaterialSettings{StockAlertTreshold: 1000},
		Orders: models.OrderSettings{
			Queues:                       []models.OrderQueueSettings{{Prefix: "A", Next: 1}},
			DefaultCostCalculationMethod: "average",
		},
		Language:           models.LanguageSettings{Code: "en", Language: "English"},
		AutoOpenCashDrawer: false,
		PaymentSources:     []models.PaymentSource{{Name: "Cash"}, {Name: "Card"}},
	}
}

func jsonBody(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

// decodeData unwraps the {"data": ...} envelope used by the handlers.
func decodeData(t *testing.T, rec *httptest.ResponseRecorder, out interface{}) {
	t.Helper()
	var wrapper struct {
		Data json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &wrapper))
	require.NoError(t, json.Unmarshal(wrapper.Data, out))
}

func decodeMeta(t *testing.T, rec *httptest.ResponseRecorder) JSONAPIMeta {
	t.Helper()
	var wrapper struct {
		Meta JSONAPIMeta `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &wrapper))
	return wrapper.Meta
}

// seed inserts documents into a collection of the isolated test database.
func seed(t *testing.T, env *testutil.TestEnv, collection string, docs ...interface{}) {
	t.Helper()
	coll := env.Client.Database(env.Config.Databases[0].Database).Collection(collection)
	for _, doc := range docs {
		_, err := coll.InsertOne(context.Background(), doc)
		require.NoError(t, err)
	}
}

// withUser injects an authenticated user id into the request context so that
// handlers resolving the user from context get a deterministic value.
func withUser(req *http.Request, userID string) *http.Request {
	ctx := context.WithValue(req.Context(), middlewares.AuthContextKey, &middlewares.Claims{UserID: userID})
	return req.WithContext(ctx)
}

func newRequestWithBody(method, path string, body []byte) *http.Request {
	return httptest.NewRequest(method, path, bytes.NewReader(body))
}

func serve(t *testing.T, router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func countRecipes(t *testing.T, env *testutil.TestEnv) int64 {
	t.Helper()
	n, err := env.Client.Database(env.Config.Databases[0].Database).Collection("recipes").
		CountDocuments(context.Background(), map[string]any{})
	require.NoError(t, err)
	return n
}

// chdirRepoRoot changes the working directory to the repository root so tests
// that read the assets/ directory (language packs) can find it.
func chdirRepoRoot(t *testing.T) {
	t.Helper()
	t.Chdir("../../../")
}
