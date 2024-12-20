// Package dto contains Data Transfer Objects (DTOs) which are used to transfer data between application components.
// It is usually used for client-server communication.
package dto

import (
	"github.com/elmawardy/nutrix/modules/core/models"
)

// HttpComponent is a DTO containing the most important information about Material.
// It is used to return data from the API to the client.
type HttpComponent struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float32 `json:"quantity"`
	Company  string  `json:"company"`
}

// OrderItem is a DTO containing the most important information about OrderItem.
// It is used to return data from the API to the client.
type OrderItem struct {
	Name       string          `json:"name"`
	Components []HttpComponent `json:"components"`
}

// OrderStartRequestIngredient is a DTO used in the request body in the POST /order/start endpoint.
// It contains the information about the component, entry and quantity.
type OrderStartRequestIngredient struct {
	ComponentId string  `json:"component_id" bson:"component_id"`
	EntryId     string  `json:"entry_id" bson:"entry_id"`
	Name        string  `json:"name"`
	Quantity    float32 `json:"quantity"`
	Company     string  `json:"company"`
}

// Order is a DTO used to return data from the API to the client.
// It contains the ingredients and the order information.
type Order struct {
	models.Order `bson:",inline"`
	Ingredients  [][]OrderStartRequestIngredient `json:"ingredients"`
}

// GetComponentConsumeLogsRequest is a DTO used in the request body in the GET /material/consume_logs endpoint.
// It contains the material name.
type GetComponentConsumeLogsRequest struct {
	Name string `json:"name"`
}
