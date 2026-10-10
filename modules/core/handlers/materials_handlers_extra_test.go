package handlers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetMaterialsHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "materials",
		models.Material{Id: "m-1", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{{Id: "e1", Quantity: 3}}},
		models.Material{Id: "m-2", Name: "Milk", Unit: "L", Entries: []models.MaterialEntry{{Id: "e2", Quantity: 2}}},
	)

	router := mux.NewRouter()
	router.HandleFunc("/api/materials", GetMaterials(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/materials?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var materials []models.Material
	decodeData(t, rec, &materials)
	assert.Len(t, materials, 2)
}

func TestCalculateMaterialCostHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "materials", models.Material{
		Id:   "m-cost",
		Name: "Flour",
		Entries: []models.MaterialEntry{
			{Id: "e1", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 100},
			{Id: "e2", Quantity: 10, PurchaseQuantity: 10, PurchasePrice: 200},
		},
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/avgcost", CalculateMaterialAverageCost(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/materials/{material_id}/entries/{entry_id}/cost", CalculateMaterialExactCost(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/materials/m-cost/avgcost?quantity=2", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var avgCost float64
	decodeData(t, rec, &avgCost)
	assert.Equal(t, 30.0, avgCost) // ((100/10)+(200/10))/2 * 2

	rec = doRequest(t, router, http.MethodGet, "/api/materials/m-cost/avgcost?quantity=abc", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/materials/m-cost/entries/e1/cost?quantity=2", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var exactCost float64
	decodeData(t, rec, &exactCost)
	assert.Equal(t, 20.0, exactCost)

	rec = doRequest(t, router, http.MethodGet, "/api/materials/m-cost/entries/e1/cost", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/materials/m-cost/entries/e1/cost?quantity=bad", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/materials/m-cost/entries/nope/cost?quantity=2", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteEntryHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "materials", models.Material{
		Id:      "m-del",
		Entries: []models.MaterialEntry{{Id: "e1", Quantity: 1}, {Id: "e2", Quantity: 2}},
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/entries/{entry_id}", DeleteEntry(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodDelete, "/api/materials/m-del/entries/e1", "")
	require.Equal(t, http.StatusNoContent, rec.Code)

	stored, ok := findMaterialDoc(t, env, bson.M{"id": "m-del"})
	require.True(t, ok)
	require.Len(t, stored.Entries, 1)
	assert.Equal(t, "e2", stored.Entries[0].Id)
}

func TestGetMaterialLogsHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "logs", map[string]interface{}{
		"id":          "log-1",
		"type":        models.LogTypeMaterialGRNReceive,
		"material_id": "m-logs",
		"quantity":    5,
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{id}/logs", GetMaterialLogs(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/materials/m-logs/logs?page[number]=1&page[size]=10", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "log-1")
}

func TestEditMaterialHandler_BadJSON(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{id}", EditMaterial(env.Config, env.Logger)).Methods(http.MethodPatch)

	rec := doRequest(t, router, http.MethodPatch, "/api/materials/m-1", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
