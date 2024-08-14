package dto

import (
	"github.com/elmawardy/nutrix/modules/core/models"
)

type HttpComponent struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float32 `json:"quantity"`
	Company  string  `json:"company"`
}

type OrderItem struct {
	Name       string          `json:"name"`
	Components []HttpComponent `json:"components"`
}

type Category struct {
	Name    string          `json:"name"`
	Recipes []models.Recipe `json:"recipes"`
}

type PrepareItemResponse struct {
	ComponentId     string                  `json:"component_id"`
	Name            string                  `json:"name"`
	DefaultQuantity float32                 `json:"defaultquantity"`
	Unit            string                  `json:"unit"`
	Entries         []models.ComponentEntry `json:"entries"`
}

type OrderStartRequestIngredient struct {
	ComponentId string  `json:"component_id" bson:"component_id"`
	EntryId     string  `json:"entry_id" bson:"entry_id"`
	Name        string  `json:"name"`
	Quantity    float32 `json:"quantity"`
	Company     string  `json:"company"`
}

type OrderStartRequest struct {
	Id          string                          `json:"order_id"`
	Name        string                          `json:"name"`
	Ingredients [][]OrderStartRequestIngredient `json:"ingredients"`
}

type Order struct {
	models.Order `bson:",inline"`
	Ingredients  [][]OrderStartRequestIngredient `json:"ingredients"`
}

type RequestComponentEntryAdd struct {
	ComponentId string                  `json:"component_id"`
	Entries     []models.ComponentEntry `json:"entries"`
}

type FinishOrderRequest struct {
	Id string `json:"order_id"`
}

type GetComponentConsumeLogsRequest struct {
	Name string `json:"name"`
}
