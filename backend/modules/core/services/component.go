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

func (cs *ComponentService) GetComponents() (components []models.Component, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		var component models.Component
		err := cur.Decode(&component)
		if err != nil {
			cs.Logger.Error(err.Error())
			return components, err
		}
		components = append(components, component)
	}

	return components, nil

}

func (cs *ComponentService) AddComponent(component models.Component) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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

	// Insert the DBComponent struct into the "components" collection
	collection := client.Database("waha").Collection("components")
	_, err = collection.InsertOne(ctx, component)
	if err != nil {
		cs.Logger.Error(err.Error())
		return err
	}

	for _, entry := range component.Entries {

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
			cs.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (cs *ComponentService) PushComponentEntry(componentId primitive.ObjectID, entries []models.ComponentEntry) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

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

	filter := bson.M{"_id": componentId}

	for _, entry := range entries {

		entry.Id = primitive.NewObjectID()

		entry.PurchaseQuantity = entry.Quantity

		update := bson.M{"$push": bson.M{"entries": entry}}
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
