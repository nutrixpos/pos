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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/services"
)

func ExportSalesCSV(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

		sales, _, err := salesService.GetSalesPerday(page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make([][]string, 0)
		data = append(data, []string{"Id", "Display Id", "Date", "Cost", "Sale Price", "Payment Source", "Refunds Value", "Profit"})

		for _, sale_day := range sales {
			for _, order := range sale_day.Orders {
				submitted_at_str := order.Order.SubmittedAt.Format(time.RFC3339)
				data = append(data, []string{order.Id, order.Order.DisplayId, submitted_at_str, fmt.Sprintf("%v", order.Order.Cost), fmt.Sprintf("%v", order.Order.SalePrice), order.Order.PaymentSource, fmt.Sprintf("%v", sale_day.RefundsValue), fmt.Sprintf("%v", order.Order.SalePrice-order.Order.Cost)})
			}
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=sales.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// 4. Write data
		for _, record := range data {
			if err := writer.Write(record); err != nil {
				http.Error(w, "Error writing csv", http.StatusInternalServerError)
				return
			}
		}
	}
}

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
