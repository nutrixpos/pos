package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetComponents(config config.Config, logger logger.ILogger) http.HandlerFunc {

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

		componentService := services.ComponentService{
			Logger: logger,
			Config: config,
		}

		components, err := componentService.GetComponents()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the slice to JSON
		jsonComponents, err := json.Marshal(components)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonComponents)
	}

}

func AddComponent(config config.Config, logger logger.ILogger) http.HandlerFunc {
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
		var dbComponent models.Component
		err := json.NewDecoder(r.Body).Decode(&dbComponent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for index, entry := range dbComponent.Entries {
			dbComponent.Entries[index].PurchaseQuantity = entry.Quantity
		}

		componentService := services.ComponentService{
			Logger: logger,
			Config: config,
		}
		err = componentService.AddComponent(dbComponent)
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
		entryID, err := primitive.ObjectIDFromHex(entryIDStr)
		if err != nil {
			http.Error(w, "Invalid entry ID", http.StatusBadRequest)
			return
		}

		componentIdStr := r.URL.Query().Get("component_id")
		componentId, err := primitive.ObjectIDFromHex(componentIdStr)
		if err != nil {
			http.Error(w, "Invalid entry ID", http.StatusBadRequest)
			return
		}

		componentService := services.ComponentService{
			Logger: logger,
			Config: config,
		}

		err = componentService.DeleteEntry(entryID, componentId)
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

func PushComponentEntry(config config.Config, logger logger.ILogger) http.HandlerFunc {

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

		componentService := services.ComponentService{
			Logger: logger,
			Config: config,
		}

		err = componentService.PushComponentEntry(objectId, componentEntryRequest.Entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}
