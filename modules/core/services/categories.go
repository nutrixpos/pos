// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
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

type CategoryService struct {
	Logger logger.ILogger
	Config config.Config
}

// InsertCategory inserts a new category into the database.
func (cs *CategoryService) InsertCategory(category models.Category) (err error) {

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

	category.Id = primitive.NewObjectID().Hex()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.InsertOne(ctx, category)

	return err
}

// DeleteCategory deletes a category from the database.
func (cs *CategoryService) DeleteCategory(category_id string) (err error) {

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

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.DeleteOne(ctx, bson.M{"id": category_id})

	return err
}

// UpdateCategory updates a category in the database.
func (cs *CategoryService) UpdateCategory(category models.Category) (updatedCategory models.Category, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return updatedCategory, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return updatedCategory, err
	}

	// Connected successfully

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.UpdateOne(ctx, bson.M{"id": category.Id}, bson.M{"$set": bson.M{
		"name":     category.Name,
		"products": category.Products,
	}})

	return category, err

}

// GetCategories returns a list of categories from the database.
func (cs *CategoryService) GetCategories(page_number int, page_size int) (categories []models.Category, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return categories, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return categories, err
	}

	// Connected successfully
	cs.Logger.Info("Connected to MongoDB!")

	// Fetch categories from the database
	findOptions := options.Find()

	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(page_size))
	cur, err := client.Database(cs.Config.Databases[0].Database).Collection("categories").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return categories, err
	}
	defer cur.Close(ctx)

	// Iterate through the categories
	for cur.Next(ctx) {
		var category models.Category
		err := cur.Decode(&category)
		if err != nil {
			return categories, err
		}

		products := []models.Product{}

		// Fetch the recipes for each category using the Recipe IDs
		for _, category_product := range category.Products {
			var product models.Product
			filter := bson.D{{Key: "id", Value: category_product.Id}}
			err := client.Database(cs.Config.Databases[0].Database).Collection("recipes").FindOne(ctx, filter).Decode(&product)
			if err != nil {
				cs.Logger.Error(`\nERROR: Recipe doesn't exist id: ${%v}`, category_product)
				continue
			}

			products = append(products, models.Product{
				Id: product.Id,
			})
		}

		// Create a CategoriesContentRequest_Category with embedded recipes
		contentCategory := models.Category{
			Name:     category.Name,
			Id:       category.Id,
			Products: products,
		}

		// Append the contentCategory to the categories slice
		categories = append(categories, contentCategory)
	}

	return categories, nil
}
