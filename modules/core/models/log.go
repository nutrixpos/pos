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
)

type Log struct {
	Id     string    `json:"id" bson:"id"`
	Type   string    `json:"type" bson:"type"` // the log action type
	Date   time.Time `json:"date" bson:"date"`
	UserId string    `json:"user_id" bson:"user_id"`
}

type RefundOrderLog struct {
	Log     `json:",inline" bson:",inline"`
	Reason  string `json:"reason" bson:"reason"`
	OrderId string `json:"order_id" bson:"order_id"`
}

type OrderStartLog struct {
	Log          `json:",inline" bson:",inline"`
	OrderDetails Order `json:"order_details" bson:"order_details"`
}

type OrderFinishLog struct {
	Log          `json:",inline" bson:",inline"`
	Cost         float64       `json:"cost" bson:"cost"`
	SalePrice    float64       `json:"sale_price" bson:"sale_price"`
	Items        []ItemCost    `json:"items" bson:"items"`
	OrderId      string        `json:"order_id" bson:"order_id"`
	TimeConsumed time.Duration `json:"time_consumed" bson:"time_consumed"`
}

type OrderItemRefundLog struct {
	Log     `json:",inline" bson:",inline"`
	OrderId string `json:"order_id" bson:"order_id"`
	ItemId  string `json:"item_id" bson:"item_id"`
	Reason  string `json:"reason" bson:"reason"`
}

type MaterialInventoryReturnLog struct {
	Log      `json:",inline" bson:",inline"`
	OrderId  string  `json:"order_id" bson:"order_id"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Reason   string  `json:"reason" bson:"reason"`
}

type ProductIncreaseLog struct {
	Log       `json:",inline" bson:",inline"`
	ProductId string  `json:"product_id" bson:"product_id"`
	Quantity  float64 `json:"quantity" bson:"quantity"`
	Source    string  `json:"source" bson:"source"`
	OrderId   string  `json:"order_id" bson:"order_id"`
}

type WasteOrderItemLog struct {
	Log      `json:",inline" bson:",inline"`
	Item     OrderItem `json:"item" bson:"item"`
	OrderId  string    `json:"order_id" bson:"order_id"`
	Quantity float64   `json:"quantity" bson:"quantity"`
	Reason   string    `json:"reason" bson:"reason"`
}

type DisposalMaterialAddLog struct {
	Log      `json:",inline" bson:",inline"`
	Disposal MaterialDisposal `json:"disposal"`
}

type DisposalProductAddLog struct {
	Log      `json:",inline" bson:",inline"`
	Disposal ProductDisposal `json:"disposal"`
}

type WasteMaterialLog struct {
	Log        `json:",inline" bson:",inline"`
	MaterialId string  `json:"material_id" bson:"material_id"`
	EntryId    string  `json:"entry_id" bson:"entry_id"`
	OrderId    string  `json:"order_id" bson:"order_id"`
	Reason     string  `json:"reason" bson:"reason"`
	IsConsume  bool    `json:"is_consume" bson:"is_consume"`
	Quantity   float64 `json:"quantity" bson:"quantity"`
}

type MaterialConsumeLog struct {
	Log            `json:",inline" bson:",inline"`
	MaterialId     string  `json:"material_id" bson:"material_id"`
	EntryId        string  `json:"entry_id" bson:"entry_id"`
	OrderId        string  `json:"order_id" bson:"order_id"`
	ProductId      string  `json:"recipe_id" bson:"recipe_id"`
	OrderItemIndex int     `json:"order_item_index" bson:"order_item_index"`
	Reason         string  `json:"reason" bson:"reason"`
	Quantity       float64 `json:"quantity" bson:"quantity"`
}

type MaterialAddLog struct {
	Log        `json:",inline" bson:",inline"`
	MaterialId string  `json:"material_id" bson:"material_id"`
	EntryId    string  `json:"entry_id" bson:"entry_id"`
	Quantity   float64 `json:"quantity" bson:"quantity"`
}

type ProductWasteLog struct {
	Log         `json:",inline" bson:",inline"`
	ProductId   string    `json:"product_id" bson:"product_id"`
	OrderId     string    `json:"order_id" bson:"order_id"`
	OrderItemId string    `json:"order_item_id" bson:"order_item_id"`
	Reason      string    `json:"reason" bson:"reason"`
	Quantity    float64   `json:"quantity" bson:"quantity"`
	Item        OrderItem `json:"item" bson:"item"`
}
