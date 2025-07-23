// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MaterialService provides methods to manage and manipulate materials.
// It contains methods for calculating costs, checking availability, and
// updating material entries in the database. It relies on a logger for
// logging operations and a configuration for database connectivity.
type MaterialService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

type GetMaterialEntriesParams struct {
	// PageNumber sets the first index of the record to begin with in the select transaction
	PageNumber int
	// PageSize sets the limit of the number of desired rows
	PageSize int
	// Search is a text that the function should use to search for products that has a title contains the contians string
	Search string
}

func (rs *MaterialService) GetMaterialEntries(material_id string, params GetMaterialEntriesParams) (entries []models.MaterialEntry, totalRecords int64, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	entries = make([]models.MaterialEntry, 0)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		rs.Logger.Error(err.Error())
		return entries, totalRecords, err
	}

	collection := client.Database(rs.Config.Databases[0].Database).Collection("materials")
	// findOptions.SetSort(bson.M{"name": 1})
	// findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	// findOptions.SetLimit(int64(params.PageSize))

	// Get the total number of entries
	entryCountPipeline := []bson.M{
		{"$match": bson.M{"id": material_id}},
		{"$project": bson.M{"entryCount": bson.M{"$size": "$entries"}}},
	}

	entryCountCursor, err := collection.Aggregate(ctx, entryCountPipeline)
	if err != nil {
		return entries, totalRecords, err
	}
	defer entryCountCursor.Close(ctx)

	totalRecords = 0

	var entryCountResult []bson.M
	if err = entryCountCursor.All(ctx, &entryCountResult); err != nil {
		return entries, totalRecords, err
	}
	if len(entryCountResult) > 0 {
		totalRecords = int64(entryCountResult[0]["entryCount"].(int32))
	}

	skip := (params.PageNumber) * params.PageSize

	// Create aggregation pipeline
	pipeline := []bson.M{
		{"$match": bson.M{"id": material_id}},
		{"$project": bson.M{
			"entries": bson.M{
				"$slice": []interface{}{"$entries", skip, params.PageSize},
			},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entries, totalRecords, err
	}
	defer cursor.Close(ctx)

	// Get results
	var results []models.Material
	if err = cursor.All(ctx, &results); err != nil {
		return entries, totalRecords, err
	}

	if len(results) == 0 {
		return entries, totalRecords, err
	}

	return results[0].Entries, totalRecords, err
}

func (ms *MaterialService) GetMaterial(material_id string) (material models.Material, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ms.Config.Databases[0].Host, ms.Config.Databases[0].Port))

	timeout := 1000 * time.Second

	if ms.Config.Env == "dev" {
		timeout = 5 * time.Minute
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	collection := client.Database(ms.Config.Databases[0].Database).Collection("materials")
	err = collection.FindOne(ctx, bson.M{"id": material_id}).Decode(&material)
	if err != nil {
		return
	}

	return material, err
}

func (ms *MaterialService) Waste(entry_id, material_id string, quantity float64, order_id string, reason string, is_consume bool, user_id string) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ms.Config.Databases[0].Host, ms.Config.Databases[0].Port))

	timeout := 1000 * time.Second

	if ms.Config.Env == "dev" {
		timeout = 5 * time.Minute
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	if is_consume {
		var material models.Material
		err = client.Database(ms.Config.Databases[0].Database).Collection("materials").FindOne(context.Background(), bson.M{
			"id":         material_id,
			"entries.id": entry_id,
		},
			options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&material)
		if err != nil {
			return err
		}
		ms.ConsumeFromInventory(material, material.Entries[0].Id, quantity, reason, order_id, user_id)
	}

	filter := bson.M{"id": material_id, "entries.id": entry_id}
	// Define the update operation
	update := bson.M{
		"$inc": bson.M{
			"entries.$.quantity": -quantity,
		},
	}

	_, err = client.Database(ms.Config.Databases[0].Database).Collection("materials").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	log_material_return := models.LogWasteMaterial{
		Log: models.Log{
			Type:   models.LogTypeMaterialWaste,
			Date:   time.Now(),
			Id:     primitive.NewObjectID().Hex(),
			UserId: entry_id,
		},
		MaterialId: material_id,
		EntryId:    entry_id,
		IsConsume:  is_consume,
		OrderId:    order_id,
		Quantity:   quantity,
		Reason:     reason,
	}

	logs_collection := client.Database(ms.Config.Databases[0].Database).Collection("logs")
	_, err = logs_collection.InsertOne(ctx, log_material_return)

	if err != nil {
		return err
	}

	return nil
}

func (ms *MaterialService) InventoryReturn(entry_id, material_id string, quantity float64, order_id string, reason string, is_refunded bool, user_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ms.Config.Databases[0].Host, ms.Config.Databases[0].Port))

	timeout := 1000 * time.Second

	if ms.Config.Env == "dev" {
		timeout = 5 * time.Minute
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	filter := bson.M{"id": material_id, "entries.id": entry_id}
	// Define the update operation
	update := bson.M{
		"$inc": bson.M{
			"entries.$.quantity": quantity,
		},
	}

	_, err = client.Database(ms.Config.Databases[0].Database).Collection("materials").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if is_refunded {

		filter = bson.M{"id": order_id, "items.materials": bson.M{
			"$elemMatch": bson.M{
				"material.id": material_id,
				"entry.id":    entry_id,
			},
		}}

		update = bson.M{
			"$set": bson.M{
				"items.$[item].materials.$[material].entry.is_refunded":   true,
				"items.$[item].materials.$[material].entry.refund_reason": reason,
			},
		}

		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"item.materials.material.id": material_id,
					"item.materials.entry.id":    entry_id,
				},
				bson.M{
					"material.material.id": material_id,
					"material.entry.id":    entry_id,
				},
			},
		}

		opts := options.Update().SetArrayFilters(arrayFilters)

		_, err = client.Database(ms.Config.Databases[0].Database).Collection("orders").UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err
		}
	}

	log_material_return := models.LogMaterialInventoryReturn{
		Log: models.Log{
			Type:   models.LogTypeMaterialInventoryReturn,
			Date:   time.Now(),
			Id:     primitive.NewObjectID().Hex(),
			UserId: user_id,
		},
		OrderId:  order_id,
		Quantity: quantity,
		Reason:   reason,
	}

	logs_collection := client.Database(ms.Config.Databases[0].Database).Collection("logs")
	_, err = logs_collection.InsertOne(ctx, log_material_return)

	if err != nil {
		return err
	}

	return nil
}

func (cs *MaterialService) ConsumeFromInventory(material models.Material, entry_id string, quantity float64, reason string, order_id string, user_id string) (notifications []models.WebsocketTopicServerMessage, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return notifications, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return notifications, err
	}

	// Connected successfully

	if err != nil {
		return notifications, err
	}

	material_available_amount, err := cs.GetMaterialEntryAvailability(material.Id, entry_id)
	if err != nil {
		return notifications, err
	}

	if float64(material_available_amount) < quantity || material_available_amount < 0 {
		notifications = append(notifications, models.WebsocketTopicServerMessage{
			TopicName: "inventory_insufficient",
			Type:      "topic_message",
			Severity:  "error",
			Message:   fmt.Sprintf("Inventory for %s is insufficient, quantity requested by order_id: %s is %f, but entry %s only has %f", material.Name, order_id, quantity, entry_id, float64(material_available_amount)),
			Key:       fmt.Sprintf("inventory_insufficient@%s", material.Id),
		})
		return notifications, fmt.Errorf("entry %s is insufficient", material.Id)
	}

	filter := bson.M{"id": material.Id, "entries.id": entry_id}
	// Define the update operation
	update := bson.M{
		"$inc": bson.M{
			"entries.$.quantity": -quantity,
		},
	}

	_, err = client.Database(cs.Config.Databases[0].Database).Collection("materials").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return notifications, err
	}

	logs_data := models.LogMaterialConsume{
		Log: models.Log{
			Id:     primitive.NewObjectID().Hex(),
			Type:   "component_consume",
			Date:   time.Now(),
			UserId: user_id,
		},
		MaterialId: material.Id,
		EntryId:    entry_id,
	}
	_, err = client.Database(cs.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, logs_data)
	if err != nil {
		return notifications, err
	}

	available_quantity, err := cs.GetComponentAvailability(material.Id)
	if err != nil {
		return notifications, err
	}

	if float64(available_quantity) <= cs.Settings.Inventory.StockAlertTreshold {
		notifications = append(notifications, models.WebsocketTopicServerMessage{
			TopicName: "inventory_low",
			Type:      "topic_message",
			Severity:  "warn",
			Message:   fmt.Sprintf("Inventory for %s is low: %f", material.Name, float64(available_quantity)),
			Key:       fmt.Sprintf("low_inventiry@%s", material.Id),
		})
	}

	return notifications, err
}

// CalculateMaterialCost calculates the cost of a material entry based on its ID, material ID, and quantity.
// It connects to the MongoDB database, retrieves the specific material entry, and calculates the cost
// using the purchase price and purchase quantity.
func (cs *MaterialService) CalculateMaterialCost(entry_id, material_id string, quantity float64) (cost float64, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return 0, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return 0, err
	}

	// Connected successfully
	cs.Logger.Info("Connected to MongoDB!")

	var material models.Material
	err = client.Database(cs.Config.Databases[0].Database).Collection("materials").FindOne(context.Background(), bson.M{
		"id":         material_id,
		"entries.id": entry_id,
	},
		options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&material)
	if err != nil {
		return 0, err
	}

	if len(material.Entries) == 0 {
		return 0.0, fmt.Errorf("entry %s not found in material %s", entry_id, material_id)
	}

	cost = (material.Entries[0].PurchasePrice / float64(material.Entries[0].PurchaseQuantity)) * quantity

	return cost, nil
}

// GetMaterialEntryAvailability retrieves the quantity of a specific material entry
// from the database.
//
// The function takes a material ID and an entry ID as parameters and returns the
// quantity of the specified entry in the material. If the entry is not found in
// the material, the function returns an error.
func (cs *MaterialService) GetMaterialEntryAvailability(material_id string, entry_id string) (amount float64, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return amount, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return amount, err
	}

	var material models.Material
	err = client.Database(cs.Config.Databases[0].Database).Collection("materials").FindOne(context.Background(), bson.M{
		"id":         material_id,
		"entries.id": entry_id,
	},
		options.FindOne().SetProjection(bson.M{"entries.$": 1})).Decode(&material)
	if err != nil {
		return 0.0, err
	}

	if len(material.Entries) == 0 {
		return 0.0, fmt.Errorf("entry %s not found in material %s", entry_id, material_id)
	}

	amount = material.Entries[0].Quantity

	return amount, err

}

// ConsumeItemComponentsForOrder consumes components for an order item, and returns the notifications to be sent via websocket.
// It returns an error if something goes wrong.
func (cs *MaterialService) ConsumeItemComponentsForOrder(item models.OrderItem, order models.Order, order_item_index int, user_id string) (notifications []models.WebsocketTopicServerMessage, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return notifications, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return notifications, err
	}

	// Connected successfully
	cs.Logger.Info("Connected to MongoDB!")

	productService := RecipeService{
		Logger: cs.Logger,
		Config: cs.Config,
	}

	if item.IsConsumeFromReady {
		err = productService.ConsumeFromReady(item.Product.Id, item.Quantity)
		return notifications, err
	}

	for _, component := range item.Materials {
		if err != nil {
			return notifications, err
		}

		material_available_amount, err := cs.GetMaterialEntryAvailability(component.Material.Id, component.Entry.Id)
		if err != nil {
			return notifications, err
		}

		if float64(material_available_amount) < (component.Quantity*item.Quantity) || material_available_amount < 0 {
			notifications = append(notifications, models.WebsocketTopicServerMessage{
				TopicName: "inventory_insufficient",
				Type:      "topic_message",
				Severity:  "error",
				Message:   fmt.Sprintf("Inventory for %s is insufficient, quantity requested by order_id: %s (display_id: %s) is %f, but entry %s only has %f", component.Material.Name, order.Id, order.DisplayId, (component.Quantity * item.Quantity), component.Entry.Id, float64(material_available_amount)),
				Key:       fmt.Sprintf("inventory_insufficient@%s", component.Material.Id),
			})
			return notifications, fmt.Errorf("entry %s is insufficient", component.Material.Id)
		}

		filter := bson.M{"id": component.Material.Id, "entries.id": component.Entry.Id}
		// Define the update operation
		update := bson.M{
			"$inc": bson.M{
				"entries.$.quantity": -component.Quantity * item.Quantity,
			},
		}

		_, err = client.Database(cs.Config.Databases[0].Database).Collection("materials").UpdateOne(context.Background(), filter, update)
		if err != nil {
			return notifications, err
		}

		logs_data := bson.M{
			"type":             "component_consume",
			"date":             time.Now(),
			"id":               primitive.NewObjectID().Hex(),
			"component_id":     component.Material.Id,
			"quantity":         component.Quantity * item.Quantity,
			"entry_id":         component.Entry.Id,
			"order_id":         order.Id,
			"recipe_id":        item.Product.Id,
			"order_item_index": order_item_index,
			"user_id":          user_id,
		}
		_, err = client.Database(cs.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			return notifications, err
		}

		quantity, err := cs.GetComponentAvailability(component.Material.Id)
		if err != nil {
			return notifications, err
		}

		if float64(quantity) <= cs.Settings.Inventory.StockAlertTreshold {
			notifications = append(notifications, models.WebsocketTopicServerMessage{
				TopicName: "inventory_low",
				Type:      "topic_message",
				Severity:  "warn",
				Message:   fmt.Sprintf("Inventory for %s is low: %f", component.Material.Name, float64(quantity)),
				Key:       fmt.Sprintf("low_inventiry@%s", component.Material.Id),
			})
		}

	}

	for _, subrecipe := range item.SubItems {

		sub_notifications, err := cs.ConsumeItemComponentsForOrder(subrecipe, order, order_item_index, user_id)
		if err != nil {
			return notifications, err
		}

		notifications = append(notifications, sub_notifications...)
	}

	return notifications, nil
}

// GetComponentAvailability retrieves the total quantity of a specific component.
//
// The function takes a component ID as a parameter and returns the total
// quantity of the specified component in the database. If the component is not
// found in the database, the function returns an error.
//
// The function is used to check the availability of a specific component before
// consuming it. The component quantity is calculated by summing up the quantity
// of all entries of the component.
func (cs *MaterialService) GetComponentAvailability(componentid string) (amount float64, err error) {

	amount = 0.0

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return 0.0, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return 0.0, err
	}

	// Connected successfully
	cs.Logger.Info("Connected to MongoDB!")

	// Get the "test" collection from the database
	collection := client.Database(cs.Config.Databases[0].Database).Collection("materials")

	filter := bson.M{"id": componentid}
	var component models.Material
	err = collection.FindOne(ctx, filter).Decode(&component)
	if err != nil {
		return 0.0, err
	}

	for _, entry := range component.Entries {
		if entry.Quantity > 0 {
			amount += entry.Quantity
		}
	}

	return amount, nil

}

// GetMaterials retrieves all materials from the database.
//
// The function returns a slice of Material structs.
func (cs *MaterialService) GetMaterials(page_number int, page_size int) (materials []models.Material, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cs.Logger.Error(err.Error())
		return materials, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		cs.Logger.Error(err.Error())
		return materials, err
	}

	// Connected successfully
	fmt.Println("Connected to MongoDB!")

	materials = make([]models.Material, 0)
	// findOptions := options.Find().SetProjection(bson.M{"entries": 0})
	findOptions := options.Find()

	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(page_size))

	// Get the "materials" collection from the database
	cur, err := client.Database(cs.Config.Databases[0].Database).Collection("materials").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		cs.Logger.Error(err.Error())
		return materials, err
	}

	defer cur.Close(ctx)

	// Iterate over the documents and print them as JSON
	for cur.Next(ctx) {
		var material models.Material
		err := cur.Decode(&material)
		if err != nil {
			cs.Logger.Error(err.Error())
			return materials, err
		}

		var sum float64
		for _, entry := range material.Entries {
			sum += entry.Quantity
		}

		material.Quantity = sum
		material.Entries = make([]models.MaterialEntry, 0)

		materials = append(materials, material)
	}

	return materials, nil

}

func (cs *MaterialService) DeleteMaterial(material_id string) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	// Connected successfully

	// Delete the material with the given ID
	filter := bson.M{"id": material_id}
	_, err = client.Database(cs.Config.Databases[0].Database).Collection("materials").DeleteOne(ctx, filter)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	return nil
}

// EditMaterial updates a material in the database.
//
// The function takes a MaterialEditRequest as a parameter and updates the
// corresponding material in the database with the provided information.
//
// The function is used to edit an existing material in the database.
func (cs *MaterialService) EditMaterial(material_id string, material_to_edit models.Material) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	// Connected successfully

	// Find material in db
	var existingMaterial models.Material
	err = client.Database(cs.Config.Databases[0].Database).Collection("materials").FindOne(context.Background(), bson.M{"id": material_id}).Decode(&existingMaterial)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	existingMaterial.Settings.StockAlertTreshold = material_to_edit.Settings.StockAlertTreshold
	existingMaterial.Name = material_to_edit.Name
	existingMaterial.Unit = material_to_edit.Unit

	// Update the material
	_, err = client.Database(cs.Config.Databases[0].Database).Collection("materials").UpdateOne(context.Background(), bson.M{"id": material_id}, bson.M{"$set": existingMaterial})
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	return nil
}

// AddComponent adds a new material component to the database.
// It first inserts the material into the "materials" collection,
// then logs the addition of each entry into the "logs" collection.
// If there is any error during the database operations, it returns the error.
func (cs *MaterialService) AddComponent(material models.Material, user_id string) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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

	material.Id = primitive.NewObjectID().Hex()

	for index, _ := range material.Entries {
		material.Entries[index].Id = primitive.NewObjectID().Hex()
	}

	// Insert the DBComponent struct into the "materials" collection
	collection := client.Database(cs.Config.Databases[0].Database).Collection("materials")
	_, err = collection.InsertOne(ctx, material)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	for _, entry := range material.Entries {

		logs_data := bson.M{
			"id":          primitive.NewObjectID().Hex(),
			"type":        "component_add",
			"date":        time.Now(),
			"material_id": material.Id,
			"entry_id":    entry.Id,
			"company":     entry.Company,
			"quantity":    entry.Quantity,
			"price":       entry.PurchasePrice,
			"user_id":     user_id,
		}
		_, err = client.Database(cs.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			cs.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// PushMaterialEntry adds a new entry to a material in the database.
//
// The function takes a component ID and a slice of MaterialEntry structs as parameters.
// It then finds the material with the given ID and appends the new entries to the material's
// entries array. If the material is not found, the function will return an error.
func (cs *MaterialService) PushMaterialEntry(componentId string, entries []models.MaterialEntry, user_id string) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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

	filter := bson.M{"id": componentId}

	for _, entry := range entries {

		entry.PurchaseQuantity = entry.Quantity
		entry_id := primitive.NewObjectID().Hex()

		entry_data := bson.M{
			"id":                entry_id,
			"purchase_quantity": entry.PurchaseQuantity,
			"price":             entry.PurchasePrice,
			"quantity":          entry.Quantity,
			"company":           entry.Company,
			"sku":               entry.SKU,
			"expiration_date":   entry.ExpirationDate,
		}

		update := bson.M{"$push": bson.M{"entries": entry_data}}
		opts := options.Update().SetUpsert(false)

		_, err = client.Database(cs.Config.Databases[0].Database).Collection("materials").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}

		logs_data := bson.M{
			"type":        "component_add",
			"id":          primitive.NewObjectID().Hex(),
			"date":        time.Now(),
			"material_id": componentId,
			"entry_id":    entry_id,
			"company":     entry.Company,
			"quantity":    entry.Quantity,
			"price":       entry.PurchasePrice,
			"user_id":     user_id,
		}
		_, err = client.Database(cs.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			cs.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// DeleteEntry deletes an entry from a material in the database.
//
// The function takes a entry ID and a component ID as parameters. It then finds
// the material with the given component ID and removes the entry with the given
// entry ID from the material's entries array. If the material or the entry is not
// found, the function will return an error.
func (cs *MaterialService) DeleteEntry(entryid string, componentid string) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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

	// Connect to the database
	collection := client.Database(cs.Config.Databases[0].Database).Collection("materials")

	// Find the component document and update the entries array
	filter := bson.M{"id": componentid}
	update := bson.M{"$pull": bson.M{"entries": bson.M{"id": entryid}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
