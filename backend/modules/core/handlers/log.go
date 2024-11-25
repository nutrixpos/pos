package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/modules/core/services"
)

func GetMaterialLogs(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		logService := services.Log{
			Logger: logger,
			Config: config,
		}

		logs, err := logService.GetComponentLogs(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonLogs)

	}
}

func GetSalesLog(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logService := services.Log{
			Logger: logger,
			Config: config,
		}

		sales_logs := logService.GetSalesLogs()

		jsonLogs, err := json.Marshal(sales_logs)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonLogs)

	}
}
