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

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

// GetSalesPerDay returns a HTTP handler function to retrieve sales data per day.
// It requires two query string parameters:
// first_index: the index of the first record to be retrieved
// rows: the number of records to be retrieved
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
