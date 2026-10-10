package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowCors_OptionsShortCircuits(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	handler := AllowCors(next)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.False(t, called, "next handler must not be called for OPTIONS")
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OPTIONS,DELETE,PATCH", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization, X-Requested-With", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestAllowCors_PassesThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot)
	})
	handler := AllowCors(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.True(t, called, "next handler must be called for non-OPTIONS requests")
	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}
