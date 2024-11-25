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

func GetSalesPerDay(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		first_index, err := strconv.Atoi(r.URL.Query().Get("first_index"))
		if err != nil {
			first_index = 0
		}

		rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
		if err != nil {
			rows = 7
		}

		salesService := services.SalesService{
			Logger: logger,
			Config: config,
		}

		sales, totalRecords, err := salesService.GetSalesPerday(first_index, rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			TotalRecords int                  `json:"total_records"`
			Sales        []models.SalesPerDay `json:"sales"`
		}{
			TotalRecords: totalRecords,
			Sales:        sales,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
