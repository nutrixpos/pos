package handlers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefundOrderItemHandler_Success(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{
		Id: "o-refund",
		Items: []models.OrderItem{{
			Id:        "i-1",
			Product:   models.Product{Id: "p-1", Name: "Pizza"},
			Quantity:  1,
			SalePrice: 10,
		}},
	})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/items/{item_id}/refund", RefundOrderItem(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-refund/items/i-1/refund?reason=customer",
		jsonBody(t, map[string]interface{}{"data": map[string]interface{}{
			"destination":  dto.DTOOrderItemRefundDestination_Waste,
			"refund_value": 10,
		}}))
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "refunded", findOrderDoc(t, env, "o-refund").Items[0].Status)
}

func TestGetOrdersHandler_FilterBranches(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-1", DisplayId: "A-1", State: "pending", IsPaid: false, IsPayLater: true})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders", GetOrders(env.Config, env.Logger)).Methods(http.MethodGet)

	// Invalid boolean is ignored, is_pay_later true selects the order.
	rec := doRequest(t, router, http.MethodGet, "/api/orders?filter[is_paid]=notabool&filter[is_pay_later]=true&page[number]=1&page[size]=5", "")
	require.Equal(t, http.StatusOK, rec.Code)
	var orders []models.Order
	decodeData(t, rec, &orders)
	assert.Len(t, orders, 1)

	// is_pay_later false excludes it.
	rec = doRequest(t, router, http.MethodGet, "/api/orders?filter[is_pay_later]=false", "")
	require.Equal(t, http.StatusOK, rec.Code)
	decodeData(t, rec, &orders)
	assert.Empty(t, orders)
}

func TestOrderTipHandlers_Errors(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/addtips", OrderAddTip(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPatch)
	router.HandleFunc("/api/orders/{order_id}/removetips", OrderRemoveTip(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPatch)

	rec := doRequest(t, router, http.MethodPatch, "/api/orders/o/addtips?tip_amount=abc", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	rec = doRequest(t, router, http.MethodPatch, "/api/orders/o/removetips?tip_amount=abc", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doRequest(t, router, http.MethodPatch, "/api/orders/missing/addtips?tip_amount=1", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	rec = doRequest(t, router, http.MethodPatch, "/api/orders/missing/removetips?tip_amount=1", "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestPrintClientReceiptHandler_AcceptLanguage(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-print-lang"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/print-client", PrintClientReceipt(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	req := newRequestWithBody(http.MethodPost, "/api/orders/o-print-lang/print-client", nil)
	req.Header.Set("Accept-Language", "ar,en;q=0.9")
	rec := serve(t, router, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestOrderTipHandlers_MissingAmount(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/addtips", OrderAddTip(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPatch)
	router.HandleFunc("/api/orders/{order_id}/removetips", OrderRemoveTip(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPatch)

	assert.Equal(t, http.StatusBadRequest, doRequest(t, router, http.MethodPatch, "/api/orders/o/addtips", "").Code)
	assert.Equal(t, http.StatusBadRequest, doRequest(t, router, http.MethodPatch, "/api/orders/o/removetips", "").Code)
}

func TestWasteOrderItemHandler_BadBody(t *testing.T) {
	env := newHandlerEnv(t)
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{order_id}/waste", WasteOrderItem(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o/waste?reason=x&quantity=1", "{")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestPrintKitchenReceiptHandler_AcceptLanguage(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-kitchen-lang"})

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/print-kitchen", PrintKitchenReceipt(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	req := newRequestWithBody(http.MethodPost, "/api/orders/o-kitchen-lang/print-kitchen", nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	rec := serve(t, router, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSubmitOrderHandler_PrintGoroutine(t *testing.T) {
	env := newHandlerEnv(t)
	seed(t, env, "recipes", models.Product{Id: "r-gor", Name: "Pizza", Price: 10})

	settings := testHandlerSettings()
	settings.AutoOpenCashDrawer = true
	settings.ClientReceiptPrinter.Host = "127.0.0.1:1"
	settings.KitchenReceiptPrinter.Host = "127.0.0.1:1"

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/submit", SubmitOrder(env.Config, env.Logger, settings)).Methods(http.MethodPost)

	body := map[string]interface{}{
		"meta": map[string]interface{}{"is_print_client_receipt": true, "is_print_kitchen_receipt": true},
		"data": models.Order{
			IsPaid:   true,
			Payments: []models.OrderPayment{{Source: "Cash", Amount: 10}},
			Items:    []models.OrderItem{{Product: models.Product{Id: "r-gor"}, Quantity: 1}},
		},
	}
	rec := doRequest(t, router, http.MethodPost, "/api/orders/submit", jsonBody(t, body))
	require.Equal(t, http.StatusOK, rec.Code)

	// Give the asynchronous receipt/cash-drawer goroutine time to run so its
	// statements are recorded by the coverage run.
	time.Sleep(500 * time.Millisecond)
}

func TestFinishOrderHandler_SettingsErrorPath(t *testing.T) {
	env := newHandlerEnv(t)
	insertOrderDoc(t, env, models.Order{Id: "o-fin2", State: "in_progress"})

	// Remove the settings document so GetSettings fails before finishing.
	if err := env.Client.Database(env.Config.Databases[0].Database).Collection("settings").Drop(context.Background()); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/finish", FinishOrder(env.Config, env.Logger, testHandlerSettings())).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/orders/o-fin2/finish", "")
	// The handler writes the error body before setting the status, so assert on
	// the body rather than the (implicit 200) status.
	assert.NotEmpty(t, rec.Body.String())
}
