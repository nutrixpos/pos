package dto

import (
	"github.com/elmawardy/nutrix/modules/core/models"
)

type OrderStartRequest struct {
	Id    string             `json:"order_id"`
	Items []models.OrderItem `json:"items"`
}
