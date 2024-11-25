package services

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryService struct {
	Logger logger.ILogger
	Config config.Config
}

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

	collection := client.Database("waha").Collection("categories")
	_, err = collection.InsertOne(ctx, category)

	return err
}

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

	collection := client.Database("waha").Collection("categories")
	_, err = collection.DeleteOne(ctx, bson.M{"id": category_id})

	return err
}

func (cs *CategoryService) UpdateCategory(category models.Category) (err error) {

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

	collection := client.Database("waha").Collection("categories")
	_, err = collection.UpdateOne(ctx, bson.M{"id": category.Id}, bson.M{"$set": bson.M{
		"name":     category.Name,
		"products": category.Products,
	}})

	return err

}

func (cs *CategoryService) GetCategories(first int, rows int) (categories []models.Category, err error) {

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
	findOptions.SetSkip(int64(first))
	findOptions.SetLimit(int64(rows))
	cur, err := client.Database("waha").Collection("categories").Find(ctx, bson.D{}, findOptions)
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

		products := []models.CategoryProduct{}

		// Fetch the recipes for each category using the Recipe IDs
		for _, category_product := range category.Products {
			var product models.Product
			filter := bson.D{{Key: "id", Value: category_product.Id}}
			err := client.Database("waha").Collection("recipes").FindOne(ctx, filter).Decode(&product)
			if err != nil {
				cs.Logger.Error(`\nERROR: Recipe doesn't exist id: ${%v}`, category_product)
				continue
			}

			products = append(products, models.CategoryProduct{
				Id:   product.Id,
				Name: product.Name,
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
