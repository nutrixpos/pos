package services

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SettingsService struct {
	Config config.Config
}

// UpdateSettings updates the settings in the database
func (ss *SettingsService) UpdateSettings(settings models.Settings) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	collection := client.Database(ss.Config.Databases[0].Database).Collection("settings")
	_, err = collection.UpdateOne(ctx, bson.M{}, bson.M{"$set": settings})

	return
}

// GetSettings returns the settings from the database
func (os *SettingsService) GetSettings() (ordersettings models.Settings, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", os.Config.Databases[0].Host, os.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if os.Config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	var settings models.Settings
	err = client.Database(os.Config.Databases[0].Database).Collection("settings").FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		return
	}

	return settings, err
}
