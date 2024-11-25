package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

func UpdateProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		err = recipeService.UpdateProduct(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func DeleteProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		err := recipeService.DeleteProduct(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func InesrtNewProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		err = recipeService.InsertNew(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetProducts(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		first_index, err := strconv.Atoi(r.URL.Query().Get("first_index"))
		if err != nil {
			first_index = 0
		}

		rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
		if err != nil {
			rows = 7
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		products, totalRecords, err := recipeService.GetProducts(first_index, rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			TotalRecords int64            `json:"total_records"`
			Products     []models.Product `json:"products"`
		}{
			TotalRecords: totalRecords,
			Products:     products,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetProductReadyNumber(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		ready, err := recipeService.GetReadyNumber(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ready); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetRecipeTree(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "recipe id query string is required", http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		tree, err := recipeService.GetRecipeTree(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tree); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetRecipeAvailability(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("ids")
		if id == "" {
			http.Error(w, "recipe ids comma separated is required as query string", http.StatusBadRequest)
			return
		}

		ids := strings.Split(id, `,`)

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		availabilities, err := recipeService.CheckRecipesAvailability(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(availabilities); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
