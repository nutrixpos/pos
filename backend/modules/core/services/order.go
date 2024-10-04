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
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
}

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

			component_id, _ := primitive.ObjectIDFromHex(component.Material.Id)
			entry_id, _ := primitive.ObjectIDFromHex(component.Entry.Id)

			err = client.Database("waha").Collection("components").FindOne(
				context.Background(), bson.M{"_id": component_id, "entries._id": entry_id}, options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&component_with_specific_entry)

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

		recipeID, _ := primitive.ObjectIDFromHex(items[itemIndex].Product.Id)

		err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"_id": recipeID}).Decode(&recipe)

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

func (os *OrderService) FinishOrder(finish_order_request dto.FinishOrderRequest) (err error) {

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
	Id, err := primitive.ObjectIDFromHex(finish_order_request.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": Id}
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

	recipes_cost, err := os.CalculateCost(order.Items)
	if err != nil {
		return err
	}

	for _, recipe_cost := range recipes_cost {
		totalCost += recipe_cost.Cost
		totalSalePrice += recipe_cost.SalePrice
	}

	logs_data := bson.M{
		"type":          "order_finish",
		"date":          time.Now(),
		"cost":          totalCost,
		"sale_price":    totalSalePrice,
		"items":         recipes_cost,
		"order_id":      finish_order_request.Id,
		"time_consumed": time.Since(order.SubmittedAt),
	}
	_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

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

	var settings config.Settings
	err = client.Database("waha").Collection("settings").FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		return order_display_id, err
	}

	random_queue_index := rand.Intn((len(settings.Orders.Queues)-1)-0) + 0
	random_queue := settings.Orders.Queues[random_queue_index]

	order_display_id = fmt.Sprintf("%s-%v", random_queue.Prefix, random_queue.Next)

	settings_obj_id, err := primitive.ObjectIDFromHex(settings.Id)
	if err != nil {
		return order_display_id, err
	}

	_, err = client.Database("waha").Collection("settings").UpdateOne(
		ctx,
		bson.M{"_id": settings_obj_id, "orders.queues.prefix": random_queue.Prefix},
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

func (os *OrderService) SubmitOrder(order models.Order) (err error) {

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
		return err
	}

	order_data := bson.M{
		"items":        order.Items,
		"submitted_at": time.Now(),
		"display_id":   order.DisplayId,
	}
	_, err = client.Database("waha").Collection("orders").InsertOne(ctx, order_data)
	if err != nil {
		return err
	}

	return err
}

func (os *OrderService) GetOrders() (orders []models.Order, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return orders, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return orders, err
	}

	// Connected successfully
	os.Logger.Info("Connected to MongoDB!")

	cursor, err := client.Database("waha").Collection("orders").Find(context.Background(), bson.M{
		"state": bson.M{
			"$ne": "finished",
		},
	})
	if err != nil {
		return orders, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return orders, err
	}

	return orders, nil

}

func (os *OrderService) StartOrder(order_start_request dto.OrderStartRequest) error {
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

	order_id, err := primitive.ObjectIDFromHex(order_start_request.Id)
	if err != nil {
		panic(err)
	}

	err = client.Database("waha").Collection("orders").FindOne(context.Background(), bson.M{"_id": order_id}).Decode(&order)
	if err != nil {
		return err
	}

	// decrease the ingredient component quantity from the components inventory

	componentService := ComponentService{
		Config:   os.Config,
		Logger:   os.Logger,
		Settings: os.Settings,
	}

	refined_notifications := map[string]models.WebsocketTopicServerMessage{}

	for itemIndex, item := range order_start_request.Items {
		notifications, err := componentService.ConsumeItemComponentsForOrder(item, order.Id, itemIndex)
		if err != nil {
			return err
		}

		for _, notification := range notifications {
			if _, ok := refined_notifications[notification.Key]; !ok {
				refined_notifications[notification.Key] = notification
			}
		}

	}

	notificationService, err := SpawnNotificationService("melody", os.Logger, os.Config)
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

	update := bson.M{
		"$set": bson.M{
			"items":      order_start_request.Items,
			"state":      "in_progress",
			"started_at": time.Now(),
		},
	}

	_, err = client.Database("waha").Collection("orders").UpdateOne(context.Background(), bson.M{"_id": order_id}, update)
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

	objID, err := primitive.ObjectIDFromHex(order_id)
	if err != nil {
		os.Logger.Error(err.Error())
		return order, err
	}

	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		os.Logger.Error(err.Error())
		return order, err
	}

	return order, nil
}
