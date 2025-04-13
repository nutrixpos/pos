package models

import (
	"time"
)

const (
	DisposalAdd = "disposal_add"
)

type Log struct {
	Id     string    `json:"id" bson:"id"`
	Type   string    `json:"type" bson:"type"` // the log action type
	Date   time.Time `json:"date" bson:"date"`
	UserId string    `json:"user_id" bson:"user_id"`
}

type RefundOrderLog struct {
	Log     `json:",inline"`
	Reason  string `json:"reason" bson:"reason"`
	OrderId string `json:"order_id" bson:"order_id"`
}

type OrderItemRefundLog struct {
	Log     `json:",inline"`
	OrderId string `json:"order_id" bson:"order_id"`
	ItemId  string `json:"item_id" bson:"item_id"`
	Reason  string `json:"reason" bson:"reason"`
}

type MaterialInventoryReturnLog struct {
	Log      `json:",inline"`
	OrderId  string  `json:"order_id" bson:"order_id"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Reason   string  `json:"reason" bson:"reason"`
}

type ProductIncreaseLog struct {
	Log       `json:",inline"`
	ProductId string                 `json:"product_id" bson:"product_id"`
	Quantity  float64                `json:"quantity" bson:"quantity"`
	Source    string                 `json:"source" bson:"source"`
	Other     map[string]interface{} `json:"other" bson:"other"`
}

type WasteOrderItemLog struct {
	Log      `json:",inline"`
	Item     OrderItem              `json:"item" bson:"item"`
	OrderId  string                 `json:"order_id" bson:"order_id"`
	Quantity float64                `json:"quantity" bson:"quantity"`
	Reason   string                 `json:"reason" bson:"reason"`
	Other    map[string]interface{} `json:"other" bson:"other"`
}

type DisposalAddLog struct {
	Log      `json:",inline"`
	Disposal Disposal `json:"disposal"`
}

type WasteMaterialLog struct {
	Log        `json:",inline"`
	MaterialId string  `json:"material_id" bson:"material_id"`
	EntryId    string  `json:"entry_id" bson:"entry_id"`
	OrderId    string  `json:"order_id" bson:"order_id"`
	Reason     string  `json:"reason" bson:"reason"`
	IsConsume  bool    `json:"is_consume" bson:"is_consume"`
	Quantity   float64 `json:"quantity" bson:"quantity"`
}

type MaterialConsumeLog struct {
	Log            `json:",inline"`
	MaterialId     string  `json:"material_id" bson:"material_id"`
	EntryId        string  `json:"entry_id" bson:"entry_id"`
	OrderId        string  `json:"order_id" bson:"order_id"`
	ProductId      string  `json:"recipe_id" bson:"recipe_id"`
	OrderItemIndex int     `json:"order_item_index" bson:"order_item_index"`
	Reason         string  `json:"reason" bson:"reason"`
	Quantity       float64 `json:"quantity" bson:"quantity"`
}

type ProductWasteLog struct {
	Log            `json:",inline"`
	ProductId      string    `json:"product_id" bson:"product_id"`
	OrderId        string    `json:"order_id" bson:"order_id"`
	OrderItemIndex int       `json:"order_item_index" bson:"order_item_index"`
	Reason         string    `json:"reason" bson:"reason"`
	Quantity       float64   `json:"quantity" bson:"quantity"`
	Item           OrderItem `json:"item" bson:"item"`
}
