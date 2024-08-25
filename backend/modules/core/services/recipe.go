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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecipeService struct {
	Logger logger.ILogger
	Config config.Config
}

func (rs *RecipeService) GetRecipeComponents(recipe_id string) (components []models.RecipeComponent, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	var recipe models.Recipe

	subRecipeID, err := primitive.ObjectIDFromHex(recipe_id)
	if err != nil {
		return components, err
	}

	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"_id": subRecipeID}).Decode(&recipe)
	if err != nil {
		return components, err
	}

	return recipe.Components, nil
}

func (rs *RecipeService) GetRecipeTree(recipe_id string) (tree dto.RecipeTree, err error) {

	self_components := []dto.RecipeComponentResponse{}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	var recipe models.Recipe

	subRecipeID, err := primitive.ObjectIDFromHex(recipe_id)
	if err != nil {
		return tree, err
	}

	err = client.Database("waha").Collection("recipes").FindOne(context.Background(), bson.M{"_id": subRecipeID}).Decode(&recipe)
	if err != nil {
		return tree, err
	}

	for _, component := range recipe.Components {

		if component.Type != "recipe" {
			component_id, err := primitive.ObjectIDFromHex(component.ComponentId)
			if err != nil {
				return tree, err
			}

			var db_component models.Component
			err = client.Database("waha").Collection("components").FindOne(context.Background(), bson.M{"_id": component_id}).Decode(&db_component)
			if err != nil {
				return tree, err
			}
			self_components = append(self_components, dto.RecipeComponentResponse{
				ComponentId:     component.ComponentId,
				Name:            db_component.Name,
				DefaultQuantity: component.Quantity,
				Unit:            db_component.Unit,
				Entries:         db_component.Entries,
				Type:            component.Type,
			})

		} else {

			sub_recipe_tree, err := rs.GetRecipeTree(component.ComponentId)
			if err != nil {
				return tree, err
			}
			tree.SubRecipes = append(tree.SubRecipes, sub_recipe_tree)
		}

	}

	tree.Components = self_components
	tree.RecipeId = recipe_id
	tree.RecipeName = recipe.Name

	return tree, nil
}

//TODO - Continue

func (rs *RecipeService) ConsumeRecipe(recipes_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	return nil
}

func (rs *RecipeService) CheckRecipesAvailability(recipe_ids []string) (availabilities []dto.RecipeAvailability, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

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
	rs.Logger.Info("Connected to MongoDB!")

	coll := client.Database("waha").Collection("recipes")

	availabilitiesChan := make(chan dto.RecipeAvailability)
	errorChan := make(chan error)

	var wg sync.WaitGroup

	for _, recipe_id := range recipe_ids {

		wg.Add(1)

		go func(ch chan<- dto.RecipeAvailability, errorChan chan<- error, recipe_id string) {

			defer wg.Done()

			var recipe models.Recipe
			objID, err := primitive.ObjectIDFromHex(recipe_id)
			if err != nil {
				errorChan <- err
				return
			}
			err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe)
			if err != nil {
				errorChan <- err
				return
			}

			var recipeAvailability dto.RecipeAvailability
			recipeAvailability.ComponentRequirements = make(map[string]float64)

			self_component_requirements := make(map[string]float64)
			components_inventory := make(map[string]float64)

			var lowest_available float64

			// subrecipes_components_requirements := make(map[string]float64)
			// subrecipes_components_consumption := make(map[string]float64)

			subrecipe_availability := []dto.RecipeAvailability{}

			for _, component := range recipe.Components {

				self_component_requirements[component.ComponentId] = float64(component.Quantity)
				recipeAvailability.ComponentRequirements[component.ComponentId] += self_component_requirements[component.ComponentId]

				if component.Type == "recipe" {
					// get sub recipes availability and component consumption to reduce it from the available component quantities.

					subrecipe_available, err := rs.CheckRecipesAvailability([]string{component.ComponentId})

					if err != nil {
						errorChan <- err
						return
					}

					subrecipe_availability = append(subrecipe_availability, dto.RecipeAvailability{
						RecipeId:              component.ComponentId,
						Ready:                 subrecipe_available[0].Ready,
						ComponentRequirements: subrecipe_available[0].ComponentRequirements,
					})

				} else {
					componentService := ComponentService{
						Logger: rs.Logger,
						Config: rs.Config,
					}

					component_amount, err := componentService.GetComponentAvailability(component.ComponentId)
					if err != nil {
						errorChan <- err
						return
					}

					components_inventory[component.ComponentId] = float64(component_amount)
				}
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

				for k, v := range components_inventory {
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
						components_inventory[k] = v
					}
				}

			}

			for index, component := range recipe.Components {

				if component.Type != "recipe" {
					// subrecipes_components_requirements[component.ComponentId] += float64(component.Quantity)

					if index == 0 {
						lowest_available = components_inventory[component.ComponentId] / recipeAvailability.ComponentRequirements[component.ComponentId]
					}

					if components_inventory[component.ComponentId]/recipeAvailability.ComponentRequirements[component.ComponentId] < lowest_available {
						lowest_available = components_inventory[component.ComponentId] / recipeAvailability.ComponentRequirements[component.ComponentId]
					}
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
