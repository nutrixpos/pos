package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
)

func DeleteCustomer(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		customers_svc := services.CustomersService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		err := customers_svc.DeleteCustomer(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetCustomers(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := services.GetCustomersParams{}

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

		customers_svc := services.CustomersService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		customers, total_records, err := customers_svc.GetCustomers(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := JSONApiOkResponse{
			Data: customers,
			Meta: JSONAPIMeta{
				TotalRecords: total_records,
			},
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func AddCustomer(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Data models.Customer `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		customers_svc := services.CustomersService{
			Logger: logger,
			Config: config,
		}

		addedCustomer, err := customers_svc.InsertNew(request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: addedCustomer,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)

	}
}

func UpdateCustomer(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id_param := params["id"]

		request := struct {
			Data models.Customer `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		customers_svc := services.CustomersService{
			Logger: logger,
			Config: config,
		}

		updatedCustomer, err := customers_svc.UpdateCustomer(request.Data, id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: updatedCustomer,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func GetCustomer(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id_param := params["id"]

		customers_svc := services.CustomersService{
			Logger: logger,
			Config: config,
		}

		customer, err := customers_svc.GetCustomer(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: customer,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}
