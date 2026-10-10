package handlers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurchaseOrderReadHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	insertMaterialDoc(t, env, models.Material{Id: "mat-po-h", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/purchase-orders", GetPurchaseOrders(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/purchase-orders/{id}", GetPurchaseOrder(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/grns", GetGRNs(env.Config, env.Logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/grns/{id}", GetGRN(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders",
		jsonBody(t, map[string]interface{}{"data": models.PurchaseOrder{
			Supplier:    "Acme",
			AutoReceive: true,
			Items:       []models.PurchaseOrderItem{{MaterialId: "mat-po-h", Quantity: 5, PurchasePrice: 2}},
		}}))
	require.Equal(t, http.StatusCreated, rec.Code)
	var po models.PurchaseOrder
	decodeData(t, rec, &po)
	require.NotEmpty(t, po.Id)

	rec = doRequest(t, router, http.MethodGet, "/api/purchase-orders?page[number]=1&page[size]=10&filter[search]=Acme&filter[status]=received", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var list []models.PurchaseOrder
	decodeData(t, rec, &list)
	assert.NotEmpty(t, list)

	rec = doRequest(t, router, http.MethodGet, "/api/purchase-orders/"+po.Id, "")
	require.Equal(t, http.StatusOK, rec.Code)
	rec = doRequest(t, router, http.MethodGet, "/api/purchase-orders/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	rec = doRequest(t, router, http.MethodGet, "/api/grns?page[number]=1&page[size]=10&filter[purchase_order_id]="+po.Id, "")
	require.Equal(t, http.StatusOK, rec.Code)
	var grns []models.GRN
	decodeData(t, rec, &grns)
	require.NotEmpty(t, grns)

	rec = doRequest(t, router, http.MethodGet, "/api/grns/"+grns[0].Id, "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestCancelPurchaseOrderHandler_Success(t *testing.T) {
	env := newHandlerEnv(t)
	insertMaterialDoc(t, env, models.Material{Id: "mat-po-c", Name: "Flour", Unit: "kg", Entries: []models.MaterialEntry{}})

	router := mux.NewRouter()
	router.HandleFunc("/api/purchase-orders", CreatePurchaseOrder(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/purchase-orders/{id}/cancel", CancelPurchaseOrder(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/purchase-orders",
		jsonBody(t, map[string]interface{}{"data": models.PurchaseOrder{
			Supplier: "Acme",
			Items:    []models.PurchaseOrderItem{{MaterialId: "mat-po-c", Quantity: 5, PurchasePrice: 2}},
		}}))
	require.Equal(t, http.StatusCreated, rec.Code)
	var po models.PurchaseOrder
	decodeData(t, rec, &po)

	rec = doRequest(t, router, http.MethodPost, "/api/purchase-orders/"+po.Id+"/cancel", "")
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/purchase-orders/missing/cancel", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
