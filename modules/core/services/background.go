// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckExpirationDates is a background job that checks all materials if they are expired
// and informs the admin about it. The function is designed to be called periodically
// by the job scheduler.
func CheckExpirationDates(log logger.ILogger, conf config.Config, notification_svc INotificationService) {

	log.Info("core:background: Checking expiration dates")

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", conf.Databases[0].Host, conf.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Connected successfully

	// Get the "test" collection from the database
	collection := client.Database(conf.Databases[0].Host).Collection("materials")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer cursor.Close(ctx)

	// Iterate over the documents and print them as JSON
	for cursor.Next(ctx) {
		var component models.Material
		err := cursor.Decode(&component)
		if err != nil {
			log.Error(err.Error())
			return
		}

		for _, entry := range component.Entries {
			t := time.Until(entry.ExpirationDate)
			if t <= 14*24*time.Hour {

				msg := fmt.Sprintf("Material %s, entry %s will expire within 2 weeks", component.Name, entry.Id)

				log.Warning(msg)

				topic_msg := &models.WebsocketTopicServerMessage{
					Type:      "topic_message",
					TopicName: "expire_soon",
					Message:   msg,
					Severity:  "warn",
					Date:      time.Now(),
				}

				jsonstr, err := json.Marshal(topic_msg)
				if err != nil {
					log.Error(err.Error())
					return
				}

				notification_svc.SendToTopic("expire_soon", string(jsonstr))
			}
		}
	}
}
