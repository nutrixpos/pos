package dto

import (
	"github.com/elmawardy/nutrix/modules/core/models"
)

type OrderStartRequest2 struct {
	Id    string                    `json:"order_id"`
	Items []models.RecipeSelections `json:"items"`
}
