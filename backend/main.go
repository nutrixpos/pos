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

	// deleting an entry from a component
	router.HandleFunc("/api/entry", func(w http.ResponseWriter, r *http.Request) {

		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
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

		// Retrieve the entry ID from the query string
		entryIDStr := r.URL.Query().Get("entry_id")
		entryID, err := primitive.ObjectIDFromHex(entryIDStr)
		if err != nil {
			http.Error(w, "Invalid entry ID", http.StatusBadRequest)
			return
		}

		componentIdStr := r.URL.Query().Get("component_id")
		componentId, err := primitive.ObjectIDFromHex(componentIdStr)
		if err != nil {
			http.Error(w, "Invalid entry ID", http.StatusBadRequest)
			return
		}

		// Connect to the database
		collection := client.Database("waha").Collection("components")

		// Find the component document and update the entries array
		filter := bson.M{"_id": componentId}
		update := bson.M{"$pull": bson.M{"entries": bson.M{"_id": entryID}}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
			return
		}

		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry deleted successfully"))

	}).Methods("GET")

	router.HandleFunc("/api/componententry", func(w http.ResponseWriter, r *http.Request) {

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

		var componentEntryRequest RequestComponentEntryAdd
		err = json.NewDecoder(r.Body).Decode(&componentEntryRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(componentEntryRequest.ComponentId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": objectId}

		for _, entry := range componentEntryRequest.Entries {

			entry.Id = primitive.NewObjectID()

			entry.PurchaseQuantity = entry.Quantity

			update := bson.M{"$push": bson.M{"entries": entry}}
			opts := options.Update().SetUpsert(false)

			_, err = client.Database("waha").Collection("components").UpdateOne(ctx, filter, update, opts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	})

	router.HandleFunc("/api/component", func(w http.ResponseWriter, r *http.Request) {

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

		// Parse the request body into a DBComponent struct
		var dbComponent DBComponent
		err = json.NewDecoder(r.Body).Decode(&dbComponent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for index, entry := range dbComponent.Entries {
			dbComponent.Entries[index].PurchaseQuantity = entry.Quantity
		}

		// Insert the DBComponent struct into the "components" collection
		collection := client.Database("waha").Collection("components")
		_, err = collection.InsertOne(ctx, dbComponent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, entry := range dbComponent.Entries {

			entry.Id = primitive.NewObjectID()

			logs_data := bson.M{
				"type":     "component_add",
				"date":     time.Now(),
				"company":  entry.Company,
				"quantity": entry.Quantity,
				"price":    entry.Price,
			}
			_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Return a success response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "component adding saved successfully")

	})

	router.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {

		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
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

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id quert string is required", http.StatusBadRequest)
			return
		}

		coll := client.Database("waha").Collection("orders")
		var order Order

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrder)
	})

	router.HandleFunc("/api/componentlogs", func(w http.ResponseWriter, r *http.Request) {

		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
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

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		filter := bson.M{"type": "component_consume", "name": name}
		cur, err := client.Database("waha").Collection("logs").Find(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		var logs []ComponentConsumeLogs
		if err = cur.All(ctx, &logs); err != nil {
			log.Fatal(err)
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonLogs)

	})

	router.HandleFunc("/api/components", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
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

		// Get the "test" collection from the database
		cur, err := client.Database("waha").Collection("components").Find(ctx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}

		defer cur.Close(ctx)

		// Iterate over the documents and print them as JSON
		var components []DBComponent
		for cur.Next(ctx) {
			var component DBComponent
			err := cur.Decode(&component)
			if err != nil {
				log.Fatal(err)
			}
			components = append(components, component)
		}

		// Convert the slice to JSON
		jsonComponents, err := json.Marshal(components)
		if err != nil {
			log.Fatal(err)
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonComponents)
	})

	router.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {

		// an example API handler
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, OPTIONS")
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

		// Get categories from the database, decode to the Category struct
		// and use the Recipes array of string representing the Ids of the recipes
		// then bring those recipes and embed them into a CategoriesContentRequest_Category slice

		// Define a slice to hold the categories
		var categories []CategoriesContentRequest_Category

		// Fetch categories from the database
		cur, err := client.Database("waha").Collection("categories").Find(ctx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		// Iterate through the categories
		for cur.Next(ctx) {
			var category Category
			err := cur.Decode(&category)
			if err != nil {
				log.Fatal(err)
			}

			// Fetch the recipes for each category using the Recipe IDs
			var recipes []Recipe
			for _, recipeID := range category.Recipes {
				var recipe Recipe
				obj_id, _ := primitive.ObjectIDFromHex(recipeID)
				filter := bson.D{{Key: "_id", Value: obj_id}}
				err := client.Database("waha").Collection("recipes").FindOne(ctx, filter).Decode(&recipe)
				if err != nil {
					log.Fatal(err)
				}
				recipes = append(recipes, recipe)
			}

			// Create a CategoriesContentRequest_Category with embedded recipes
			contentCategory := CategoriesContentRequest_Category{
				Name:    category.Name,
				Recipes: recipes,
			}

			// Append the contentCategory to the categories slice
			categories = append(categories, contentCategory)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(categories); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}).Methods("GET")

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
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
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

		// decrease the ingredient component quantity from the components inventory

		for itemIndex, ingredient := range order_start_request.Ingredients {
			for componentIndex, component := range ingredient {
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				filter := bson.M{"name": component.Name, "entries.company": component.Company}
				// Define the update operation
				update := bson.M{
					"$inc": bson.M{
						"entries.$.quantity": -component.Quantity,
					},
				}

				var db_component DBComponent

				for _, entry := range db_component.Entries {
					if entry.Company == component.Company {
						component.EntryId = entry.Id.Hex()
					}
				}

				err = client.Database("waha").Collection("components").FindOne(context.Background(), filter).Decode(&db_component)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				component.ComponentId = db_component.Id
				order_start_request.Ingredients[itemIndex][componentIndex] = component

				_, err = client.Database("waha").Collection("components").UpdateOne(context.Background(), filter, update)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// upsertedId := fmt.Sprintf("%s", result.UpsertedID)
				// updatedId, _ := primitive.ObjectIDFromHex(upsertedId)

				logs_data := bson.M{
					"type":             "component_consume",
					"date":             time.Now(),
					"name":             component.Name,
					"quantity":         component.Quantity,
					"company":          component.Company,
					"order_id":         order_start_request.Id,
					"item_name":        order.Items[itemIndex].Name,
					"item_order_index": itemIndex,
				}
				_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
				if err != nil {
					log.Fatal(err)
				}

			}
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

		logs_data := bson.M{
			"type":          "order_Start",
			"date":          time.Now(),
			"order_details": order,
		}
		_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			log.Fatal(err)
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
				ComponentId:     db_component.Id,
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

		cursor, err := client.Database("waha").Collection("orders").Find(context.Background(), bson.M{
			"state": bson.M{
				"$ne": "finished",
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var order Order
			if err := cursor.Decode(&order); err != nil {
				log.Fatal(err)
			}

			for i, item := range order.Items {

				var recipe Recipe

				err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"name": item.Name}).Decode(&recipe)
				if err == nil {
					order.Items[i].Recipe = recipe
				}

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
		var finish_order_request FinishOrderRequest
		err = decoder.Decode(&finish_order_request)
		if err != nil {
			panic(err)
		}

		collection := client.Database("waha").Collection("orders")
		Id, err := primitive.ObjectIDFromHex(finish_order_request.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filter := bson.M{"_id": Id}
		update := bson.M{"$set": bson.M{"state": "finished"}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			panic(err)
		}

		var order Order
		err = collection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalCost := 0.0
		totalSalePrice := 0.0

		itemsCost := []struct {
			ItemName   string
			Cost       float64
			SalePrice  float64
			Components []struct {
				ComponentName string
				Cost          float64
			}
		}{}

		for itemIndex, itemIngredients := range order.Ingredients {

			itemCost := struct {
				ItemName   string
				Cost       float64
				SalePrice  float64
				Components []struct {
					ComponentName string
					Cost          float64
				}
			}{
				ItemName: order.Items[itemIndex].Name,
				Cost:     0.0,
			}

			for _, ingredient := range itemIngredients {

				itemComponent := struct {
					ComponentName string
					Cost          float64
				}{
					ComponentName: ingredient.Name,
				}

				var component_with_specific_entry DBComponent

				component_id, _ := primitive.ObjectIDFromHex(ingredient.ComponentId)
				entry_id, _ := primitive.ObjectIDFromHex(ingredient.EntryId)

				err = client.Database("waha").Collection("components").FindOne(
					context.Background(), bson.M{"_id": component_id, "entries._id": entry_id}, options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&component_with_specific_entry)

				if err == nil {
					quantity_cost := (component_with_specific_entry.Entries[0].Price / float64(component_with_specific_entry.Entries[0].PurchaseQuantity)) * float64(ingredient.Quantity)
					totalCost += quantity_cost
					itemCost.Cost += quantity_cost
					itemComponent.Cost = quantity_cost
				}

				itemCost.Components = append(itemCost.Components, itemComponent)

			}

			var recipe Recipe

			recipeID, _ := primitive.ObjectIDFromHex(order.Items[itemIndex].Id)

			err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"_id": recipeID}).Decode(&recipe)

			if err != nil {
				panic(err)
			}

			itemCost.SalePrice = recipe.Price
			totalSalePrice += recipe.Price

			itemsCost = append(itemsCost, itemCost)
		}

		logs_data := bson.M{
			"type":          "order_finish",
			"date":          time.Now(),
			"cost":          totalCost,
			"sale_price":    totalSalePrice,
			"items":         itemsCost,
			"order_id":      finish_order_request.Id,
			"time_consumed": time.Since(order.SubmittedAt),
		}
		_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			log.Fatal(err)
		}

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
