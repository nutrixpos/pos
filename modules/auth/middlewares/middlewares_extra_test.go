package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("s3cret")
	require.NoError(t, err)
	assert.True(t, CheckPassword("s3cret", hash))
	assert.False(t, CheckPassword("wrong", hash))
}

func TestInternalAuth_AllowAnyOfRoles_Superuser(t *testing.T) {
	auth := NewInternalAuth(testConfig(), NewJWTUtil("test-secret", 24))
	handler := auth.AllowAnyOfRoles(echoClaimsHandler(t), "admin")

	rec := serve(t, handler, "Bearer "+generateTestToken(t, "superuser"))
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = serve(t, handler, "Bearer not-a-token")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestInternalAuth_AllowAuthenticated_ZitadelInvalidToken(t *testing.T) {
	cfg := testConfig()
	cfg.Zitadel.Enabled = true
	auth := NewInternalAuth(cfg, NewJWTUtil("test-secret", 24))
	handler := auth.AllowAuthenticated(echoClaimsHandler(t))

	rec := serve(t, handler, "Bearer not-a-token")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestNoAuth_AllowAnyOfRoles(t *testing.T) {
	auth := NewNoAuth(testConfig())
	handler := auth.AllowAnyOfRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}), "admin")

	rec := serve(t, handler, "")
	assert.Equal(t, http.StatusTeapot, rec.Code)
}

func TestZitadelAuth_Disabled(t *testing.T) {
	// Disabled config avoids initialising the Zitadel SDK.
	za := NewZitadelAuth(testConfig())
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusTeapot) })

	rec := serve(t, za.AllowAuthenticated(next), "")
	assert.Equal(t, http.StatusTeapot, rec.Code)

	roles := za.AllowAnyOfRoles(next, "admin")
	rec = serve(t, roles, "")
	assert.Equal(t, http.StatusTeapot, rec.Code)
}

func TestZitadelAuth_AllowAnyOfRoles_EnabledNoHeader(t *testing.T) {
	cfg := testConfig()
	cfg.Zitadel.Enabled = true
	za := &ZitadelAuth{Config: cfg}

	handler := za.AllowAnyOfRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}), "admin")

	rec := serve(t, handler, "")
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// No roles requested and no authz: loop is skipped, unauthorized.
	handler = za.AllowAnyOfRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec = serve(t, handler, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
