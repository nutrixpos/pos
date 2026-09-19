package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func userIDFromContext(config config.Config, r *http.Request) string {
	if config.Zitadel.Enabled {
		return r.Context().Value("auth_ctx").(oidc.IntrospectionResponse).Subject
	}
	return "0"
}

// CreatePurchaseOrder returns a HTTP handler function to create a purchase order.
// When the order is created with auto_receive=true the goods are received
// immediately, generating a GRN and adding the quantities to inventory.
func CreatePurchaseOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id := userIDFromContext(config, r)

		request := struct {
			Data models.PurchaseOrder `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		po, err := purchaseOrderService.CreatePurchaseOrder(request.Data, user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := JSONApiOkResponse{
			Data: po,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetPurchaseOrders returns a HTTP handler function to list purchase orders.
func GetPurchaseOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		search := r.URL.Query().Get("filter[search]")
		status := r.URL.Query().Get("filter[status]")

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		pos, totalRecords, err := purchaseOrderService.GetPurchaseOrders(services.GetPurchaseOrdersParams{
			PageNumber: page_number,
			PageSize:   page_size,
			Search:     search,
			Status:     status,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: pos,
			Meta: JSONAPIMeta{
				TotalRecords: int(totalRecords),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetPurchaseOrder returns a HTTP handler function to get a single purchase order.
func GetPurchaseOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		purchase_order_id := params["id"]

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		po, err := purchaseOrderService.GetPurchaseOrder(purchase_order_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: po,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CancelPurchaseOrder returns a HTTP handler function to cancel a purchase order.
func CancelPurchaseOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		purchase_order_id := params["id"]

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		err := purchaseOrderService.CancelPurchaseOrder(purchase_order_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ReceivePurchaseOrder returns a HTTP handler function to receive goods against
// a purchase order. The optional body contains per-item received quantities;
// an empty body receives all remaining quantities.
func ReceivePurchaseOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id := userIDFromContext(config, r)

		params := mux.Vars(r)
		purchase_order_id := params["id"]

		request := struct {
			Data []services.ReceiveItem `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil && !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		grn, err := purchaseOrderService.ReceivePurchaseOrder(purchase_order_id, user_id, request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := JSONApiOkResponse{
			Data: grn,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetGRNs returns a HTTP handler function to list goods received notes.
func GetGRNs(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		search := r.URL.Query().Get("filter[search]")
		purchase_order_id := r.URL.Query().Get("filter[purchase_order_id]")

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		grns, totalRecords, err := purchaseOrderService.GetGRNs(services.GetGRNsParams{
			PageNumber:      page_number,
			PageSize:        page_size,
			Search:          search,
			PurchaseOrderId: purchase_order_id,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: grns,
			Meta: JSONAPIMeta{
				TotalRecords: int(totalRecords),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetGRN returns a HTTP handler function to get a single goods received note.
func GetGRN(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		grn_id := params["id"]

		purchaseOrderService := services.PurchaseOrderService{
			Logger: logger,
			Config: config,
		}

		grn, err := purchaseOrderService.GetGRN(grn_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: grn,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}