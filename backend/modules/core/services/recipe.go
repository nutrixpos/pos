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

			var lowest_available float64

			for index, component := range recipe.Components {
				componentService := ComponentService{
					Logger: rs.Logger,
					Config: rs.Config,
				}

				component_amount, err := componentService.GetComponentAvailability(component.ComponentId)
				if err != nil {
					errorChan <- err
					return
				}

				if index == 0 {
					lowest_available = float64(component_amount / component.Quantity)
				}

				if float64(component_amount/component.Quantity) < lowest_available {
					lowest_available = float64(component_amount / component.Quantity)
				}
			}

			var recipeAvailability dto.RecipeAvailability
			recipeAvailability.RecipeId = recipe_id
			recipeAvailability.Available = lowest_available

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

			availabilities = append(availabilities, dto.RecipeAvailability{
				RecipeId:  availability.RecipeId,
				Available: availability.Available,
				Ready:     availability.Ready,
			})
		case err := <-errorChan:
			return availabilities, err
		}
	}
}
