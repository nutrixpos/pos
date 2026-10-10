package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-1", DisplayId: "A-1", State: "pending"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}", GetOrder(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/orders/o-1", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var o models.Order
	decodeData(t, rec, &o)
	assert.Equal(t, "A-1", o.DisplayId)

	rec = doRequest(t, router, http.MethodGet, "/api/orders/missing", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-del"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}", DeleteOrder(env.Config, env.Logger)).Methods(http.MethodDelete)

	rec := doRequest(t, router, http.MethodDelete, "/api/orders/o-del", "")
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCancelOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-cancel", State: "pending"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/cancel", CancelOrder(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-cancel/cancel", "")
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "cancelled", findOrderDoc(t, env, "o-cancel").State)
}

func TestGetUnpaidOrdersHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-unpaid", IsPayLater: true, IsPaid: false, State: "pending"})
	insertOrderDoc(t, env, models.Order{Id: "o-paid", IsPayLater: true, IsPaid: true, State: "pending"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/unpaid", GetUnpaidOrders(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/orders/unpaid", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Orders []models.Order `json:"orders"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Orders, 1)
	assert.Equal(t, "o-unpaid", resp.Orders[0].Id)
}

func TestPayOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-pay", SalePrice: 10, IsPaid: false})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/pay", Payorder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	// mismatched total -> bad request
	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-pay/pay",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"payments": []models.OrderPayment{{Source: "Cash", Amount: 3}}}}))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// exact total -> paid
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-pay/pay",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"payments": []models.OrderPayment{{Source: "Cash", Amount: 10}}}}))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// already paid -> conflict
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-pay/pay",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"payments": []models.OrderPayment{{Source: "Cash", Amount: 10}}}}))
	assert.Equal(t, http.StatusConflict, rec.Code)

	// bad json
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-pay/pay", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrderLogsHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "logs", models.LogWasteOrderItem{
		Log:      models.Log{Id: "log-1", Type: "waste_orderitem"},
		OrderId:  "o-logs",
		Quantity: 1,
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/logs", GetOrderLogs(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/orders/o-logs/logs", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "waste_orderitem")
}

func TestWasteOrderItemHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-waste"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/waste", WasteOrderItem(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-waste/waste?reason=burnt&quantity=1",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{"order_item": models.OrderItem{Id: "i-1"}}}))
	assert.Equal(t, http.StatusOK, rec.Code)

	// missing reason
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-waste/waste?quantity=1", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// invalid quantity
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-waste/waste?reason=x&quantity=abc", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateOrderCustomDataHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-custom"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/custom", UpdateOrderCustomData(env.Config, env.Logger)).Methods(http.MethodPatch)

	rec := doRequest(t, router, http.MethodPatch, "/api/orders/o-custom/custom",
		jsonBody(t, map[string]interface{}{"data": map[string]string{"table": "5"}}))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(t, router, http.MethodPatch, "/api/orders/o-custom/custom", "{")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrdersHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-1", DisplayId: "A-1", State: "pending", IsPaid: false})
	insertOrderDoc(t, env, models.Order{Id: "o-2", DisplayId: "A-2", State: "finished", IsPaid: true})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders", GetOrders(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/orders?page[number]=1&page[size]=10&filter[state]=pending&filter[is_paid]=false&filter[is_pay_later]=false&filter[display_id]=A-1", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var orders []models.Order
	decodeData(t, rec, &orders)
	assert.Len(t, orders, 1)

	// defaults + negative state filter
	rec = doRequest(t, router, http.MethodGet, "/api/orders?filter[state]=!finished", "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestStartOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-start", State: "pending"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/start", StartOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-start/start",
		jsonBody(t, map[string]interface{}{"data": []models.OrderItem{{Id: "i-1", Quantity: 1}}}))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "in_progress", findOrderDoc(t, env, "o-start").State)

	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-start/start", "{")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSubmitOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-1", Name: "Pizza", Price: 10})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/submit", SubmitOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	order := models.Order{
		IsPayLater: true,
		Items:      []models.OrderItem{{Product: models.Product{Id: "r-1"}, Quantity: 2}},
	}
	rec := doRequest(t, router, http.MethodPost, "/api/orders/submit",
		jsonBody(t, map[string]interface{}{"data": order}))
	require.Equal(t, http.StatusOK, rec.Code)

	var created models.Order
	decodeData(t, rec, &created)
	assert.NotEmpty(t, created.Id)
	assert.Equal(t, 20.0, created.SalePrice)
}

func TestFinishOrderHandler(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-1", Name: "Pizza", Price: 10})
	insertOrderDoc(t, env, models.Order{
		Id:    "o-fin",
		State: "in_progress",
		Items: []models.OrderItem{{
			Id:       "i-1",
			Product:  models.Product{Id: "r-1", EnableInventoryConsumption: false},
			Quantity: 1,
		}},
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/finish", FinishOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-fin/finish", "")
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "finished", findOrderDoc(t, env, "o-fin").State)
}

func TestPrintReceiptHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-print"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/print-client", PrintClientReceipt(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)
	router.HandleFunc("/api/orders/{id}/print-kitchen", PrintKitchenReceipt(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	// No printer reachable: the handlers get the order then fail at print time.
	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-print/print-client", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/orders/missing/print-kitchen", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Existing order through the kitchen receipt path (fails at print time).
	rec = doRequest(t, router, http.MethodPost, "/api/orders/o-print/print-kitchen", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRefundOrderItemHandler_Validation(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/items/{item_id}/refund", RefundOrderItem(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o/items/i/refund", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodPost, "/api/orders/o/items/i/refund?reason=x", "{")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestStartOrderHandler_MissingOrder(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/start", StartOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/nope/start",
		jsonBody(t, map[string]interface{}{"data": []models.OrderItem{}}))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSubmitOrderHandler_AutoStartFinish(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-2", Name: "Pizza", Price: 10})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/submit", SubmitOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	order := models.Order{
		IsAutoStart:  true,
		IsAutoFinish: true,
		IsPaid:       true,
		Payments:     []models.OrderPayment{{Source: "Cash", Amount: 10}},
		Items:        []models.OrderItem{{Product: models.Product{Id: "r-2"}, Quantity: 1}},
	}
	rec := doRequest(t, router, http.MethodPost, "/api/orders/submit",
		jsonBody(t, map[string]interface{}{"data": order}))
	require.Equal(t, http.StatusOK, rec.Code)

	var created models.Order
	decodeData(t, rec, &created)
	assert.Equal(t, "finished", findOrderDoc(t, env, created.Id).State)
}

func TestUserIDFromContext(t *testing.T) {
	req := newRequestWithBody(http.MethodGet, "/", nil)
	assert.Equal(t, "0", userIDFromContext(config.Config{}, req))

	authReq := withUser(req, "user-42")
	cfg := config.Config{}
	cfg.Auth.Enabled = true
	assert.Equal(t, "user-42", userIDFromContext(cfg, authReq))
}
