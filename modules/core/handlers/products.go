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
	"net/http"
	"strconv"
	"strings"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

// UpdateProduct returns a HTTP handler function to update a product in the database.
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

// DeleteProduct returns a HTTP handler function to delete a product from the database.
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

// InesrtNewProduct returns a HTTP handler function to insert a new product in the database.
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

// GetProducts returns a HTTP handler function to retrieve a list of products from the database.
// It requires two query string parameters:
// first_index: the index of the first product to be retrieved
// rows: the number of products to be retrieved
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

		search := r.URL.Query().Get("search")

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		args := services.GetProductsParams{
			FirstIndex: first_index,
			Rows:       rows,
			Search:     search,
		}

		products, totalRecords, err := recipeService.GetProducts(args)
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

// GetProductReadyNumber returns a HTTP handler function to retrieve the ready number for a product.
// The product ID is required as query string.
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

// GetRecipeTree returns a HTTP handler function to retrieve a recipe tree.
// The recipe ID is required as query string.
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

// GetRecipeAvailability returns a HTTP handler function to check the availability of multiple recipes.
// The recipe IDs are required as query string, comma separated.
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
