// Package handlers contains HTTP handlers for the core module of nutrix.
//
// The handlers in this package are used to handle incoming HTTP requests for
// the core module of nutrix. They interact with the services package, which
// contains the business logic of the core module.
//
// The handlers in this package create a RESTful API for the core module of
// nutrix. The API endpoints are documented using the Swagger specification.
// Each handler function is responsible for processing HTTP requests, calling
// the appropriate service methods, and returning HTTP responses.
package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
)

func GetOrderLogs(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		order_id_param := params["order_id"]

		order_svc := services.OrderService{
			Config:   config,
			Logger:   logger,
			Settings: settings,
		}

		logs, err := order_svc.GetLogs(order_id_param)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: logs,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func WasteOrderItem(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		order_id_param := params["order_id"]

		decoder := json.NewDecoder(r.Body)

		reason := r.URL.Query().Get("reason")
		if reason == "" {
			http.Error(w, "reason query string is required", http.StatusBadRequest)
			return
		}

		quantityStr := r.URL.Query().Get("quantity")
		if reason == "" {
			http.Error(w, "quantity query string is required", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		request := struct {
			Data struct {
				models.OrderItem `json:"order_item"`
				Other            map[string]interface{}
			} `json:"data"`
		}{}

		err = decoder.Decode(&request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		order_svc := services.OrderService{
			Config:   config,
			Logger:   logger,
			Settings: settings,
		}

		err = order_svc.WasteOrderItem(request.Data.OrderItem, order_id_param, quantity, reason, request.Data.Other)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func RefundOrderItem(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		order_id_param := params["order_id"]
		item_id_param := params["item_id"]

		reason := r.URL.Query().Get("reason")
		if reason == "" {
			http.Error(w, "reason query string is required", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)

		request := struct {
			Data dto.OrderItemRefundRequest `json:"data"`
		}{}

		err := decoder.Decode(&request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		order_svc := services.OrderService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		err = order_svc.RefundItem(order_id_param, item_id_param, reason, request.Data)

		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func PrintKitchenReceipt(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		order, err := orderService.GetOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lang_svc := services.LanguageService{
			Config:   config,
			Logger:   logger,
			Settings: settings,
		}

		acceptLanguage := r.Header.Get("Accept-Language")
		lang := "en"
		if acceptLanguage != "" {
			langs := strings.Split(acceptLanguage, ",")
			if len(langs) > 0 {

				for i := range langs {
					code := strings.TrimSpace(strings.Split(langs[i], ";")[0])
					if len(strings.Split(code, "-")) > 0 {
						code = strings.Split(code, "-")[0]
					}

					code = strings.ToLower(code)
					if _, err := lang_svc.GetLanguage(code); err == nil {
						lang = code
					}
				}
			}
		}

		pwd, err := os.Getwd()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		err = orderService.PrintReceipt(order, pwd+"/assets/core/templates/kitchen_receipt_0.mustache", lang)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PrintClientReceipt(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		order, err := orderService.GetOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lang_svc := services.LanguageService{
			Config:   config,
			Logger:   logger,
			Settings: settings,
		}

		acceptLanguage := r.Header.Get("Accept-Language")
		lang := "en"
		if acceptLanguage != "" {
			langs := strings.Split(acceptLanguage, ",")
			if len(langs) > 0 {

				for i := range langs {
					code := strings.TrimSpace(strings.Split(langs[i], ";")[0])
					if len(strings.Split(code, "-")) > 0 {
						code = strings.Split(code, "-")[0]
					}

					code = strings.ToLower(code)
					if _, err := lang_svc.GetLanguage(code); err == nil {
						lang = code
					}
				}
			}
		}

		pwd, err := os.Getwd()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		err = orderService.PrintReceipt(order, pwd+"/assets/core/templates/order_receipt_0.mustache", lang)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteOrder an http handler to delete an order resource
func DeleteOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.DeleteOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}

}

// Payorder returns a HTTP handler function to pay an unpaid order.
func Payorder(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.PayUnpaidOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetUnpaidOrders returns a HTTP handler function to get all unpaid orders.
func GetUnpaidOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		unpaidOrders, err := orderService.GetUnpaidOrders()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to get unpaid orders", http.StatusInternalServerError)
			return
		}

		response := struct {
			Orders []models.Order `json:"orders"`
		}{
			Orders: unpaidOrders,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal unpaid orders response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

// CancelOrder returns a HTTP handler function to cancel an order.
func CancelOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.CancelOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// FinishOrder returns a HTTP handler function to finish an order.
func FinishOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.FinishOrder(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		display_id, err := orderService.GetOrderDisplayId()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := models.WebsocketOrderFinishServerMessage{
			OrderId: display_id,
			WebsocketTopicServerMessage: models.WebsocketTopicServerMessage{
				Type:      "topic_message",
				TopicName: "order_finished",
				Severity:  "info",
			},
		}

		msgJson, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc, err := services.SpawnNotificationSingletonSvc("melody", logger, config)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc.SendToTopic("order_finished", string(msgJson))

		w.WriteHeader(http.StatusNoContent)
	}
}

// SubmitOrder returns a HTTP handler function to submit an order.
func SubmitOrder(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		acceptLanguage := r.Header.Get("Accept-Language")

		decoder := json.NewDecoder(r.Body)
		var order models.Order

		request := struct {
			Meta models.SubmitOrderMeta `json:"meta"`
			Data models.Order           `json:"data"`
		}{}

		err := decoder.Decode(&request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err = orderService.SubmitOrder(request.Data)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if request.Data.IsAutoStart {
			err = orderService.StartOrder(order.Id, request.Data.Items)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		response := JSONApiOkResponse{
			Data: order,
			Meta: JSONAPIMeta{
				TotalRecords: 1,
			},
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		receipt_svc := services.ReceiptService{
			Config:   config,
			Logger:   logger,
			Settings: settings,
		}

		go func() {

			lang_svc := services.LanguageService{
				Config:   config,
				Logger:   logger,
				Settings: settings,
			}

			lang := "en"
			if acceptLanguage != "" {
				langs := strings.Split(acceptLanguage, ",")
				if len(langs) > 0 {

					for i := range langs {
						code := strings.TrimSpace(strings.Split(langs[i], ";")[0])
						if len(strings.Split(code, "-")) > 0 {
							code = strings.Split(code, "-")[0]
						}

						code = strings.ToLower(code)
						if _, err := lang_svc.GetLanguage(code); err == nil {
							lang = code
						}
					}
				}
			}

			pwd, err := os.Getwd()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			if !order.IsPayLater && request.Meta.IsPrintClientReceipt {
				err = receipt_svc.Print(order, order.Discount, 0, order.SubmittedAt, lang, pwd+"/assets/core/templates/order_receipt_0.mustache")
				if err != nil {
					logger.Error(err.Error())
					return
				}
			}

			if request.Meta.IsPrintKitchenReceipt {
				err = receipt_svc.Print(order, order.Discount, 0, order.SubmittedAt, lang, pwd+"/assets/core/templates/kitchen_receipt_0.mustache")
				if err != nil {
					logger.Error(err.Error())
					return
				}
			}
		}()

		msg := models.WebsocketOrderSubmitServerMessage{
			Order: order,
			WebsocketTopicServerMessage: models.WebsocketTopicServerMessage{
				Type:      "topic_message",
				TopicName: "order_submitted",
				Severity:  "info",
			},
		}

		msgJson, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc, err := services.SpawnNotificationSingletonSvc("melody", logger, config)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc.SendToTopic("order_submitted", string(msgJson))

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}

}

// GetOrders returns a HTTP handler function to retrieve a list of orders.
// to use pagination, send a "first" and "rows" query string
// to select all rows, send a "rows" query string with value -1
// to filter for orders that contains a specific display_id, just send a display_id query string
func GetOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := services.GetOrdersParameters{}

		filter_states := r.URL.Query()["filter[state]"]
		if len(filter_states) > 0 {
			params.FilterState = filter_states
		}

		filter_displayId := r.URL.Query().Get("filter[display_id]")
		if filter_displayId != "" {
			params.OrderDisplayIdContains = filter_displayId
		}

		filter_isPaid := r.URL.Query().Get("filter[is_paid]")
		if filter_isPaid != "" {
			filter_isPaid_bool, err := strconv.ParseBool(filter_isPaid)
			if err == nil {
				if filter_isPaid_bool {
					params.FilterIsPaid = 1
				} else {
					params.FilterIsPaid = 0
				}
			}
		} else {
			params.FilterIsPaid = -1
		}

		filter_isPaylater := r.URL.Query().Get("filter[is_pay_later]")
		if filter_isPaylater != "" {
			filter_isPayLater_bool, err := strconv.ParseBool(filter_isPaylater)
			if err == nil {
				if filter_isPayLater_bool {
					params.IsPayLater = 1
				} else {
					params.IsPayLater = 0
				}
			}
		} else {
			params.IsPayLater = -1
		}

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			params.PageNumber = 1
		} else {
			params.PageNumber = page_number
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			params.PageSize = 50
		} else {
			params.PageSize = page_size
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		orders, total_records, err := orderService.GetOrders(params)

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: orders,
			Meta: JSONAPIMeta{
				TotalRecords: int(total_records),
				PageNumber:   params.PageNumber,
				PageSize:     params.PageSize,
				PageCount:    int(math.Ceil(float64(total_records) / float64(params.PageSize))),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

// StartOrder returns a HTTP handler function to start an order.
func StartOrder(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		request_body := struct {
			Data []models.OrderItem `json:"data"`
		}{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request_body)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		err = orderService.StartOrder(id_param, request_body.Data)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			response := struct {
				Data string `json:"body"`
			}{
				Data: err.Error(),
			}

			json_response, err := json.Marshal(response)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(json_response)
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

// GetOrder returns a HTTP handler function to retrieve an order.
func GetOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err := orderService.GetOrder(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: order,
		}

		jsonOrder, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrder)
	}
}
