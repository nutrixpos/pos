// Package dto contains all the data transfer objects that are used in the API.
// These struct will be used to unmarshal json data from the request body
// and return data to the client.
package dto

import (
	"github.com/elmawardy/nutrix/backend/modules/core/models"
)

// OrderStartRequest is a DTO used in the request body in the POST /order/start endpoint.
// It contains the order id and the items to start.
type OrderStartRequest struct {
	Id    string             `json:"order_id"`
	Items []models.OrderItem `json:"items"`
}

// OrderStashRequest is a DTO used in the request body in the POST /order/stash endpoint.
// It contains the order to stash.
type OrderStashRequest struct {
	Order models.Order `json:"order"`
}

// OrderStashResponse is a DTO used in the response body in the POST /order/stash endpoint.
// It contains the stashed order.
type OrderStashResponse struct {
	Order models.Order `json:"order"`
}

// OrderRemoveStashRequest is a DTO used in the request body in the POST /order/remove_stash endpoint.
// It contains the order display id to remove from stash.
type OrderRemoveStashRequest struct {
	OrderDisplayId string `json:"order_display_id"`
}
