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
	"strconv"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
	"github.com/gorilla/mux"
)

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
func Payorder(config config.Config, logger logger.ILogger) http.HandlerFunc {

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
func SubmitOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order models.Order

		request := struct {
			Data models.Order `json:"data"`
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

		filter_displayId := r.URL.Query().Get("filter[display_id]")
		if filter_displayId != "" {
			params.OrderDisplayIdContains = filter_displayId
		}

		filter_isPaid := r.URL.Query().Get("filter[is_paid]")
		if filter_isPaid != "" {
			filter_isPaid_bool, err := strconv.ParseBool(filter_isPaid)
			if err == nil {
				params.FilterIsPaid = filter_isPaid_bool
			}
		}

		filter_isStashed := r.URL.Query().Get("filter[is_stashed]")
		if filter_isStashed != "" {
			filter_isStashed_bool, err := strconv.ParseBool(filter_isStashed)
			if err == nil {
				params.FilterIsStashed = filter_isStashed_bool
			}
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
