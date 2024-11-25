package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

func InsertCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Category models.Category `json:"category"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err = categoryService.InsertCategory(request.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func DeleteCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err := categoryService.DeleteCategory(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func UpdateCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body := struct {
			Category models.Category `json:"category"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err = categoryService.UpdateCategory(body.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetCategories(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		first_index, err := strconv.Atoi(r.URL.Query().Get("first_index"))
		if err != nil {
			first_index = 0
		}

		rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
		if err != nil {
			rows = 9999999
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		categories, err := categoryService.GetCategories(first_index, rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Categories   []models.Category `json:"categories"`
			TotalRecords int               `json:"total_records"`
		}{
			Categories:   categories,
			TotalRecords: len(categories),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
