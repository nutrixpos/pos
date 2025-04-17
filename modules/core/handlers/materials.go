// Package handlers contains HTTP handlers for the core module of nutrix.
//
// The handlers in this package are used to handle incoming HTTP requests for
// the core module of nutrix. The handlers are used to interact with the services
// package, which contains the business logic of the core module.
//
// The handlers in this package are used to create a RESTful API for the core
// module of nutrix. The API endpoints are documented using the Swagger
// specification.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
)

// CalculateMaterialCost returns a HTTP handler function to calculate the cost of a material entry.
//
// This handler retrieves the entry ID, material ID, and quantity from the query string,
// and uses the MaterialService to calculate and return the cost.
func CalculateMaterialCost(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		material_id_param := params["material_id"]
		entry_id_param := params["entry_id"]

		quantityStr := r.URL.Query().Get("quantity")
		if quantityStr == "" {
			http.Error(w, "quantity query string is required", http.StatusBadRequest)
			return
		}

		var quantity float64
		quantity, err := strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		cost, err := materialService.CalculateMaterialCost(entry_id_param, material_id_param, quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: cost,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

// GetMaterials returns a HTTP handler function to retrieve a list of materials from the database.
func GetMaterials(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		materials, err := materialService.GetMaterials(page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: materials,
			Meta: JSONAPIMeta{
				TotalRecords: len(materials),
			},
		}

		// Convert the slice to JSON
		jsonMaterials, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonMaterials)
	}

}

// AddMaterial returns a HTTP handler function to add a new material to the database.
func AddMaterial(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parse the request body into a DBComponent struct

		request := struct {
			Data models.Material `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for index, entry := range request.Data.Entries {
			request.Data.Entries[index].PurchaseQuantity = entry.Quantity
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}
		err = materialService.AddComponent(request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a success response
		fmt.Fprint(w, "component adding saved successfully")

		w.WriteHeader(http.StatusCreated)

	}
}

// DeleteEntry returns a HTTP handler function to delete an entry in the material database.
func DeleteEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		material_id_param := params["material_id"]
		entry_id_param := params["entry_id"]

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err := materialService.DeleteEntry(entry_id_param, material_id_param)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteMaterial(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id_param := params["id"]

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err := materialService.DeleteMaterial(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// EditMaterial returns a HTTP handler function to edit an existing material in the database.
func EditMaterial(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		request := struct {
			Data models.Material `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err = materialService.EditMaterial(id_param, request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// PushMaterialEntry returns a HTTP handler function to add a new entry to a material in the database.
func PushMaterialEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		material_id := params["id"]

		request := struct {
			Data []models.MaterialEntry `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err = materialService.PushMaterialEntry(material_id, request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}

func GetMaterialLogs(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		material_id := params["id"]

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		logService := services.LogService{
			Logger: logger,
			Config: config,
		}

		logs, total_records, err := logService.GetComponentLogs(material_id, page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: logs,
			Meta: JSONAPIMeta{
				TotalRecords: int(total_records),
			},
		}

		jsonLogs, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonLogs)
	}
}
