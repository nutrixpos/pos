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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CalculateMaterialCost(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

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

func GetMaterials(config config.Config, logger logger.ILogger) http.HandlerFunc {

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

func AddMaterial(config config.Config, logger logger.ILogger) http.HandlerFunc {
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
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "component adding saved successfully")

	}
}

func DeleteEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry deleted successfully"))

	}
}

func PushMaterialEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {

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

		var componentEntryRequest dto.RequestComponentEntryAdd
		err := json.NewDecoder(r.Body).Decode(&componentEntryRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(componentEntryRequest.ComponentId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		materialService := services.MaterialService{
			Logger: logger,
			Config: config,
		}

		err = materialService.PushMaterialEntry(objectId, componentEntryRequest.Entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}
