package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elmawardy/waha-backend/globals"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ComponentEntry struct {
	Quantity float32 `json:"quantity"`
	Company  string  `json:"company"`
	Unit     string  `json:"unit"`
}

type DBComponent struct {
	Name    string           `json:"name"`
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
	Name    string `json:"name"`
	Comment string `json:"comment"`
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
	Name     string  `bson:"name"`
	Quantity float32 `bson:"quantity"`
	Unit     string  `bson:"unit"`
	Type     string  `bson:"type"`
}

type Recipe struct {
	Name       string            `bson:"name"`
	Components []RecipeComponent `bson:"components"`
}

type PrepareItemResponse struct {
	Name            string           `json:"name"`
	DefaultQuantity float32          `json:"defaultquantity"`
	Unit            string           `json:"unit"`
	Entries         []ComponentEntry `json:"entries"`
}

type OrderStartRequestIngredient struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
}

type OrderStartRequest struct {
	Id          string                          `json:"order_id"`
	Name        string                          `json:"name"`
	Ingredients [][]OrderStartRequestIngredient `json:"ingredients"`
}

func main() {

	const defaultPort = "8000"

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	DBHost := os.Getenv("SWP_DB_HOST")
	if DBHost == "" {
		DBHost = "localhost"
		// panic("SWP_DB_HOST env var not set.")
	}

	// ADMIN_EMAIL := os.Getenv("SWP_ADMIN_EMAIL")
	// if ADMIN_EMAIL == "" {
	// 	panic("SWP_ADMIN_EMAIL env var not set.")
	// }

	// SMTP_FROM := os.Getenv("SWP_SMTP_FROM")
	// if SMTP_FROM == "" {
	// 	panic("SWP_SMTP_FROM env var not set.")
	// }

	// SMTP_PASSWORD := os.Getenv("SWP_SMTP_PASSWORD")
	// if SMTP_PASSWORD == "" {
	// 	panic("SWP_SMTP_PASSWORD env var not set.")
	// }

	globals.Init(DBHost)

	router := mux.NewRouter()

	router.HandleFunc("/api/startorder", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", globals.DBHost))

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Connected successfully
		fmt.Println("Connected to MongoDB!")

		decoder := json.NewDecoder(r.Body)
		var order_start_request OrderStartRequest
		err = decoder.Decode(&order_start_request)
		if err != nil {
			panic(err)
		}

		var order Order

		order_id, err := primitive.ObjectIDFromHex(order_start_request.Id)
		if err != nil {
			panic(err)
		}

		err = client.Database("waha").Collection("orders").FindOne(context.Background(), bson.M{"_id": order_id}).Decode(&order)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"ingredients": order_start_request.Ingredients,
				"state":       "in_progress",
				"started_at":  time.Now(),
			},
		}
		_, err = client.Database("waha").Collection("orders").UpdateOne(context.Background(), bson.M{"_id": order_id}, update)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	})

	router.HandleFunc("/api/prepareitem", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", globals.DBHost))

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Connected successfully
		fmt.Println("Connected to MongoDB!")

		decoder := json.NewDecoder(r.Body)
		var orderitem OrderItem
		err = decoder.Decode(&orderitem)
		if err != nil {
			panic(err)
		}

		var recipe Recipe

		err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"name": orderitem.Name}).Decode(&recipe)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var prepare_item_responses []PrepareItemResponse

		for _, component := range recipe.Components {

			var db_component DBComponent

			if component.Type != "recipe" {
				err = client.Database("waha").Collection("components").FindOne(context.Background(), bson.M{"name": component.Name}).Decode(&db_component)
				if err != nil {
					log.Fatal(err)
				}
			}

			res := PrepareItemResponse{
				Name:            component.Name,
				DefaultQuantity: component.Quantity,
				Unit:            component.Unit,
			}

			if component.Type != "recipe" {
				res.Entries = db_component.Entries
			} else {
				res.Entries = []ComponentEntry{}
			}

			prepare_item_responses = append(prepare_item_responses, res)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(prepare_item_responses); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	})

	router.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", globals.DBHost))

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Connected successfully
		fmt.Println("Connected to MongoDB!")

		var orders []Order

		cursor, err := client.Database("waha").Collection("orders").Find(context.Background(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var order Order
			if err := cursor.Decode(&order); err != nil {
				log.Fatal(err)
			}
			orders = append(orders, order)
		}

		// Check for errors during iteration
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/api/submitorder", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", globals.DBHost))

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Connected successfully
		fmt.Println("Connected to MongoDB!")

		decoder := json.NewDecoder(r.Body)
		var order Order
		err = decoder.Decode(&order)
		if err != nil {
			panic(err)
		}

		order_data := bson.M{
			"items":        order.Items,
			"id":           order.Id,
			"submitted_at": time.Now(),
		}
		_, err = client.Database("waha").Collection("orders").InsertOne(ctx, order_data)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)

	})

	router.HandleFunc("/api/finishorder", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", globals.DBHost))

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Connected successfully
		fmt.Println("Connected to MongoDB!")

		decoder := json.NewDecoder(r.Body)
		var order OrderItem
		err = decoder.Decode(&order)
		if err != nil {
			panic(err)
		}

		for _, component := range order.Components {

			filter := bson.M{"name": component.Name, "entries.company": component.Company}
			update := bson.M{"$inc": bson.M{"entries.$.quantity": -component.Quantity}}
			_, err := client.Database("waha").Collection("components").UpdateOne(context.TODO(), filter, update)

			if err != nil {
				log.Fatal(err)
			}

			logs_data := bson.M{
				"type":       "new_order",
				"date":       time.Now(),
				"order_name": order.Name,
				"component":  component,
			}
			_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
			if err != nil {
				log.Fatal(err)
			}
		}

		// _, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		w.WriteHeader(http.StatusOK)

	}).Methods("OPTIONS", "POST")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
