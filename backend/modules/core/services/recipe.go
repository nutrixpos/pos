package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
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

func (rs *RecipeService) CheckRecipesAvailability(recipe_ids []string) (availabilities map[string]int, err error) {

	availabilities = make(map[string]int)

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

	availabilitiesChan := make(chan map[string]interface{})
	errorChan := make(chan error)

	var wg sync.WaitGroup

	for _, recipe_id := range recipe_ids {

		wg.Add(1)

		go func(ch chan<- map[string]interface{}, errorChan chan<- error, recipe_id string) {

			defer wg.Done()

			var recipe models.Recipe
			objID, err := primitive.ObjectIDFromHex(recipe_id)
			if err != nil {
				errorChan <- err
				close(ch)
			}
			err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe)
			if err != nil {
				errorChan <- err
				close(ch)
			}

			lowest_available := 0

			for index, component := range recipe.Components {
				componentService := ComponentService{
					Logger: rs.Logger,
					Config: rs.Config,
				}

				component_amount, err := componentService.GetComponentAvailability(component.ComponentId)
				if err != nil {
					errorChan <- err
					close(ch)
					break
				}

				if index == 0 {
					lowest_available = int(component_amount / component.Quantity)
				}

				if int(component_amount/component.Quantity) < lowest_available {
					lowest_available = int(component_amount / component.Quantity)
				}
			}

			recipeAvailability := make(map[string]interface{})
			recipeAvailability["id"] = recipe_id
			recipeAvailability["available"] = lowest_available

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
			recipe_id := availability["id"].(string)
			availabilities[recipe_id] = availability["available"].(int)
		case err := <-errorChan:
			return availabilities, err
		}
	}
}
