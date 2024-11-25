package dto

import (
	"github.com/elmawardy/nutrix/backend/modules/core/models"
)

type OrderStartRequest struct {
	Id    string             `json:"order_id"`
	Items []models.OrderItem `json:"items"`
}

type OrderStashRequest struct {
	Order models.Order `json:"order"`
}

type OrderStashResponse struct {
	Order models.Order `json:"order"`
}

type OrderRemoveStashRequest struct {
	OrderDisplayId string `json:"order_display_id"`
}
