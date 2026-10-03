package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func insertMaterialDoc(t *testing.T, env *testutil.TestEnv, material models.Material) {
	t.Helper()
	_, err := env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		InsertOne(context.Background(), material)
	require.NoError(t, err)
}

func findMaterialDoc(t *testing.T, env *testutil.TestEnv, filter bson.M) (models.Material, bool) {
	t.Helper()
	var m models.Material
	err := env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), filter).Decode(&m)
	if err != nil {
		return models.Material{}, false
	}
	return m, true
}

func doRequest(t *testing.T, router *mux.Router, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestGetMaterialEntriesHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	entries := make([]models.MaterialEntry, 0, 5)
	for i := 0; i < 5; i++ {
		entries = append(entries, models.MaterialEntry{Id: fmt.Sprintf("entry-%d", i), Quantity: float64(i)})
	}
	insertMaterialDoc(t, env, models.Material{Id: "mat-1", Entries: entries})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/entries", GetMaterialEntries(env.Config, env.Logger)).Methods("GET")

	rec := doRequest(t, router, http.MethodGet, "/api/materials/mat-1/entries?page%5Bnumber%5D=0&page%5Bsize%5D=2", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data []models.MaterialEntry `json:"data"`
		Meta JSONAPIMeta            `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 5, resp.Meta.TotalRecords)
}

func TestGetMaterialEntriesHandler_NegativePage(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-2", Entries: []models.MaterialEntry{{Id: "e", Quantity: 1}}})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/entries", GetMaterialEntries(env.Config, env.Logger)).Methods("GET")

	rec := doRequest(t, router, http.MethodGet, "/api/materials/mat-2/entries?page%5Bnumber%5D=-1", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAddMaterialHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	router := mux.NewRouter()
	router.HandleFunc("/api/materials", AddMaterial(env.Config, env.Logger)).Methods("POST")

	rec := doRequest(t, router, http.MethodPost, "/api/materials", `{"data":{"id":"mat-new","name":"Flour","unit":"kg"}}`)
	require.Equal(t, http.StatusCreated, rec.Code)

	_, ok := findMaterialDoc(t, env, bson.M{"name": "Flour"})
	assert.True(t, ok, "material should be persisted")
}

func TestAddMaterialHandler_BadJSON(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	router := mux.NewRouter()
	router.HandleFunc("/api/materials", AddMaterial(env.Config, env.Logger)).Methods("POST")

	rec := doRequest(t, router, http.MethodPost, "/api/materials", `{"data": `)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEditMaterialHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-edit", Name: "Old", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{id}", EditMaterial(env.Config, env.Logger)).Methods("PATCH")

	rec := doRequest(t, router, http.MethodPatch, "/api/materials/mat-edit", `{"data":{"name":"New","unit":"L"}}`)
	require.Equal(t, http.StatusOK, rec.Code)

	stored, ok := findMaterialDoc(t, env, bson.M{"id": "mat-edit"})
	require.True(t, ok)
	assert.Equal(t, "New", stored.Name)
	assert.Equal(t, "L", stored.Unit)
}

func TestDeleteMaterialHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-del", Name: "X", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{id}", DeleteMaterial(env.Config, env.Logger)).Methods("DELETE")

	rec := doRequest(t, router, http.MethodDelete, "/api/materials/mat-del", "")
	require.Equal(t, http.StatusNoContent, rec.Code)

	_, ok := findMaterialDoc(t, env, bson.M{"id": "mat-del"})
	assert.False(t, ok, "material should be deleted")
}

func TestCalculateMaterialAverageCostHandler_MissingQuantity(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	router := mux.NewRouter()
	router.HandleFunc("/api/materials/{material_id}/avgcost", CalculateMaterialAverageCost(env.Config, env.Logger)).Methods("GET")

	rec := doRequest(t, router, http.MethodGet, "/api/materials/mat-x/avgcost", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
