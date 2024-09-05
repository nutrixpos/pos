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

func FinishOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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

		w.WriteHeader(http.StatusOK)

	}
}

func SubmitOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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

		w.WriteHeader(http.StatusOK)

	}

}

func GetOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		orders, err := orderService.GetOrders()

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}

func StartOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var order_start_request dto.OrderStartRequest2
		err := decoder.Decode(&order_start_request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.StartOrder(order_start_request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}

func GetOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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
