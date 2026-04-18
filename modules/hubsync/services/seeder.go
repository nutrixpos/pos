package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
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
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx := context.Background()

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
			Settings: models.Settings{
				Enabled:      false,
				BufferSize:   100,
				SyncInterval: 60,
			},
			LastSynced:   time.Now().Unix(),
			SyncProgress: 0,
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
				SyncProgress: 0,
				Settings: models.Settings{
					SyncInterval: 60,
					BufferSize:   100,
				},
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
