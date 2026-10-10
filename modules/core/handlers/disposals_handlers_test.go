package handlers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisposalHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "disposals", models.Disposal{Id: "d-1", Type: models.TypeDisposalMaterial, Quantity: 2, OrderId: "o-1"})

	router := mux.NewRouter()
	router.HandleFunc("/api/disposals", GetDisposals(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/disposals", InsertDisposal(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/disposals/{id}", GetDisposal(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/disposals/{id}", UpdateDisposal(env.Config, env.Logger)).Methods(http.MethodPatch)
	router.HandleFunc("/api/disposals/{id}", DeleteDisposal(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodGet, "/api/disposals?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "d-1")

	rec = doRequest(t, router, http.MethodGet, "/api/disposals/d-1", "")
	require.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(t, router, http.MethodPatch, "/api/disposals/d-1",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"type": models.TypeDisposalMaterial, "material_id": "m-1"}}))
	require.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(t, router, http.MethodDelete, "/api/disposals/d-1", "")
	require.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/disposals/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// InsertDisposal decodes into an interface{} so a JSON body never matches the
	// concrete model types and the handler answers 400.
	rec = doRequest(t, router, http.MethodPost, "/api/disposals",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"id": "d-2", "type": models.TypeDisposalMaterial}}))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodPatch, "/api/disposals/d-1", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
