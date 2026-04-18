// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryService struct {
	Logger logger.ILogger
	Config config.Config
}

// InsertCategory inserts a new category into the database.
func (cs *CategoryService) InsertCategory(category models.Category) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return err
	}

	category.Id = primitive.NewObjectID().Hex()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.InsertOne(ctx, category)

	return err
}

// DeleteCategory deletes a category from the database.
func (cs *CategoryService) DeleteCategory(category_id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return err
	}

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.DeleteOne(ctx, bson.M{"id": category_id})

	return err
}

// UpdateCategory updates a category in the database.
func (cs *CategoryService) UpdateCategory(category models.Category) (updatedCategory models.Category, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return updatedCategory, err
	}

	collection := client.Database(cs.Config.Databases[0].Database).Collection("categories")
	_, err = collection.UpdateOne(ctx, bson.M{"id": category.Id}, bson.M{"$set": bson.M{
		"name":     category.Name,
		"products": category.Products,
	}})

	return category, err

}

func (cs *CategoryService) GetCategories(page_number int, page_size int) (categories []models.Category, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return categories, err
	}

	categories = make([]models.Category, 0)

	findOptions := options.Find()

	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(page_size))
	cur, err := client.Database(cs.Config.Databases[0].Database).Collection("categories").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return categories, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var category models.Category
		err := cur.Decode(&category)
		if err != nil {
			return categories, err
		}

		products := []models.Product{}

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

		contentCategory := models.Category{
			Name:     category.Name,
			Id:       category.Id,
			Products: products,
		}

		categories = append(categories, contentCategory)
	}

	return categories, nil
}
