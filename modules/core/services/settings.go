package services

import (
	"context"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SettingsService struct {
	Config config.Config
	Logger logger.ILogger
}

func (ss *SettingsService) UpdateSettings(settings models.Settings) (err error) {
	client, err := common.GetDatabaseClient(ss.Logger, &ss.Config)
	if err != nil {
		return
	}

	ctx := context.Background()

	collection := client.Database(ss.Config.Databases[0].Database).Collection("settings")
	_, err = collection.UpdateOne(ctx, bson.M{}, bson.M{"$set": settings})

	return
}

func (os *SettingsService) GetSettings() (ordersettings models.Settings, err error) {
	client, err := common.GetDatabaseClient(os.Logger, &os.Config)
	if err != nil {
		return
	}

	ctx := context.Background()

	var settings models.Settings
	err = client.Database(os.Config.Databases[0].Database).Collection("settings").FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		return
	}

	return settings, err
}