// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderService is the service to interact with the orders collection in the database.
type OrderService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

func (os *OrderService) PrintReceipt(order models.Order, template string, lang_code string) (err error) {
	receipt_svc := ReceiptService{
		Config:   os.Config,
		Logger:   os.Logger,
		Settings: os.Settings,
	}

	err = receipt_svc.Print(order, order.Discount, 0, order.SubmittedAt, lang_code, template)
	if err != nil {
		return err
	}

	return nil
}

func (os *OrderService) DeleteOrder(order_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))
	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	collection := client.Database("waha").Collection("orders")
	_, err = collection.DeleteOne(ctx, bson.M{"id": order_id})

	return
}

// PayUnpaidOrder sets the is_paid field of the order with the given order_id to true.
func (os *OrderService) PayUnpaidOrder(order_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))
	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	collection := client.Database("waha").Collection("orders")

	filter := bson.M{"id": order_id}
	update := bson.M{"$set": bson.M{"is_paid": true}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return
}

// GetUnpaidOrders returns all orders that are not paid and their state is not cancelled.
func (os *OrderService) GetUnpaidOrders() (orders []models.Order, err error) {

	orders = make([]models.Order, 0)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))
	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	collection := client.Database("waha").Collection("orders")

	cursor, err := collection.Find(ctx, bson.M{"is_pay_later": true, "is_paid": false, "state": bson.M{"$not": bson.M{"$in": []string{"cancelled"}}}})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		var order models.Order
		err = cursor.Decode(&order)
		if err != nil {
			return
		}
		orders = append(orders, order)
	}

	return
}

// CancelOrder sets the state of the order with the given order_id to "cancelled".
func (os *OrderService) CancelOrder(order_id string) (err error) {
	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Connected successfully
	collection := client.Database("waha").Collection("orders")

	// Define filter and update for setting order state to "cancelled"
	filter := bson.M{"id": order_id}
	update := bson.M{"$set": bson.M{"state": "cancelled"}}

	// Update the order in the database
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

// CalculateCost calculates the cost of each item in the provided list of order items.
func (os *OrderService) CalculateCost(items []models.OrderItem) (cost []models.ItemCost, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))
	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return cost, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return cost, err
	}

	// Connected successfully
	os.Logger.Info("Connected to MongoDB!")

	for itemIndex, item := range items {

		itemCost := models.ItemCost{
			ItemName: items[itemIndex].Product.Name,
			Cost:     0.0,
			RecipeId: items[itemIndex].Product.Id,
		}

		for _, component := range item.Materials {

			itemComponent := struct {
				ComponentName string
				ComponentId   string
				EntryId       string
				Quantity      float64
				Cost          float64
			}{
				ComponentName: component.Material.Name,
				ComponentId:   component.Material.Id,
				EntryId:       component.Entry.Id,
				Quantity:      component.Quantity,
			}

			var component_with_specific_entry models.Material

			err = client.Database("waha").Collection("materials").FindOne(
				context.Background(), bson.M{"id": component.Material.Id, "entries.id": component.Entry.Id}, options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&component_with_specific_entry)

			if err == nil {
				quantity_cost := (component_with_specific_entry.Entries[0].PurchasePrice / float64(component_with_specific_entry.Entries[0].PurchaseQuantity)) * float64(component.Quantity)

				// check if cost is positive or negative infinity (semantic bug in calculation that causes problems later on)
				if math.IsInf(quantity_cost, 0) || math.IsInf(quantity_cost, -1) {
					quantity_cost = 0
				}

				itemCost.Cost += quantity_cost
				itemComponent.Cost = quantity_cost

			} else {
				return cost, err
			}

			itemCost.Components = append(itemCost.Components, itemComponent)

		}

		var recipe models.Product
		err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"id": items[itemIndex].Product.Id}).Decode(&recipe)

		if err != nil {
			panic(err)
		}

		for _, subrecipe := range item.SubItems {
			total_cost, err := os.CalculateCost([]models.OrderItem{subrecipe})
			if err != nil {
				return cost, err
			}

			for _, subrecipe_cost := range total_cost {
				itemCost.Cost += subrecipe_cost.Cost * float64(subrecipe.Quantity)
				subrecipe_cost.Quantity = subrecipe.Quantity
				itemCost.DownstreamCost = append(itemCost.DownstreamCost, subrecipe_cost)
			}
		}

		itemCost.SalePrice = recipe.Price
		cost = append(cost, itemCost)
	}

	return cost, err
}

// FinishOrder sets the state of the order with the given order_id to "finished".
func (os *OrderService) FinishOrder(order_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Connected successfully
	os.Logger.Info("Connected to MongoDB!")

	collection := client.Database("waha").Collection("orders")

	filter := bson.M{"id": order_id}
	update := bson.M{"$set": bson.M{"state": "finished"}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	var order models.Order
	err = collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return err
	}

	totalCost := 0.0
	totalSalePrice := 0.0

	items_cost, err := os.CalculateCost(order.Items)
	if err != nil {
		return err
	}

	for index, recipe_cost := range items_cost {

		order.Items[index].Cost = recipe_cost.Cost
		order.Items[index].SalePrice = recipe_cost.SalePrice

		totalCost += recipe_cost.Cost
		totalSalePrice += recipe_cost.SalePrice
	}

	logs_data := bson.M{
		"type":          "order_finish",
		"date":          time.Now(),
		"cost":          totalCost,
		"sale_price":    totalSalePrice,
		"items":         items_cost,
		"order_id":      order_id,
		"time_consumed": time.Since(order.SubmittedAt),
	}
	_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
	if err != nil {
		log.Fatal(err)
	}

	order.Cost = totalCost
	order.SalePrice = totalSalePrice

	salesSvc := SalesService{Config: os.Config, Logger: os.Logger}
	err = salesSvc.AddOrderToSalesDay(order, items_cost)
	if err != nil {
		return err
	}

	return err
}

// GetOrderDisplayId returns a new order display id and increments the current value in the database.
func (os *OrderService) GetOrderDisplayId() (order_display_id string, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return order_display_id, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return order_display_id, err
	}

	var settings models.Settings
	err = client.Database("waha").Collection("settings").FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		return order_display_id, err
	}

	random_queue_index := 0

	if len(settings.Orders.Queues) > 1 {
		random_queue_index = rand.Intn((len(settings.Orders.Queues)-1)-0) + 0
	}
	random_queue := settings.Orders.Queues[random_queue_index]

	order_display_id = fmt.Sprintf("%s-%v", random_queue.Prefix, random_queue.Next)

	_, err = client.Database("waha").Collection("settings").UpdateOne(
		ctx,
		bson.M{"id": settings.Id, "orders.queues.prefix": random_queue.Prefix},
		bson.M{
			"$inc": bson.M{
				"orders.queues.$.next": 1,
			},
		},
	)

	if err != nil {
		return order_display_id, err
	}

	return order_display_id, err

}

// SubmitOrder adds an order to the database and creates a display id.
func (os *OrderService) SubmitOrder(order models.Order) (models.Order, error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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
	os.Logger.Info("Connected to MongoDB!")

	order.DisplayId, err = os.GetOrderDisplayId()
	if err != nil {
		return order, err
	}

	totalCost := 0.0
	totalSalePrice := 0.0

	items_cost, err := os.CalculateCost(order.Items)
	if err != nil {
		return order, err
	}

	for index, recipe_cost := range items_cost {

		order.Items[index].Cost = recipe_cost.Cost
		order.Items[index].SalePrice = recipe_cost.SalePrice

		totalCost += recipe_cost.Cost
		totalSalePrice += recipe_cost.SalePrice
	}

	order.SalePrice = totalSalePrice - order.Discount
	order.Cost = totalCost
	order.SubmittedAt = time.Now()
	order.Id = primitive.NewObjectID().Hex()

	if order.State != "stashed" {
		order.State = "pending"
	}

	_, err = client.Database("waha").Collection("orders").InsertOne(ctx, order)
	if err != nil {
		return order, err
	}

	return order, err
}

// GetOrdersParameters is the struct to hold the parameters for the GetOrders method.
type GetOrdersParameters struct {
	// OrderDisplayIdContains is the string to search for in the display_id field.
	OrderDisplayIdContains string
	// page_number is used in pagination to set the index of the first record to be returned.
	PageNumber int
	// Rows is to set the desired row count limit.
	PageSize int
	// FilterIsPaid is to filter for orders that are paid 0 (false), 1 (true), -1 (any)
	FilterIsPaid int8
	// IsPayLater is used to filter for is_pay_later orders 0 (unpaid), 1 (paid), -1 (any)
	IsPayLater int8
	// FilterState is used to filter for a specific state in_progress, finished, stashed, pending, cancelled, !stashed (! notation can be used to filter for negative values)
	FilterState []string
}

// GetOrders retrieves all orders from the database by default,
// If the OrderDisplayIdContains parameter is not empty,
// the function will check if display_id is not ""
// then it will filter for all order that contains the specified string
func (os *OrderService) GetOrders(params GetOrdersParameters) (orders []models.Order, totalRecords int64, err error) {

	orders = make([]models.Order, 0)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return orders, 0, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return orders, 0, err
	}

	// Connected successfully
	os.Logger.Info("Connected to MongoDB!")
	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetLimit(int64(params.PageSize))
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))

	if params.FilterIsPaid != -1 {
		if params.FilterIsPaid == 1 {
			filter["is_paid"] = true
		}
		if params.FilterIsPaid == 0 {
			filter["is_paid"] = false
		}
	}

	if params.OrderDisplayIdContains != "" {
		filter["display_id"] = bson.M{
			"$regex": fmt.Sprintf("(?i).*%s.*", params.OrderDisplayIdContains),
		}
	}

	positiveStateFilters := []string{}
	negativeStateFilter := []string{}

	for _, state := range params.FilterState {
		if state[0] == '!' {
			negativeStateFilter = append(negativeStateFilter, state[1:])
		} else {
			positiveStateFilters = append(positiveStateFilters, state)
		}
	}

	if len(positiveStateFilters) > 0 {
		filter["$and"] = []bson.M{
			{"state": bson.M{"$in": positiveStateFilters}},
			{"state": bson.M{"$nin": negativeStateFilter}},
		}
	}

	if params.IsPayLater == 1 {
		filter["is_pay_later"] = bson.M{"$eq": true}
	} else if params.IsPayLater == 0 {
		filter["is_pay_later"] = bson.M{"$eq": false}
	}

	totalRecords, err = client.Database("waha").Collection("orders").CountDocuments(ctx, bson.D{})
	if err != nil {
		os.Logger.Error(err.Error())
		return
	}

	cursor, err := client.Database("waha").Collection("orders").Find(context.Background(), filter, findOptions)
	if err != nil {
		return orders, 0, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return orders, 0, err
		}

		orders = append(orders, order)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return orders, 0, err
	}

	return orders, totalRecords, nil

}

func (os *OrderService) ConsumeOrderComponents(order models.Order) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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
	// Connected to Mongo
	materialService := MaterialService{
		Config:   os.Config,
		Logger:   os.Logger,
		Settings: os.Settings,
	}

	refined_notifications := map[string]models.WebsocketTopicServerMessage{}

	for itemIndex, item := range order.Items {
		notifications, err := materialService.ConsumeItemComponentsForOrder(item, order, itemIndex)
		for _, notification := range notifications {
			if _, ok := refined_notifications[notification.Key]; !ok {
				refined_notifications[notification.Key] = notification
			}
		}

		if len(refined_notifications) > 0 {
			notificationService, err := SpawnNotificationSingletonSvc("melody", os.Logger, os.Config)
			if err != nil {
				return err
			}
			for _, notification := range refined_notifications {

				json_notification, err := json.Marshal(notification)
				if err != nil {
					return err
				}

				notificationService.SendToTopic(notification.TopicName, string(json_notification))
			}
		}

		if err != nil {
			return err
		}

	}

	return nil
}

// StartOrder sets the state of the order with the given order_id to "in_progress",
// and updates the "started_at" field with the current time.
// It also consumes the item components from the inventory and sends a notification to the websockets.
func (os *OrderService) StartOrder(order_id string, order_items []models.OrderItem) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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
	os.Logger.Info("Connected to MongoDB!")

	var order models.Order

	err = client.Database("waha").Collection("orders").FindOne(context.Background(), bson.M{"id": order_id}).Decode(&order)
	if err != nil {
		return err
	}

	// decrease the ingredient component quantity from the components inventory
	err = os.ConsumeOrderComponents(order)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"items":      order_items,
			"state":      "in_progress",
			"started_at": time.Now(),
		},
	}

	_, err = client.Database("waha").Collection("orders").UpdateOne(context.Background(), bson.M{"id": order_id}, update)
	if err != nil {
		return err
	}

	logs_data := bson.M{
		"type":          "order_Start",
		"date":          time.Now(),
		"order_details": order,
	}
	_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
	if err != nil {
		return err
	}

	return nil
}

// GetOrder retrieves an order from the database with the given order_id.
func (os *OrderService) GetOrder(order_id string) (models.Order, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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

	coll := client.Database("waha").Collection("orders")
	var order models.Order

	err = coll.FindOne(ctx, bson.M{"id": order_id}).Decode(&order)
	if err != nil {
		os.Logger.Error(err.Error())
		return order, err
	}

	return order, nil
}
