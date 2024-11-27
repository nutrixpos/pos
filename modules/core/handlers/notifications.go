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
	"net/http"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/modules/core/services"
)

// HandleNotificationsWsRequest returns a HTTP handler function to handle WebSocket requests.
//
// The function takes a configuration object, a logger object, and a INotificationService object as input.
// It returns a HTTP handler function that handles WebSocket requests.
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
