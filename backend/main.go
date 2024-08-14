package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules"
	"github.com/elmawardy/nutrix/modules/core"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type GetComponentConsumeLogsRequest struct {
	Name string `json:"name"`
}

type DBComponent struct {
	Id      string           `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string           `json:"name"`
	Unit    string           `json:"unit"`
	Entries []ComponentEntry `json:"entries"`
}

type HttpComponent struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float32 `json:"quantity"`
	Company  string  `json:"company"`
}

type OrderItem struct {
	Name       string          `json:"name"`
	Components []HttpComponent `json:"components"`
}

type Orderitem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Recipe  Recipe `json:"recipe"`
}

type Order struct {
	SubmittedAt time.Time                       `json:"submitted_at" bson:"submitted_at"`
	Id          string                          `json:"id" bson:"_id,omitempty"`
	Items       []Orderitem                     `json:"items"`
	State       string                          `json:"state"`
	Ingredients [][]OrderStartRequestIngredient `json:"ingredients"`
	StartedAt   time.Time                       `json:"started_at" bson:"started_at"`
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

type Category struct {
	Name    string   `json:"name"`
	Recipes []string `json:"recipes"`
}

type CategoriesContentRequest_Category struct {
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes"`
}

type PrepareItemResponse struct {
	ComponentId     string           `json:"component_id"`
	Name            string           `json:"name"`
	DefaultQuantity float32          `json:"defaultquantity"`
	Unit            string           `json:"unit"`
	Entries         []ComponentEntry `json:"entries"`
}

type OrderStartRequestIngredient struct {
	ComponentId string  `json:"component_id" bson:"component_id"`
	EntryId     string  `json:"entry_id" bson:"entry_id"`
	Name        string  `json:"name"`
	Quantity    float32 `json:"quantity"`
	Company     string  `json:"company"`
}

type OrderStartRequest struct {
	Id          string                          `json:"order_id"`
	Name        string                          `json:"name"`
	Ingredients [][]OrderStartRequestIngredient `json:"ingredients"`
}

type RequestComponentEntryAdd struct {
	ComponentId string           `json:"component_id"`
	Entries     []ComponentEntry `json:"entries"`
}

type FinishOrderRequest struct {
	Id string `json:"order_id"`
}

func main() {
	logger := logger.NewZeroLog()
	config := config.ConfigFactory("viper", "config.yaml", &logger)

	router := mux.NewRouter()

	modules_manager := modules.ModulesManager{}
	modules_manager.RegisterModule("core", &logger, core.NewBuilder(config), router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
