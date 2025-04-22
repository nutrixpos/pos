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

// UpdateDisposal returns a HTTP handler function to update a disposal in the database.
func UpdateDisposal(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		request := struct {
			Data interface{} `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		disposal_svc := services.DisposalService{
			Logger: logger,
			Config: config,
		}

		updated_disposal, err := disposal_svc.UpdateDisposal(id_param, request.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: updated_disposal,
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

func GetDisposal(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id_param := params["id"]

		disposalService := services.DisposalService{
			Logger: logger,
			Config: config,
		}

		disposal, err := disposalService.GetDisposal(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: disposal,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteDisposal returns a HTTP handler function to delete a disposal from the database.
func DeleteDisposal(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id_param := params["id"]

		disposalService := services.DisposalService{
			Logger: logger,
			Config: config,
		}

		err := disposalService.DeleteDisposal(id_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// InsertDisposal returns a HTTP handler function to insert a new disposal in the database.
func InsertDisposal(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Data interface{} `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if material_disposal, ok := request.Data.(models.MaterialDisposal); ok {
			recipeService := services.DisposalService{
				Logger: logger,
				Config: config,
			}

			err = recipeService.AddMaterialDisposal(material_disposal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		if product_disposal, ok := request.Data.(models.ProductDisposal); ok {
			recipeService := services.DisposalService{
				Logger: logger,
				Config: config,
			}

			err = recipeService.AddProductDisposal(product_disposal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetDisposals returns a HTTP handler function to retrieve a list of disposals from the database.
func GetDisposals(config config.Config, logger logger.ILogger) http.HandlerFunc {

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

		disposal_service := services.DisposalService{
			Logger: logger,
			Config: config,
		}

		args := services.GetDisposalsParameters{
			PageNumber: page_number,
			PageSize:   page_size,
			Search:     search,
		}

		disposals, totalRecords, err := disposal_service.GetDisposals(args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: disposals,
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
