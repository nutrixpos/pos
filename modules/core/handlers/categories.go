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

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
	PageNumber   int `json:"page_number"`
	PageSize     int `json:"page_size"`
	PageCount    int `json:"page_count"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

// InsertCategory returns a HTTP handler function to insert a Category into the database.
func InsertCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Data models.Category `json:"data"`
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

		err = categoryService.InsertCategory(request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}

// DeleteCategory returns a HTTP handler function to delete a Category from the database.
func DeleteCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err := categoryService.DeleteCategory(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// UpdateCategory returns a HTTP handler function to update a Category in the database.
func UpdateCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body := struct {
			Data models.Category `json:"data"`
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

		category, err := categoryService.UpdateCategory(body.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: category,
			Meta: JSONAPIMeta{
				TotalRecords: 1,
			},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

// GetCategories returns a HTTP handler function to retrieve a list of Categories from the database.
func GetCategories(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		categories, err := categoryService.GetCategories(page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		product_svc := services.RecipeService{
			Logger: logger,
			Config: config,
		}

		for i, category := range categories {
			for j, product := range category.Products {
				product, err := product_svc.GetProduct(product.Id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				categories[i].Products[j] = product
			}
		}

		response := JSONApiOkResponse{
			Data: categories,
			Meta: JSONAPIMeta{
				TotalRecords: len(categories),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
