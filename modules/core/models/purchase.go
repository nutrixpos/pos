package models

import (
	"time"
)

const (
	PurchaseOrderStatusOpen     = "open"
	PurchaseOrderStatusPartial  = "partially_received"
	PurchaseOrderStatusReceived = "received"
	PurchaseOrderStatusCancelled = "cancelled"
)

// PurchaseOrderItem represents a single line item in a purchase order.
// Quantity is the ordered quantity while ReceivedQuantity tracks how much
// has been actually received through one or more GRNs.
type PurchaseOrderItem struct {
	ItemId           string    `json:"item_id" bson:"item_id" mapstructure:"item_id"`
	MaterialId       string    `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName     string    `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	Unit             string    `json:"unit" bson:"unit" mapstructure:"unit"`
	Company          string    `json:"company" bson:"company" mapstructure:"company"`
	SKU              string    `json:"sku" bson:"sku" mapstructure:"sku"`
	Quantity         float64   `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	ReceivedQuantity float64   `json:"received_quantity" bson:"received_quantity" mapstructure:"received_quantity"`
	PurchasePrice    float64   `json:"purchase_price" bson:"purchase_price" mapstructure:"purchase_price"`
	ExpirationDate   time.Time `json:"expiration_date" bson:"expiration_date" mapstructure:"expiration_date"`
	EntryIds         []string  `json:"entry_ids" bson:"entry_ids" mapstructure:"entry_ids"`
}

// PurchaseOrder represents an organized order placed with a supplier.
type PurchaseOrder struct {
	Id         string              `json:"id" bson:"id" mapstructure:"id"`
	DisplayId  string              `json:"display_id" bson:"display_id" mapstructure:"display_id"`
	Supplier   string              `json:"supplier" bson:"supplier" mapstructure:"supplier"`
	Status     string              `json:"status" bson:"status" mapstructure:"status"`
	Notes      string              `json:"notes" bson:"notes" mapstructure:"notes"`
	Items      []PurchaseOrderItem `json:"items" bson:"items" mapstructure:"items"`
	Total      float64             `json:"total" bson:"total" mapstructure:"total"`
	AutoReceive bool               `json:"auto_receive" bson:"auto_receive" mapstructure:"auto_receive"`
	CreatedAt  time.Time           `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	CreatedBy  string              `json:"created_by" bson:"created_by" mapstructure:"created_by"`
	ReceivedAt *time.Time          `json:"received_at,omitempty" bson:"received_at,omitempty" mapstructure:"received_at,omitempty"`
	ReceivedBy string              `json:"received_by,omitempty" bson:"received_by,omitempty" mapstructure:"received_by,omitempty"`
}

// GRNItem represents a single received line item on a Goods Received Note.
// EntryId references the material entry that was created on receipt.
type GRNItem struct {
	ItemId         string    `json:"item_id" bson:"item_id" mapstructure:"item_id"`
	MaterialId     string    `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName   string    `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	Unit           string    `json:"unit" bson:"unit" mapstructure:"unit"`
	Company        string    `json:"company" bson:"company" mapstructure:"company"`
	SKU            string    `json:"sku" bson:"sku" mapstructure:"sku"`
	EntryId        string    `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
	Quantity       float64   `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	PurchasePrice  float64   `json:"purchase_price" bson:"purchase_price" mapstructure:"purchase_price"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date" mapstructure:"expiration_date"`
}

// GRN (Goods Received Note) is an automatically generated document created
// when goods are received against a purchase order. Receiving a PO pushes the
// received quantities into the materials inventory and records them in the
// materials history.
type GRN struct {
	Id                      string    `json:"id" bson:"id" mapstructure:"id"`
	DisplayId               string    `json:"display_id" bson:"display_id" mapstructure:"display_id"`
	PurchaseOrderId         string    `json:"purchase_order_id" bson:"purchase_order_id" mapstructure:"purchase_order_id"`
	PurchaseOrderDisplayId  string    `json:"purchase_order_display_id" bson:"purchase_order_display_id" mapstructure:"purchase_order_display_id"`
	Supplier                string    `json:"supplier" bson:"supplier" mapstructure:"supplier"`
	Items                   []GRNItem `json:"items" bson:"items" mapstructure:"items"`
	ReceivedAt              time.Time `json:"received_at" bson:"received_at" mapstructure:"received_at"`
	ReceivedBy              string    `json:"received_by" bson:"received_by" mapstructure:"received_by"`
}