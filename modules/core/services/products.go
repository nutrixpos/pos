// Package services contains the business logic of the core module of nutrix.
//
// It implements the required interfaces for the core module of nutrix.
//
// The services in this package are used to interact with the database using the
// models package, and to interact with the outside world using the dto package.
package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/customerrors"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RecipeService provides methods to manage recipes, including logging and configuration.
type RecipeService struct {
	Logger logger.ILogger
	Config config.Config
}

func (rs *RecipeService) GetProduct(product_id string) (product models.Product, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return product, err
	}
	// connected to db

	collection := client.Database("waha").Collection("recipes")
	err = collection.FindOne(ctx, bson.M{"id": product_id}).Decode(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

// UpdateProduct updates a product in the database.
//
// It takes a product and updates it in the database.
//
// If the product is not found, it will return an error.
func (rs *RecipeService) UpdateProduct(product_id string, product models.Product) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// connected to db

	collection := client.Database("waha").Collection("recipes")

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"id": product_id},
		bson.M{
			"$set": bson.M{
				"name":         product.Name,
				"materials":    product.Materials,
				"sub_products": product.SubProducts,
				"ready":        product.Ready,
				"recipeId":     product.Id,
				"price":        product.Price,
				"image_url":    product.ImageURL,
			},
		},
	)

	return err
}

// DeleteProduct deletes a product from the database.
//
// It takes a product_id and deletes it from the database.
func (rs *RecipeService) DeleteProduct(product_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// connected to db

	collection := client.Database("waha").Collection("recipes")
	_, err = collection.DeleteOne(ctx, bson.M{"id": product_id})
	if err != nil {
		return err
	}

	return err

}

// InsertNew inserts a new product into the database.
//
// It takes a product and inserts it into the database.
// It returns an error if the product could not be inserted.
func (rs *RecipeService) InsertNew(product models.Product) (afterInsert models.Product, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return afterInsert, err
	}
	// connected to db

	collection := client.Database("waha").Collection("recipes")

	product.Id = primitive.NewObjectID().Hex()

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return afterInsert, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&afterInsert)
	if err != nil {
		return afterInsert, err
	}

	return afterInsert, err
}

type GetProductsParams struct {
	// PageNumber sets the first index of the record to begin with in the select transaction
	PageNumber int
	// PageSize sets the limit of the number of desired rows
	PageSize int
	// Search is a text that the function should use to search for products that has a title contains the contians string
	Search string
}

// GetProducts retrieves a list of products from the database.
//
// It takes a first_index and rows and returns a slice of products.
// It also returns the total number of records in the database.
// It returns an error if the products could not be retrieved.
func (rs *RecipeService) GetProducts(params GetProductsParams) (products []models.Product, totalRecords int64, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		rs.Logger.Error(err.Error())
		return products, totalRecords, err
	}

	collection := client.Database("waha").Collection("recipes")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"name": 1})
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	filter := bson.M{}
	if params.Search != "" {
		filter["name"] = bson.M{
			"$regex": fmt.Sprintf("(?i).*%s.*", params.Search),
		}
	}

	totalRecords, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		rs.Logger.Error(err.Error())
		return products, totalRecords, err
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		rs.Logger.Error(err.Error())
		return products, totalRecords, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(context.Background()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return products, totalRecords, err
		}

		for index, sub_product := range product.SubProducts {

			err = collection.FindOne(ctx, bson.M{"id": sub_product.Id}, options.FindOne()).Decode(&sub_product)
			if err != nil {
				return products, totalRecords, err
			}
			product.SubProducts[index].Name = sub_product.Name
		}

		products = append(products, product)
	}

	return products, totalRecords, err
}

// ConsumeFromReady consumes a quantity from the ready stock of a product.
func (rs *RecipeService) ConsumeFromReady(product_id string, quantity float64) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	var product models.Product
	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"id": product_id}).Decode(&product)
	if err != nil {
		return err
	}

	ready := product.Ready

	if ready < quantity {
		return customerrors.ErrInsufficientReady
	}

	ready -= quantity

	_, err = client.Database("waha").Collection("recipes").UpdateOne(
		context.Background(),
		bson.M{"id": product_id},
		bson.M{
			"$set": bson.M{"ready": ready},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// FillRecipeDesign fills the recipe design for an order item with its product and sub products.
func (rs *RecipeService) FillRecipeDesign(item models.OrderItem) (models.OrderItem, error) {

	self_recipe, err := rs.GetRecipeTree(item.Product.Id)
	if err != nil {
		return item, err
	}

	item.Product = self_recipe

	for i, subrecipe := range item.SubItems {
		sub_item_recipe, err := rs.FillRecipeDesign(subrecipe)
		if err != nil {
			return item, err
		}
		item.SubItems[i] = sub_item_recipe
	}

	return item, err

}

// GetRecipeMaterials returns all materials for a given recipe.
func (rs *RecipeService) GetRecipeMaterials(recipe_id string) (materials []models.Material, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	var recipe models.Product

	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"id": recipe}).Decode(&recipe)
	if err != nil {
		return materials, err
	}

	return recipe.Materials, nil
}

// GetRecipeTree returns the recipe tree for a given recipe_id.
func (rs *RecipeService) GetRecipeTree(recipe_id string) (tree models.Product, err error) {

	self_materials := []models.Material{}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	var recipe models.Product

	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"id": recipe_id}).Decode(&recipe)
	if err != nil {
		rs.Logger.Error("GetRecipeTree@getting recipe" + err.Error())
		return tree, err
	}

	for _, material := range recipe.Materials {

		var db_component models.Material
		err = client.Database("waha").Collection("materials").FindOne(context.Background(), bson.M{"id": material.Id}).Decode(&db_component)
		if err != nil {
			return tree, err
		}

		valid_entries := []models.MaterialEntry{}
		for _, entry := range db_component.Entries {
			if entry.Quantity > 0 {
				valid_entries = append(valid_entries, entry)
			}
		}

		self_materials = append(self_materials, models.Material{
			Id:       material.Id,
			Name:     db_component.Name,
			Quantity: material.Quantity,
			Entries:  valid_entries,
			Unit:     db_component.Unit,
		})
	}

	for _, sub_product := range recipe.SubProducts {
		sub_recipe, err := rs.GetRecipeTree(sub_product.Id)
		if err != nil {
			return tree, err
		}

		//TODO - Fix Quantity vs Ready from Product
		sub_recipe.Quantity = float64(sub_product.Quantity)
		tree.SubProducts = append(tree.SubProducts, sub_recipe)
	}

	tree.Materials = self_materials
	tree.Id = recipe_id
	tree.Name = recipe.Name
	tree.Quantity = recipe.Quantity
	tree.Price = recipe.Price
	tree.Ready = recipe.Ready

	return tree, nil
}

// GetReadyNumber returns the ready number of a given recipe
func (rs *RecipeService) GetReadyNumber(recipe_id string) (ready float64, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	var product models.Product
	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"id": recipe_id}).Decode(&product)
	if err != nil {
		return ready, err
	}

	ready = product.Ready

	return ready, nil
}

// CheckRecipesAvailability checks the availability of a list of recipes.
// It returns a slice of RecipeAvailability with the available and ready number for each recipe.
func (rs *RecipeService) CheckRecipesAvailability(recipe_ids []string) (availabilities []dto.RecipeAvailability, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	coll := client.Database("waha").Collection("recipes")

	availabilitiesChan := make(chan dto.RecipeAvailability)
	errorChan := make(chan error)

	var wg sync.WaitGroup

	for _, recipe_id := range recipe_ids {

		wg.Add(1)

		go func(ch chan<- dto.RecipeAvailability, errorChan chan<- error, recipe_id string) {

			defer wg.Done()

			var recipe models.Product
			err = coll.FindOne(ctx, bson.M{"id": recipe_id}).Decode(&recipe)
			if err != nil {
				errorChan <- err
				return
			}

			var recipeAvailability dto.RecipeAvailability
			recipeAvailability.ComponentRequirements = make(map[string]float64)

			self_component_requirements := make(map[string]float64)
			materials_inventory := make(map[string]float64)

			var lowest_available float64

			// subrecipes_components_requirements := make(map[string]float64)
			// subrecipes_components_consumption := make(map[string]float64)

			subrecipe_availability := []dto.RecipeAvailability{}

			for _, material := range recipe.Materials {

				self_component_requirements[material.Id] = float64(material.Quantity)
				recipeAvailability.ComponentRequirements[material.Id] += self_component_requirements[material.Id]

				materialService := MaterialService{
					Logger: rs.Logger,
					Config: rs.Config,
				}

				component_amount, err := materialService.GetComponentAvailability(material.Id)
				if err != nil {
					errorChan <- err
					return
				}

				materials_inventory[material.Id] = float64(component_amount)
			}

			for _, product := range recipe.SubProducts {

				self_component_requirements[product.Id] = float64(product.Quantity)
				subrecipe_available, err := rs.CheckRecipesAvailability([]string{product.Id})

				if err != nil {
					errorChan <- err
					return
				}

				subrecipe_availability = append(subrecipe_availability, dto.RecipeAvailability{
					RecipeId:              product.Id,
					Ready:                 subrecipe_available[0].Ready,
					ComponentRequirements: subrecipe_available[0].ComponentRequirements,
				})
			}

			satisfied := false
			for _, sra := range subrecipe_availability {
				for k, v := range sra.ComponentRequirements {
					recipeAvailability.ComponentRequirements[k] += v * self_component_requirements[sra.RecipeId]
				}
			}

			for !satisfied {
				satisfied = true
				temp_component_requirements := make(map[string]float64)
				temp_component_inventory := make(map[string]float64)

				for k, v := range recipeAvailability.ComponentRequirements {
					temp_component_requirements[k] = v
				}

				for k, v := range materials_inventory {
					temp_component_inventory[k] = v
				}

				for index, subrecipe := range subrecipe_availability {
					if subrecipe.Ready > 0 {
						satisfied = false

						subrecipe_reminder := 0.0

						if self_component_requirements[subrecipe.RecipeId] > subrecipe.Ready {
							subrecipe_reminder = self_component_requirements[subrecipe.RecipeId] - subrecipe.Ready
						}

						for component_id, value := range subrecipe.ComponentRequirements {
							if _, ok := temp_component_requirements[component_id]; ok {
								temp_component_requirements[component_id] -= value * (self_component_requirements[subrecipe.RecipeId] - subrecipe_reminder)
							}
						}

						subrecipe_availability[index].Ready -= self_component_requirements[subrecipe.RecipeId] - subrecipe_reminder
					}
				}

				increase_availability := false
				for k, v := range temp_component_requirements {
					if temp_component_inventory[k] >= v {
						temp_component_inventory[k] -= v
						increase_availability = true
					}
				}

				if increase_availability && !satisfied {
					recipeAvailability.Available += 1
					for k, v := range temp_component_inventory {
						materials_inventory[k] = v
					}
				}

			}

			for index, ingredient := range recipe.Materials {

				if index == 0 {
					lowest_available = materials_inventory[ingredient.Id] / recipeAvailability.ComponentRequirements[ingredient.Id]
				}

				if materials_inventory[ingredient.Id]/recipeAvailability.ComponentRequirements[ingredient.Id] < lowest_available {
					lowest_available = materials_inventory[ingredient.Id] / recipeAvailability.ComponentRequirements[ingredient.Id]
				}
			}

			recipeAvailability.RecipeId = recipe_id
			recipeAvailability.Available += lowest_available + recipe.Ready
			recipeAvailability.Ready = recipe.Ready

			availabilitiesChan <- recipeAvailability

		}(availabilitiesChan, errorChan, recipe_id)

	}

	go func() {
		wg.Wait()
		close(availabilitiesChan)
	}()

	for {
		select {
		case availability, hasMore := <-availabilitiesChan:
			if !hasMore {
				return availabilities, err
			}

			availabilities = append(availabilities, availability)
		case err := <-errorChan:
			return availabilities, err
		}
	}
}
