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

type ComponentEntry struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PurchaseQuantity float32            `json:"purchase_quantity" bson:"purchase_quantity"`
	Quantity         float32            `json:"quantity"`
	Company          string             `json:"company"`
	Unit             string             `json:"unit"`
	Price            float64            `json:"price"`
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

type Component struct {
	Id      string           `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string           `json:"name"`
	Unit    string           `json:"unit"`
	Entries []ComponentEntry `json:"entries"`
}

type Orderitem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Recipe  Recipe `json:"recipe"`
}

type Category struct {
	Name    string   `json:"name"`
	Recipes []string `json:"recipes"`
}

type Order struct {
	SubmittedAt time.Time   `json:"submitted_at" bson:"submitted_at"`
	Id          string      `json:"id" bson:"_id,omitempty"`
	Items       []Orderitem `json:"items" bson:"items"`
	State       string      `json:"state" bson:"state"`
	StartedAt   time.Time   `json:"started_at" bson:"started_at"`
}

type RecipeComponent struct {
	ComponentId string  `bson:"component_id"`
	Name        string  `bson:"name"`
	Quantity    float32 `bson:"quantity"`
	Unit        string  `bson:"unit"`
	Type        string  `bson:"type"`
}

type Recipe struct {
	Id         string            `bson:"_id,omitempty"`
	Name       string            `bson:"name"`
	Components []RecipeComponent `bson:"components"`
	Price      float64           `bson:"price"`
	ImageURL   string            `bson:"image_url"`
}

type SalesLogs struct {
	Id        string    `json:"_id" bson:"_id,omitempty"`
	SalePrice JSONFloat `json:"sale_price" bson:"sale_price"`
	Items     []struct {
		ItemName   string    `json:"itemname"`
		Cost       JSONFloat `json:"cost"`
		SalePrice  JSONFloat `json:"sale_price"`
		Components []struct {
			ComponentName string    `json:"componentname"`
			Cost          JSONFloat `json:"cost"`
		}
	}
	OrderId      string    `json:"order_id"`
	TimeConsumed time.Time `json:"time_consumed"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	Cost         JSONFloat `json:"cost"`
}
