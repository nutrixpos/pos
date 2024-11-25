package handlers

import (
	"net/http"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/modules/core/services"
)

func HandleNotificationsWsRequest(config config.Config, logger logger.ILogger, notificationService services.INotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := notificationService.HandleHttpRequest(w, r)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
