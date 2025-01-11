// Package models contains the data models for the application.
//
// The models are used to store data in the database and to marshal/unmarshal
// data to/from JSON.
//
// The models are divided into the following categories:
//
// - Sales: Models related to sales, such as SalesPerDay and SalesPerDayOrder.
//
// - Items: Models related to items, such as ItemCost and OrderItem.
//
// - Orders: Models related to orders, such as Order and OrderItem.
//
// - Materials: Models related to materials, such as Material and MaterialEntry.
//
// - Products: Models related to products, such as Product and ProductEntry.
//
// - SalesLogs: Models related to sales logs, such as SalesLogs.
package models

import (
	"encoding/json"
	"math"
	"time"
)

// JSONFloat is a float64 that is marshaled to JSON as a string
// to avoid the JSON number type limitations.
type JSONFloat float64

// MarshalJSON implements the json.Marshaler interface.
func (j JSONFloat) MarshalJSON() ([]byte, error) {
	v := float64(j)
	if math.IsInf(v, 0) {
		// handle infinity, assign desired value to v
		// or say +/- indicates infinity
		s := -1
		if math.IsInf(v, -1) {
			s = -1
		}
		return json.Marshal(s)
	}
	return json.Marshal(v) // marshal result as standard float64
}

// SalesPerDayOrder represents an order and its associated costs for a specific day.
type SalesPerDayOrder struct {
	Order Order      `json:"order" bson:"order,inline"`
	Costs []ItemCost `json:"costs" bson:"costs"`
}

// SalesPerDay aggregates sales data for a specific day, including total costs and sales.
type SalesPerDay struct {
	Id         string             `json:"id" bson:"id,omitempty"`
	Date       string             `json:"date" bson:"date"`
	Orders     []SalesPerDayOrder `json:"orders" bson:"orders"`
	Costs      float64            `json:"costs" bson:"costs"`
	TotalSales float64            `json:"total_sales" bson:"total_sales"`
}

type ComponentConsumeLogs struct {
	Id             string    `json:"id,omitempty" bson:"id,omitempty"`
	Date           time.Time `json:"date" bson:"date"`
	Name           string    `json:"component_name" bson:"name"`
	Quantity       float32   `json:"quantity" bson:"quantity"`
	Company        string    `json:"company" bson:"company"`
	ItemName       string    `json:"item_name" bson:"item_name"`
	OrderItemIndex uint      `json:"order_item_index" bson:"order_item_index"`
	OrderId        string    `json:"order_id" bson:"order_id"`
}

// CategoryProduct represents the product in a category.
type CategoryProduct struct {
	Id       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	ImageURL string `json:"image_url" bson:"image_url"`
}

// Category represents the category of products.
type Category struct {
	Id       string            `json:"id" bson:"id"`
	Name     string            `json:"name"`
	Products []CategoryProduct `json:"products"` // product ids
}

// ItemCost represents the cost of an item, including the recipe cost, sale price, quantity,
// and the costs of the components.
type ItemCost struct {
	RecipeId   string
	ItemName   string
	Cost       float64
	SalePrice  float64
	Quantity   float64
	Components []struct {
		ComponentName string
		ComponentId   string
		EntryId       string
		Quantity      float64
		Cost          float64
	}
	DownstreamCost []ItemCost
}

// OrderItemMaterial represents the material, entry, and quantity associated with an order item.
type OrderItemMaterial struct {
	Material Material      `json:"material"`
	Entry    MaterialEntry `json:"entry"`
	Quantity float64       `json:"quantity" bson:"quantity"`
}

// OrderItem represents an item in an order, including product details, materials, and pricing.
type OrderItem struct {
	Id                 string              `json:"id" bson:"id"`
	Product            Product             `json:"product"`
	Price              float64             `json:"price" bson:"price"`
	Materials          []OrderItemMaterial `json:"materials" bson:"materials"`
	IsConsumeFromReady bool                `json:"is_consume_from_ready"`
	SubItems           []OrderItem         `json:"sub_items" bson:"sub_items"`
	Quantity           float64             `json:"quantity" bson:"quantity"`
	Comment            string              `json:"comment" bson:"comment"`
	SalePrice          float64             `json:"sale_price" bson:"sale_price"`
	Cost               float64             `json:"cost" bson:"cost"`
}

// Order represents a customer order, containing order details, items, and financial information.
type Order struct {
	SubmittedAt time.Time   `json:"submitted_at" bson:"submitted_at"`
	Id          string      `json:"id" bson:"id,omitempty"`
	DisplayId   string      `json:"display_id" bson:"display_id"`
	Items       []OrderItem `json:"items" bson:"items"`
	Discount    float64     `json:"discount" bson:"discount"`
	State       string      `json:"state" bson:"state"`
	StartedAt   time.Time   `json:"started_at" bson:"started_at"`
	Comment     string      `json:"comment" bson:"comment"`
	Cost        float64     `json:"cost" bson:"cost"`
	SalePrice   float64     `json:"sale_price" bson:"sale_price"`
	IsPayLater  bool        `json:"is_pay_later" bson:"is_pay_later"`
	IsPaid      bool        `json:"is_paid" bson:"is_paid"`
}

// MaterialEntry represents an entry of material, detailing purchase and quantity information.
type MaterialEntry struct {
	Id               string    `json:"id,omitempty" bson:"id,omitempty"`
	PurchaseQuantity float32   `json:"purchase_quantity" bson:"purchase_quantity"`
	PurchasePrice    float64   `json:"purchase_price" bson:"price"`
	Quantity         float32   `json:"quantity"`
	Company          string    `json:"company"`
	SKU              string    `json:"sku"`
	ExpirationDate   time.Time `json:"expiration_date" bson:"expiration_date"`
}

// MaterialSettings represents settings associated with a material, such as stock alert threshold.
type MaterialSettings struct {
	StockAlertTreshold float64 `json:"stock_alert_treshold" bson:"stock_alert_treshold"`
}

// Material represents a material with its details, including entries and settings.
type Material struct {
	Id       string           `json:"id,omitempty" bson:"id,omitempty"`
	Name     string           `json:"name"`
	Entries  []MaterialEntry  `json:"entries" bson:"entries"`
	Quantity float64          `json:"quantity"`
	Settings MaterialSettings `json:"settings" bson:"settings"`
	Unit     string           `json:"unit" bson:"unit"`
}

// ProductEntry represents an entry of a product, detailing purchase and quantity information.
type ProductEntry struct {
	Id               string  `json:"id,omitempty" bson:"id,omitempty"`
	PurchaseQuantity float32 `json:"purchase_quantity" bson:"purchase_quantity"`
	PurchasePrice    float64 `json:"purchase_price"`
	Quantity         float32 `json:"quantity"`
	Company          string  `json:"company"`
	Unit             string  `json:"unit"`
	SKU              string  `json:"sku"`
}

// Product represents a product with its details, including materials, entries, and pricing.
type Product struct {
	Id          string         `bson:"id,omitempty" json:"id"`
	Name        string         `bson:"name" json:"name"`
	Materials   []Material     `bson:"materials" json:"materials"`
	SubProducts []Product      `bson:"sub_products" json:"sub_products"`
	Entries     []ProductEntry `bson:"entries" json:"entries"`
	Price       float64        `bson:"price" json:"price"`
	ImageURL    string         `bson:"image_url" json:"image_url"`
	Unit        string         `bson:"unit" json:"unit"`
	Quantity    float64        `bson:"quantity" json:"quantity"`
	Ready       float64        `bson:"ready" json:"ready"`
}

// SalesLogs represents logs of sales, capturing sale price, items, and consumption details.
type SalesLogs struct {
	Id           string    `json:"id" bson:"id,omitempty"`
	SalePrice    JSONFloat `json:"sale_price" bson:"sale_price"`
	Items        []ItemCost
	OrderId      string    `json:"order_id"`
	TimeConsumed time.Time `json:"time_consumed"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	Cost         JSONFloat `json:"cost"`
}
