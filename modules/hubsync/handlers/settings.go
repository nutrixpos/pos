package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
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

		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			logger.Error(fmt.Sprintf("error in updating settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

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
		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			logger.Error(fmt.Sprintf("error in getting settings: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

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
