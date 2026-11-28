package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/services"
)

// GetPrinters returns the printers available on the host running the backend.
// The settings screen uses it to populate the USB printer selection.
func GetPrinters(logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		printers, err := services.ListAvailablePrinters()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(JSONApiOkResponse{Data: printers})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "failed to marshal printers response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(response)
	}
}
