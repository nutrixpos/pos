package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ComponentService struct {
	Logger logger.ILogger
	Config config.Config
}

func (cs *ComponentService) CalculateMaterialCost(entry_id, material_id string, quantity float64) (cost float64, err error) {
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

	material_id_hex, err := primitive.ObjectIDFromHex(material_id)
	if err != nil {
		return 0, err
	}

	entry_id_hex, err := primitive.ObjectIDFromHex(entry_id)
	if err != nil {
		return 0, err
	}

	var material models.Material
	err = client.Database("waha").Collection("components").FindOne(context.Background(), bson.M{
		"_id":         material_id_hex,
		"entries._id": entry_id_hex,
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

func (cs *ComponentService) ConsumeItemComponentsForOrder(rs models.OrderItem, order_id string, item_order_index int) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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
	cs.Logger.Info("Connected to MongoDB!")

	for _, component := range rs.Materials {
		if err != nil {
			return err
		}

		component_ob_id, err := primitive.ObjectIDFromHex(component.Material.Id)
		if err != nil {
			return err
		}

		filter := bson.M{"_id": component_ob_id, "entries._id": component.Entry.Id}
		// Define the update operation
		update := bson.M{
			"$inc": bson.M{
				"entries.$.quantity": -component.Quantity,
			},
		}

		_, err = client.Database("waha").Collection("components").UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		logs_data := bson.M{
			"type":             "component_consume",
			"date":             time.Now(),
			"component_id":     component.Material.Id,
			"quantity":         component.Quantity,
			"entry_id":         component.Entry.Id,
			"order_id":         order_id,
			"recipe_id":        rs.Product.Id,
			"item_order_index": item_order_index,
		}
		_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			return err
		}

	}

	for _, subrecipe := range rs.SubItems {

		err = cs.ConsumeItemComponentsForOrder(subrecipe, order_id, item_order_index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs *ComponentService) GetComponentAvailability(componentid string) (amount float32, err error) {

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
	collection := client.Database("waha").Collection("components")
	objectID, err := primitive.ObjectIDFromHex(componentid)
	if err != nil {
		return 0.0, err
	}

	filter := bson.M{"_id": objectID}
	var component models.Material
	err = collection.FindOne(ctx, filter).Decode(&component)
	if err != nil {
		return 0.0, err
	}

	for _, entry := range component.Entries {
		amount += entry.Quantity
	}

	return amount, nil

}

func (cs *ComponentService) GetComponents() (components []models.Material, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cs.Logger.Error(err.Error())
		return components, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		cs.Logger.Error(err.Error())
		return components, err
	}

	// Connected successfully
	fmt.Println("Connected to MongoDB!")

	// Get the "test" collection from the database
	cur, err := client.Database("waha").Collection("components").Find(ctx, bson.D{})
	if err != nil {
		cs.Logger.Error(err.Error())
		return components, err
	}

	defer cur.Close(ctx)

	// Iterate over the documents and print them as JSON
	for cur.Next(ctx) {
		var component models.Material
		err := cur.Decode(&component)
		if err != nil {
			cs.Logger.Error(err.Error())
			return components, err
		}
		components = append(components, component)
	}

	return components, nil

}

func (cs *ComponentService) AddComponent(component models.Material) error {

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

	// Insert the DBComponent struct into the "components" collection
	collection := client.Database("waha").Collection("components")
	_, err = collection.InsertOne(ctx, component)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	for _, entry := range component.Entries {

		logs_data := bson.M{
			"type":     "component_add",
			"date":     time.Now(),
			"company":  entry.Company,
			"quantity": entry.Quantity,
			"price":    entry.PurchasePrice,
		}
		_, err = client.Database("waha").Collection("logs").InsertOne(ctx, logs_data)
		if err != nil {
			cs.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (cs *ComponentService) PushComponentEntry(componentId primitive.ObjectID, entries []models.MaterialEntry) error {

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

	filter := bson.M{"_id": componentId}

	for _, entry := range entries {

		entry.PurchaseQuantity = entry.Quantity

		entry_data := bson.M{
			"_id":               primitive.NewObjectID(),
			"purchase_quantity": entry.PurchaseQuantity,
			"price":             entry.PurchasePrice,
			"quantity":          entry.Quantity,
			"company":           entry.Company,
			"sku":               entry.SKU,
		}

		update := bson.M{"$push": bson.M{"entries": entry_data}}
		opts := options.Update().SetUpsert(false)

		_, err = client.Database("waha").Collection("components").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs *ComponentService) DeleteEntry(entryid primitive.ObjectID, componentid primitive.ObjectID) error {
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
	collection := client.Database("waha").Collection("components")

	// Find the component document and update the entries array
	filter := bson.M{"_id": componentid}
	update := bson.M{"$pull": bson.M{"entries": bson.M{"_id": entryid}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
