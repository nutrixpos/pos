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

type ComponentEntrySelection struct {
	Name        string         `json:"Name"`
	ComponentId string         `json:"ComponentId"`
	Entry       ComponentEntry `json:"Entry"`
	Unit        string         `json:"Unit"`
	Quantity    float64        `json:"Quantity"`
}

type ItemCost struct {
	RecipeId   string
	ItemName   string
	Cost       float64
	SalePrice  float64
	Components []struct {
		ComponentName string
		ComponentId   string
		EntryId       string
		Quantity      float64
		Cost          float64
	}
	DownstreamCost []ItemCost
}

type RecipeSelections struct {
	RecipeName         string                    `json:"recipe_name"`
	Id                 string                    `json:"Id"`
	Ready              float64                   `json:"Ready"`
	Components         []Component               `json:"Components" bson:"-"`
	IsConsumeFromReady bool                      `json:"isConsumeFromReady"`
	Selections         []ComponentEntrySelection `json:"Selections"`
	SubRecipes         []RecipeSelections        `json:"SubRecipes"`
	Quantity           float64                   `json:"Quantity"`
}

type Order2 struct {
	SubmittedAt time.Time          `json:"submitted_at" bson:"submitted_at"`
	Id          string             `json:"id" bson:"_id,omitempty"`
	Items       []RecipeSelections `json:"items" bson:"items"`
	State       string             `json:"state" bson:"state"`
	StartedAt   time.Time          `json:"started_at" bson:"started_at"`
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
	Id              string            `bson:"_id,omitempty"`
	Name            string            `bson:"name"`
	Components      []RecipeComponent `bson:"components"`
	Price           float64           `bson:"price"`
	ImageURL        string            `bson:"image_url"`
	Ready           float64           `bson:"ready"`
	MeasuringUnit   string            `bson:"unit"`
	QuantityPerUnit float64           `bson:"quantity_per_unit"`
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
