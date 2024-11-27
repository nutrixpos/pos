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

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

// CalculateMaterialCost returns a HTTP handler function to calculate the cost of a material entry.
//
// This handler retrieves the entry ID, material ID, and quantity from the query string,
// and uses the MaterialService to calculate and return the cost.
func CalculateMaterialCost(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Retrieve the entry ID from the query string
		entryIDStr := r.URL.Query().Get("entry_id")
		if entryIDStr == "" {
			http.Error(w, "entry_id query string is required", http.StatusBadRequest)
			return
		}

		materialIdStr := r.URL.Query().Get("material_id")
		if materialIdStr == "" {
			http.Error(w, "material_id query string is required", http.StatusBadRequest)
			return
		}

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

		cost, err := materialService.CalculateMaterialCost(entryIDStr, materialIdStr, quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("%f", cost)))
	}
}

// GetMaterials returns a HTTP handler function to retrieve a list of materials from the database.
func GetMaterials(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		materials, err := materialService.GetMaterials()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the slice to JSON
		jsonMaterials, err := json.Marshal(materials)
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
		var dbComponent models.Material
		err := json.NewDecoder(r.Body).Decode(&dbComponent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for index, entry := range dbComponent.Entries {
			dbComponent.Entries[index].PurchaseQuantity = entry.Quantity
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}
		err = materialService.AddComponent(dbComponent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a success response
		fmt.Fprint(w, "component adding saved successfully")

	}
}

// DeleteEntry returns a HTTP handler function to delete an entry in the material database.
func DeleteEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the entry ID from the query string
		entryIDStr := r.URL.Query().Get("entry_id")
		componentIdStr := r.URL.Query().Get("component_id")

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err := materialService.DeleteEntry(entryIDStr, componentIdStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}

		// Send a success response
		w.Write([]byte("Entry deleted successfully"))

	}
}

// EditMaterial returns a HTTP handler function to edit an existing material in the database.
func EditMaterial(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var material_edit_request dto.MaterialEditRequest
		err := json.NewDecoder(r.Body).Decode(&material_edit_request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err = materialService.EditMaterial(material_edit_request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
}

// PushMaterialEntry returns a HTTP handler function to add a new entry to a material in the database.
func PushMaterialEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var componentEntryRequest dto.RequestComponentEntryAdd
		err := json.NewDecoder(r.Body).Decode(&componentEntryRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err = materialService.PushMaterialEntry(componentEntryRequest.ComponentId, componentEntryRequest.Entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

}
