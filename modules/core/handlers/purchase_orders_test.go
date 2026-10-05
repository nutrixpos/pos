package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreatePurchaseOrderHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-po", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods("POST")

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders",
		`{"data":{"supplier":"Acme","auto_receive":true,"items":[{"material_id":"mat-po","quantity":5,"purchase_price":2}]}}`)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp struct {
		Data models.PurchaseOrder `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, models.PurchaseOrderStatusReceived, resp.Data.Status)
	assert.NotEmpty(t, resp.Data.Id)
}

func TestCreatePurchaseOrderHandler_InvalidBody(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods("POST")

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders", `{"data": `)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReceivePurchaseOrderHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-rec", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods("POST")
	router.HandleFunc("/api/purchase-orders/{id}/receive", ReceivePurchaseOrder(env.Config, env.Logger)).Methods("POST")

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders",
		`{"data":{"supplier":"Acme","items":[{"material_id":"mat-rec","quantity":10,"purchase_price":1}]}}`)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created struct {
		Data models.PurchaseOrder `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))

	// Empty body receives all remaining quantities.
	rec = doRequest(t, router, http.MethodPost, "/api/purchase-orders/"+created.Data.Id+"/receive", "")
	require.Equal(t, http.StatusCreated, rec.Code)

	var grnResp struct {
		Data models.GRN `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &grnResp))
	assert.NotEmpty(t, grnResp.Data.Id)
	assert.Equal(t, created.Data.Id, grnResp.Data.PurchaseOrderId)

	// Receiving again must fail.
	rec = doRequest(t, router, http.MethodPost, "/api/purchase-orders/"+created.Data.Id+"/receive", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelPurchaseOrderHandler(t *testing.T) {
	env := testutil.NewTestEnv(t, testutil.BackendFromEnv())
	insertMaterialDoc(t, env, models.Material{Id: "mat-cancel", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods("POST")
	router.HandleFunc("/api/purchase-orders/{id}", CancelPurchaseOrder(env.Config, env.Logger)).Methods("DELETE")

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders",
		`{"data":{"supplier":"Acme","items":[{"material_id":"mat-cancel","quantity":10,"purchase_price":1}]}}`)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created struct {
		Data models.PurchaseOrder `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))

	rec = doRequest(t, router, http.MethodDelete, "/api/purchase-orders/"+created.Data.Id, "")
	require.Equal(t, http.StatusNoContent, rec.Code)

	var stored models.PurchaseOrder
	err := env.Client.Database(env.Config.Databases[0].Database).Collection("purchase_orders").
		FindOne(context.Background(), bson.M{"id": created.Data.Id}).Decode(&stored)
	require.NoError(t, err)
	assert.Equal(t, models.PurchaseOrderStatusCancelled, stored.Status)
}
