package models

import (
	"encoding/json"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JSONFloat float64

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

type ComponentConsumeLogs struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date           time.Time          `json:"date" bson:"date"`
	Name           string             `json:"component_name" bson:"name"`
	Quantity       float32            `json:"quantity" bson:"quantity"`
	Company        string             `json:"company" bson:"company"`
	ItemName       string             `json:"item_name" bson:"item_name"`
	ItemOrderIndex uint               `json:"item_order_index" bson:"item_order_index"`
	OrderId        string             `json:"order_id" bson:"order_id"`
}

type Category struct {
	Name    string   `json:"name"`
	Recipes []string `json:"recipes"`
}

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

type OrderItemMaterial struct {
	Material Material      `json:"material"`
	Entry    MaterialEntry `json:"entry"`
	Quantity float64       `json:"quantity" bson:"quantity"`
}

type OrderItem struct {
	Id                 string              `json:"id" bson:"id"`
	Product            Product             `json:"product"`
	Materials          []OrderItemMaterial `json:"materials" bson:"materials"`
	IsConsumeFromReady bool                `json:"is_consume_from_ready"`
	SubItems           []OrderItem         `json:"sub_items" bson:"sub_items"`
	Quantity           float64             `json:"quantity" bson:"quantity"`
	Comment            string              `json:"comment" bson:"comment"`
}

type Order struct {
	SubmittedAt time.Time   `json:"submitted_at" bson:"submitted_at"`
	Id          string      `json:"id" bson:"_id,omitempty"`
	Items       []OrderItem `json:"items" bson:"items"`
	State       string      `json:"state" bson:"state"`
	StartedAt   time.Time   `json:"started_at" bson:"started_at"`
	Comment     string      `json:"comment" bson:"comment"`
}

type MaterialEntry struct {
	Id               string  `json:"_id,omitempty" bson:"_id,omitempty"`
	PurchaseQuantity float32 `json:"purchase_quantity" bson:"purchase_quantity"`
	PurchasePrice    float64 `json:"purchase_price" bson:"price"`
	Quantity         float32 `json:"quantity"`
	Company          string  `json:"company"`
	SKU              string  `json:"sku"`
}

type Material struct {
	Id       string          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string          `json:"name"`
	Entries  []MaterialEntry `json:"entries" bson:"entries"`
	Quantity float64         `json:"quantity"`
	Unit     string          `json:"unit" bson:"unit"`
}

type ProductEntry struct {
	Id               string  `json:"_id,omitempty" bson:"_id,omitempty"`
	PurchaseQuantity float32 `json:"purchase_quantity" bson:"purchase_quantity"`
	PurchasePrice    float64 `json:"purchase_price"`
	Quantity         float32 `json:"quantity"`
	Company          string  `json:"company"`
	Unit             string  `json:"unit"`
	SKU              string  `json:"sku"`
}

type Product struct {
	Id          string         `bson:"_id,omitempty" json:"id"`
	Name        string         `bson:"name" json:"name"`
	Materials   []Material     `bson:"materials" json:"materials"`
	SubProducts []Product      `bson:"sub_products" json:"sub_products"`
	Entries     []ProductEntry `bson:"entries" json:"entries"`
	Price       float64        `bson:"price" json:"price"`
	ImageURL    string         `bson:"imageurl" json:"image_url"`
	Unit        string         `bson:"unit" json:"unit"`
	Quantity    float64        `bson:"quantity" json:"quantity"`
	Ready       float64        `bson:"ready" json:"ready"`
}

type SalesLogs struct {
	Id           string    `json:"_id" bson:"_id,omitempty"`
	SalePrice    JSONFloat `json:"sale_price" bson:"sale_price"`
	Items        []ItemCost
	OrderId      string    `json:"order_id"`
	TimeConsumed time.Time `json:"time_consumed"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	Cost         JSONFloat `json:"cost"`
}
