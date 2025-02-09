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

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/services"
)

// GetSalesPerDay returns a HTTP handler function to retrieve sales data per day.
// It requires two query string parameters:
// first_index: the index of the first record to be retrieved
// rows: the number of records to be retrieved
func GetSalesPerDay(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil || page_number == 0 {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		salesService := services.SalesService{
			Logger: logger,
			Config: config,
		}

		sales, totalRecords, err := salesService.GetSalesPerday(page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Meta: JSONAPIMeta{
				TotalRecords: totalRecords,
			},
			Data: sales,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
