package services

import (
	"context"
	"fmt"
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

type CategoryService struct {
	Logger logger.ILogger
	Config config.Config
}

func (cs *CategoryService) GetCategories() (categories []dto.Category, err error) {

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
	cur, err := client.Database("waha").Collection("categories").Find(ctx, bson.D{})
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

		recipes := []models.Product{}

		// Fetch the recipes for each category using the Recipe IDs
		for _, recipeID := range category.Recipes {
			var recipe models.Product
			obj_id, _ := primitive.ObjectIDFromHex(recipeID)
			filter := bson.D{{Key: "_id", Value: obj_id}}
			err := client.Database("waha").Collection("recipes").FindOne(ctx, filter).Decode(&recipe)
			if err != nil {
				cs.Logger.Error(`\nERROR: Recipe doesn't exist id: ${%v}`, recipeID)
				continue
			}

			recipes = append(recipes, recipe)
		}

		// Create a CategoriesContentRequest_Category with embedded recipes
		contentCategory := dto.Category{
			Name:    category.Name,
			Recipes: recipes,
		}

		// Append the contentCategory to the categories slice
		categories = append(categories, contentCategory)
	}

	return categories, nil
}
