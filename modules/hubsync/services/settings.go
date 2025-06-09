package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SettingsSvc struct {
	Config config.Config
	Logger logger.ILogger
}

func (s *SettingsSvc) Get() (settings models.Hubsync, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	// connected to db

	collection := client.Database(s.Config.Databases[0].Database).Collection("hubsync")
	err = collection.FindOne(ctx, bson.D{{}}).Decode(&settings)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error in getting settings: %v", err))
		return settings, err
	}

	return
}
