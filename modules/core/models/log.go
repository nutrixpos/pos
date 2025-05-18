package models

import (
	"time"
)

const (
	LogTypeDisposalAdd             = "disposal_add"
	LogTypeMaterialInventoryReturn = "material_inventory_return"
	LogTypeOrderItemRefunded       = "order_item_refunded"
	LogTypeOrderStart              = "order_Start"
	LogTypeOrderFinish             = "order_finish"
	LogTypeMaterialConsume         = "component_consume"
	LogTypeMaterialAdd             = "component_add"
	LogTypeMaterialWaste           = "material_waste"
	LogTypeProductIncrease         = "product_increase"
	LogTypeSalesPerDayOrder        = "sales_per_day_order"
	LogTypeSalesPerDayRefund       = "sales_per_day_refund"
)

type Log struct {
	Id     string    `json:"id" bson:"id" mapstructure:"id"`
	Type   string    `json:"type" bson:"type" mapstructure:"type"` // the log action type
	Date   time.Time `json:"date" bson:"date" mapstructure:"date"`
	UserId string    `json:"user_id" bson:"user_id" mapstructure:"user_id"`
}

type LogRefundOrder struct {
	Log     `json:",inline" bson:",inline" mapstructure:",squash"`
	Reason  string `json:"reason" bson:"reason" mapstructure:"reason"`
	OrderId string `json:"order_id" bson:"order_id" mapstructure:"order_id"`
}

type LogOrderStart struct {
	Log          `json:",inline" bson:",inline" mapstructure:",squash"`
	OrderDetails Order `json:"order_details" bson:"order_details" mapstructure:"order_details"`
}

type LogOrderFinish struct {
	Log          `json:",inline" bson:",inline" mapstructure:",squash"`
	Cost         float64       `json:"cost" bson:"cost" mapstructure:"cost"`
	SalePrice    float64       `json:"sale_price" bson:"sale_price" mapstructure:"sale_price"`
	Items        []ItemCost    `json:"items" bson:"items" mapstructure:"items"`
	OrderId      string        `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	TimeConsumed time.Duration `json:"time_consumed" bson:"time_consumed" mapstructure:"time_consumed"`
}

type LogOrderItemRefund struct {
	Log     `json:",inline" bson:",inline" mapstructure:",squash"`
	OrderId string `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	ItemId  string `json:"item_id" bson:"item_id" mapstructure:"item_id"`
	Reason  string `json:"reason" bson:"reason" mapstructure:"reason"`
}

type LogMaterialInventoryReturn struct {
	Log      `json:",inline" bson:",inline" mapstructure:",squash"`
	OrderId  string  `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	Quantity float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Reason   string  `json:"reason" bson:"reason" mapstructure:"reason"`
}

type LogProductIncrease struct {
	Log       `json:",inline" bson:",inline" mapstructure:",squash"`
	ProductId string  `json:"product_id" bson:"product_id" mapstructure:"product_id"`
	Quantity  float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Source    string  `json:"source" bson:"source" mapstructure:"source"`
	OrderId   string  `json:"order_id" bson:"order_id" mapstructure:"order_id"`
}

type LogWasteOrderItem struct {
	Log      `json:",inline" bson:",inline" mapstructure:",squash"`
	Item     OrderItem `json:"item" bson:"item" mapstructure:"item"`
	OrderId  string    `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	Quantity float64   `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Reason   string    `json:"reason" bson:"reason" mapstructure:"reason"`
}

type LogDisposalMaterialAdd struct {
	Log      `json:",inline" bson:",inline" mapstructure:",squash"`
	Disposal MaterialDisposal `json:"disposal" mapstructure:"disposal"`
}

type LogDisposalProductAdd struct {
	Log      `json:",inline" bson:",inline" mapstructure:",squash"`
	Disposal ProductDisposal `json:"disposal" mapstructure:"disposal"`
}

type LogWasteMaterial struct {
	Log        `json:",inline" bson:",inline" mapstructure:",squash"`
	MaterialId string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	EntryId    string  `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
	OrderId    string  `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	Reason     string  `json:"reason" bson:"reason" mapstructure:"reason"`
	IsConsume  bool    `json:"is_consume" bson:"is_consume" mapstructure:"is_consume"`
	Quantity   float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
}

type LogMaterialConsume struct {
	Log            `json:",inline" bson:",inline" mapstructure:",squash"`
	MaterialId     string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	EntryId        string  `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
	OrderId        string  `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	ProductId      string  `json:"recipe_id" bson:"recipe_id" mapstructure:"recipe_id"`
	OrderItemIndex int     `json:"order_item_index" bson:"order_item_index" mapstructure:"order_item_index"`
	Reason         string  `json:"reason" bson:"reason" mapstructure:"reason"`
	Quantity       float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
}

type LogMaterialAdd struct {
	Log        `json:",inline" bson:",inline" mapstructure:",squash"`
	MaterialId string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	EntryId    string  `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
	Quantity   float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
}

type LogProductWaste struct {
	Log         `json:",inline" bson:",inline" mapstructure:",squash"`
	ProductId   string    `json:"product_id" bson:"product_id" mapstructure:"product_id"`
	OrderId     string    `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	OrderItemId string    `json:"order_item_id" bson:"order_item_id" mapstructure:"order_item_id"`
	Reason      string    `json:"reason" bson:"reason" mapstructure:"reason"`
	Quantity    float64   `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Item        OrderItem `json:"item" bson:"item" mapstructure:"item"`
}

type LogSalesPerDayOrder struct {
	Log              `json:",inline" bson:",inline" mapstructure:",squash"`
	SalesPerDayOrder SalesPerDayOrder `json:"sales_per_day_order" bson:"sales_per_day_order" mapstructure:"sales_per_day_order"`
}

type LogSalesPerDayRefund struct {
	Log               `json:",inline" bson:",inline" mapstructure:",squash"`
	SalesPerDayRefund ItemRefund `json:"sales_per_day_refund" bson:"sales_per_day_refund" mapstructure:"sales_per_day_refund"`
}
