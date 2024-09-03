package services

import (
	"context"
	"fmt"
	"log"
	"math"
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
	Logger logger.ILogger
	Config config.Config
}

func (os *OrderService) CalculateCost(items []models.RecipeSelections) (cost []models.ItemCost, err error) {

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
			ItemName: items[itemIndex].RecipeName,
			Cost:     0.0,
			RecipeId: items[itemIndex].Id,
		}

		for _, selection := range item.Selections {

			itemComponent := struct {
				ComponentName string
				ComponentId   string
				EntryId       string
				Quantity      float64
				Cost          float64
			}{
				ComponentName: selection.Name,
				ComponentId:   selection.ComponentId,
				EntryId:       selection.Entry.Id.Hex(),
				Quantity:      selection.Quantity,
			}

			var component_with_specific_entry models.Component

			component_id, _ := primitive.ObjectIDFromHex(selection.ComponentId)
			entry_id := selection.Entry.Id

			err = client.Database("waha").Collection("components").FindOne(
				context.Background(), bson.M{"_id": component_id, "entries._id": entry_id}, options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&component_with_specific_entry)

			if err == nil {
				quantity_cost := (component_with_specific_entry.Entries[0].Price / float64(component_with_specific_entry.Entries[0].PurchaseQuantity)) * float64(selection.Quantity)

				// check if cost is positive or negative infinity (semantic bug in calculation that causes problems later on)
				if math.IsInf(quantity_cost, 0) || math.IsInf(quantity_cost, -1) {
					quantity_cost = 0
				}

				itemCost.Cost += quantity_cost
				itemComponent.Cost = quantity_cost

			}

			itemCost.Components = append(itemCost.Components, itemComponent)

		}

		var recipe models.Recipe

		recipeID, _ := primitive.ObjectIDFromHex(items[itemIndex].Id)

		err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"_id": recipeID}).Decode(&recipe)

		if err != nil {
			panic(err)
		}

		for _, subrecipe := range item.SubRecipes {
			total_cost, err := os.CalculateCost([]models.RecipeSelections{subrecipe})
			if err != nil {
				return cost, err
			}

			for _, subrecipe_cost := range total_cost {
				itemCost.Cost += subrecipe_cost.Cost * float64(subrecipe.Quantity)
				cost = append(cost, subrecipe_cost)
			}
		}

		itemCost.SalePrice = recipe.Price
		cost = append(cost, itemCost)
	}

	return cost, err
}

func (os *OrderService) FinishOrder2(finish_order_request dto.FinishOrderRequest) (err error) {

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

	var order models.Order2
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
		"items":         []models.ItemCost{recipes_cost[0]},
		"order_id":      finish_order_request.Id,
		"time_consumed": time.Since(order.SubmittedAt),
	}
	_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (os *OrderService) FinishOrder(finish_order_request dto.FinishOrderRequest) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	var order dto.Order
	err = collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return err
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

			var component_with_specific_entry models.Component

			component_id, _ := primitive.ObjectIDFromHex(ingredient.ComponentId)
			entry_id, _ := primitive.ObjectIDFromHex(ingredient.EntryId)

			err = client.Database("waha").Collection("components").FindOne(
				context.Background(), bson.M{"_id": component_id, "entries._id": entry_id}, options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&component_with_specific_entry)

			if err == nil {
				quantity_cost := (component_with_specific_entry.Entries[0].Price / float64(component_with_specific_entry.Entries[0].PurchaseQuantity)) * float64(ingredient.Quantity)

				// check if cost is positive or negative infinity (semantic bug in calculation that causes problems later on)
				if math.IsInf(quantity_cost, 0) || math.IsInf(quantity_cost, -1) {
					quantity_cost = 0
				}

				totalCost += quantity_cost
				itemCost.Cost += quantity_cost
				itemComponent.Cost = quantity_cost
			}

			itemCost.Components = append(itemCost.Components, itemComponent)

		}

		var recipe models.Recipe

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

	return err
}

func (os *OrderService) SubmitOrder(order dto.Order) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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
	os.Logger.Info("Connected to MongoDB!")

	order_data := bson.M{
		"items":        order.Items,
		"submitted_at": time.Now(),
	}
	_, err = client.Database("waha").Collection("orders").InsertOne(ctx, order_data)
	if err != nil {
		return err
	}

	return err
}

func (os *OrderService) GetOrders() (orders []dto.Order, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		var order dto.Order
		if err := cursor.Decode(&order); err != nil {
			return orders, err
		}

		for i, item := range order.Items {

			var recipe models.Recipe

			err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"name": item.Name}).Decode(&recipe)
			if err == nil {
				order.Items[i].Recipe = recipe
			}

		}

		orders = append(orders, order)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return orders, err
	}

	return orders, nil

}

func (os *OrderService) StartOrder2(order_start_request dto.OrderStartRequest2) error {
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

	var order models.Order2

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
		Config: os.Config,
		Logger: os.Logger,
	}

	for itemIndex, item := range order_start_request.Items {
		err := componentService.ConsumeItemComponentsForOrder(item, order.Id, itemIndex)
		if err != nil {
			return err
		}
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

	var order dto.Order

	order_id, err := primitive.ObjectIDFromHex(order_start_request.Id)
	if err != nil {
		panic(err)
	}

	err = client.Database("waha").Collection("orders").FindOne(context.Background(), bson.M{"_id": order_id}).Decode(&order)
	if err != nil {
		return err
	}

	// decrease the ingredient component quantity from the components inventory

	for itemIndex, ingredient := range order_start_request.Ingredients {
		for componentIndex, component := range ingredient {
			if err != nil {
				return err
			}

			filter := bson.M{"name": component.Name, "entries.company": component.Company}
			// Define the update operation
			update := bson.M{
				"$inc": bson.M{
					"entries.$.quantity": -component.Quantity,
				},
			}

			var db_component models.Component

			for _, entry := range db_component.Entries {
				if entry.Company == component.Company {
					component.EntryId = entry.Id.Hex()
				}
			}

			err = client.Database("waha").Collection("components").FindOne(context.Background(), filter).Decode(&db_component)
			if err != nil {
				return err
			}

			component.ComponentId = db_component.Id
			order_start_request.Ingredients[itemIndex][componentIndex] = component

			_, err = client.Database("waha").Collection("components").UpdateOne(context.Background(), filter, update)
			if err != nil {
				return err
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
				return err
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

func (os *OrderService) GetOrder(order_id string) (dto.Order, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

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

	coll := client.Database("waha").Collection("orders")
	var order dto.Order

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
