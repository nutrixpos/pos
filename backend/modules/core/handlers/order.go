package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

func PayUnpaidOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.URL.Query().Get("id")
		if orderId == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.PayUnpaidOrder(orderId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

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

func CancelOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id quert string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.CancelOrder(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetStashedOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		stashed_orders, err := orderService.GetStashedOrders()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(stashed_orders)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func OrderRemoveFromStash(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_remove_stash_request dto.OrderRemoveStashRequest
		err := decoder.Decode(&order_remove_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.RemoveStashedOrder(order_remove_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func OrderStash(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_stash_request dto.OrderStashRequest
		err := decoder.Decode(&order_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err := orderService.StashOrder(order_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := dto.OrderStashResponse{
			Order: order,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func FinishOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var finish_order_request dto.FinishOrderRequest
		err := decoder.Decode(&finish_order_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.FinishOrder(finish_order_request)
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

	}
}

func SubmitOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order models.Order
		err := decoder.Decode(&order)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.SubmitOrder(order)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

func GetOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := services.GetOrdersParameters{}

		orderDisplayId := r.URL.Query().Get("display_id")
		if orderDisplayId != "" {
			params.OrderDisplayIdContains = orderDisplayId
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		orders, err := orderService.GetOrders(params)

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := struct {
			Orders       []models.Order `json:"orders"`
			TotalRecords int            `json:"total_records"`
		}{
			Orders:       orders,
			TotalRecords: len(orders),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func StartOrder(config config.Config, logger logger.ILogger, settings config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_start_request dto.OrderStartRequest
		err := decoder.Decode(&order_start_request)
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

		err = orderService.StartOrder(order_start_request)
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

func GetOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id quert string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err := orderService.GetOrder(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrder)
	}
}
