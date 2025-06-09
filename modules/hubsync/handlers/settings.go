package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PatchSettings(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Data models.Hubsync `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))

		deadline := 5 * time.Second
		if config.Env == "dev" {
			deadline = 1000 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Error(fmt.Sprintf("error in updating settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
		// connected to db

		collection := client.Database(config.Databases[0].Database).Collection("hubsync")
		_, err = collection.UpdateOne(ctx, bson.D{{"settings", bson.D{{"$exists", true}}}}, bson.D{{"$set", bson.D{{"settings", data.Data.Settings}}}})
		if err != nil {
			logger.Error(fmt.Sprintf("error in updating settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func GetSettings(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))

		deadline := 5 * time.Second
		if config.Env == "dev" {
			deadline = 1000 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Error(fmt.Sprintf("error in getting settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		collection := client.Database(config.Databases[0].Database).Collection("hubsync")
		var hubsync models.Hubsync
		err = collection.FindOne(ctx, bson.M{}).Decode(&hubsync)
		if err != nil {
			logger.Error(fmt.Sprintf("error in getting settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Data models.Hubsync `json:"data"`
		}{
			Data: hubsync,
		}

		response, err := json.Marshal(data)
		if err != nil {
			logger.Error(fmt.Sprintf("error marshalling settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
