package services

import (
	"context"
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SettingsSvc struct {
	Config config.Config
	Logger logger.ILogger
}

func (s *SettingsSvc) Get() (settings models.Hubsync, err error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return
	}

	ctx := context.Background()

	collection := client.Database(s.Config.Databases[0].Database).Collection("hubsync")
	err = collection.FindOne(ctx, bson.D{{}}).Decode(&settings)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error in getting settings: %v", err))
		return settings, err
	}

	return
}
