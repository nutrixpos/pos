package dto

import (
	"github.com/elmawardy/nutrix/backend/modules/core/models"
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

type OrderStartRequestIngredient struct {
	ComponentId string  `json:"component_id" bson:"component_id"`
	EntryId     string  `json:"entry_id" bson:"entry_id"`
	Name        string  `json:"name"`
	Quantity    float32 `json:"quantity"`
	Company     string  `json:"company"`
}

type Order struct {
	models.Order `bson:",inline"`
	Ingredients  [][]OrderStartRequestIngredient `json:"ingredients"`
}

type RequestComponentEntryAdd struct {
	ComponentId string                 `json:"component_id"`
	Entries     []models.MaterialEntry `json:"entries"`
}

type FinishOrderRequest struct {
	Id string `json:"order_id"`
}

type GetComponentConsumeLogsRequest struct {
	Name string `json:"name"`
}
