package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeederService struct {
	Logger    logger.ILogger
	Config    config.Config
	Prompter  userio.Prompter
	IsNewOnly bool
}

// GetSeedables returns a list of seedables.
func (c *SeederService) GetSeedables() (entities []string, err error) {
	c.Logger.Info("Getting seedables...")

	return []string{
		"hubsync",
	}, nil
}

func (s *SeederService) Seed() error {
	err := s.seedHubsyncCollection()
	if err != nil {
		return err
	}

	return nil
}

func (s *SeederService) seedHubsyncCollection() error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// connected to db

	// Check if the "hubsync" collection exists
	db := client.Database(s.Config.Databases[0].Database)
	collectionNames, err := db.ListCollectionNames(ctx, bson.M{"name": "hubsync"})
	if err != nil {
		return err
	}

	if len(collectionNames) == 0 {

		err = db.CreateCollection(ctx, "hubsync")
		if err != nil {
			return err
		}

		// Insert a simple document into the "hubsync" collection
		hubsyncCollection := db.Collection("hubsync")
		_, err = hubsyncCollection.InsertOne(ctx, models.Hubsync{
			LastSynced:   time.Now().Unix(),
			SyncInterval: 60,
			SyncProgress: 0,
			BufferSize:   100,
		})
		if err != nil {
			return err
		}
	} else {
		hubsyncCollection := db.Collection("hubsync")
		cursor, err := hubsyncCollection.Find(ctx, bson.M{})
		if err != nil {
			return err
		}
		var tracker models.Hubsync
		if !cursor.Next(ctx) {
			_, err = hubsyncCollection.InsertOne(ctx, models.Hubsync{
				LastSynced:   0,
				SyncInterval: 60,
				SyncProgress: 0,
				BufferSize:   100,
			})
			if err != nil {
				return err
			}
		} else {
			err = cursor.Decode(&tracker)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
