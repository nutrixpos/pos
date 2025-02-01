package services

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomersService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

type GetCustomersParams struct {
	// page_number is used in pagination to set the index of the first record to be returned.
	PageNumber int
	// Rows is to set the desired row count limit.
	PageSize int
}

func (cs CustomersService) GetCustomers(params GetCustomersParams) (customers []models.Customer, customers_count int, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if cs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return customers, customers_count, err
	}
	// connected to db

	customers = make([]models.Customer, 0)

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return customers, customers_count, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return customers, customers_count, err
		}

		customers = append(customers, customer)
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return customers, customers_count, err
	}
	customers_count = int(count)

	return customers, customers_count, err
}

func (cs CustomersService) GetCustomer(customer_id string) (customer models.Customer, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if cs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	// connected to db

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	err = collection.FindOne(ctx, bson.M{"id": customer_id}).Decode(&customer)
	if err != nil {
		return customer, err
	}

	return
}

func (cs CustomersService) InsertNew(customer models.Customer) (afterInsert models.Customer, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if cs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	// connected to db

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	customer.Id = primitive.NewObjectID().Hex()

	result, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return afterInsert, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&afterInsert)
	if err != nil {
		return afterInsert, err
	}

	return
}

func (cs CustomersService) UpdateCustomer(customer models.Customer, customer_id string) (afterUpdate models.Customer, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if cs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	// connected to db

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	update := bson.M{}
	if customer.Name != "" {
		update["name"] = customer.Name
	}
	if customer.Phone != "" {
		update["phone"] = customer.Phone
	}
	if customer.Address != "" {
		update["address"] = customer.Address
	}

	update["id"] = customer_id

	_, err = collection.UpdateOne(ctx, bson.M{"id": customer_id}, bson.M{"$set": update})
	if err != nil {
		return afterUpdate, err
	}

	err = collection.FindOne(ctx, bson.M{"id": customer_id}).Decode(&afterUpdate)
	if err != nil {
		return afterUpdate, err
	}

	return
}

func (cs CustomersService) DeleteCustomer(customer_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if cs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	// connected to db

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	_, err = collection.DeleteOne(ctx, bson.M{"id": customer_id})
	if err != nil {
		return err
	}

	return
}
