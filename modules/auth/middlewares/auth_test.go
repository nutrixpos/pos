package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/auth/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func testConfig() config.Config {
	return config.Config{
		Auth: config.AuthConfig{
			Enabled:      true,
			JWTSecret:    "test-secret",
			JWTExpireHrs: 24,
		},
	}
}

func echoClaimsHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(AuthContextKey).(*Claims)
		if !ok {
			http.Error(w, "no claims in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(claims.Username))
	})
}

func generateTestToken(t *testing.T, roles ...string) string {
	t.Helper()
	j := NewJWTUtil("test-secret", 24)
	token, err := j.GenerateToken(models.User{
		ID:       primitive.NewObjectID(),
		Username: "tester",
		Roles:    roles,
	})
	require.NoError(t, err)
	return token
}

func serve(t *testing.T, handler http.Handler, authHeader string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestInternalAuth_AllowAuthenticated_ValidToken(t *testing.T) {
	auth := NewInternalAuth(testConfig(), NewJWTUtil("test-secret", 24))
	handler := auth.AllowAuthenticated(echoClaimsHandler(t))

	rec := serve(t, handler, "Bearer "+generateTestToken(t, "admin"))
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "tester", rec.Body.String())
}

func TestInternalAuth_AllowAuthenticated_MissingHeader(t *testing.T) {
	auth := NewInternalAuth(testConfig(), NewJWTUtil("test-secret", 24))
	handler := auth.AllowAuthenticated(echoClaimsHandler(t))

	rec := serve(t, handler, "")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestInternalAuth_AllowAuthenticated_InvalidToken(t *testing.T) {
	auth := NewInternalAuth(testConfig(), NewJWTUtil("test-secret", 24))
	handler := auth.AllowAuthenticated(echoClaimsHandler(t))

	rec := serve(t, handler, "Bearer not-a-token")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestInternalAuth_AllowAuthenticated_Disabled(t *testing.T) {
	cfg := testConfig()
	cfg.Auth.Enabled = false
	auth := NewInternalAuth(cfg, NewJWTUtil("test-secret", 24))
	handler := auth.AllowAuthenticated(echoClaimsHandler(t))

	// Disabled auth passes the request through without any claims.
	rec := serve(t, handler, "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestInternalAuth_AllowAnyOfRoles(t *testing.T) {
	auth := NewInternalAuth(testConfig(), NewJWTUtil("test-secret", 24))
	handler := auth.AllowAnyOfRoles(echoClaimsHandler(t), "admin", "superuser")

	rec := serve(t, handler, "Bearer "+generateTestToken(t, "admin"))
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = serve(t, handler, "Bearer "+generateTestToken(t, "cashier"))
	assert.Equal(t, http.StatusForbidden, rec.Code)

	rec = serve(t, handler, "")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestNoAuth_PassesThrough(t *testing.T) {
	auth := NewNoAuth(testConfig())
	handler := auth.AllowAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rec := serve(t, handler, "")
	assert.Equal(t, http.StatusTeapot, rec.Code)
}
