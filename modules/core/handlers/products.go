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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/helpers"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
	"github.com/gorilla/mux"
)

func UpdateProductImage(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form data
		err := r.ParseMultipartForm(32 << 20) // Max file size: 32MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := mux.Vars(r)
		id_param := params["id"]

		product_svc := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		product, err := product_svc.GetProduct(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the uploaded file
		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			logger.Error("Error uploading file %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		names := strings.Split(fileHeader.Filename, ".")
		file_extension := ""

		if len(names) > 0 {
			file_extension = "." + names[len(names)-1]
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		newpath := filepath.Join(".", "public")
		err = os.MkdirAll(newpath, os.ModePerm)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		random_string := helpers.RandStringBytesMaskImprSrc(20)

		// Create a new file on the server
		dst, err := os.Create(config.UploadsPath + "/" + random_string + file_extension)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the server file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete the old image file from the public directory
		oldImagePath := config.UploadsPath + "/" + product.ImageURL
		err = os.Remove(oldImagePath)
		if err != nil {
			logger.Error("Error deleting old image file %s", err.Error())
			// Optionally handle the error, but don't return it to the client
		}

		product.ImageURL = random_string + file_extension

		product_svc.UpdateProduct(id_param, product)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateProduct returns a HTTP handler function to update a product in the database.
func UpdateProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		request := struct {
			Data models.Product `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		err = recipeService.UpdateProduct(id_param, request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updated_product, err := recipeService.GetProduct(id_param)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: updated_product,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

// DeleteProduct returns a HTTP handler function to delete a product from the database.
func DeleteProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		err := recipeService.DeleteProduct(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// InesrtNewProduct returns a HTTP handler function to insert a new product in the database.
func InesrtNewProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Data models.Product `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		new_product, err := recipeService.InsertNew(request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: new_product,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// GetProduct gets a single product from the db
func GetProduct(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		product, err := recipeService.GetProduct(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: product,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
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

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		search := r.URL.Query().Get("filter[search]")

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		args := services.GetProductsParams{
			PageNumber: page_number,
			PageSize:   page_size,
			Search:     search,
		}

		products, totalRecords, err := recipeService.GetProducts(args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: products,
			Meta: JSONAPIMeta{
				TotalRecords: int(totalRecords),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetRecipeTree returns a HTTP handler function to retrieve a recipe tree.
// The recipe ID is required as query string.
func GetRecipeTree(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		recipeService := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		tree, err := recipeService.GetRecipeTree(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: tree,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
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

		response := JSONApiOkResponse{
			Data: availabilities,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}

}
