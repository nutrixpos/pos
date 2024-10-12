package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecipeService struct {
	Logger logger.ILogger
	Config config.Config
}

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

	return tree, nil
}

//TODO - Continue

func (rs *RecipeService) ConsumeRecipe(recipes_id string) (err error) {

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

	return nil
}

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
			recipeAvailability.Available += lowest_available + recipe.Quantity
			recipeAvailability.Ready = recipe.Quantity

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
