package config

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConfigFactory(t string, path string, logger logger.ILogger) Config {
	switch t {
	case "viper":
		viper_config := NewViperConfig(logger)
		viper_config.ReadFile(path)
		config, err := viper_config.GetConfig()
		if err != nil {
			logger.Error("can't reat config")
		}
		return config
	}

	return Config{}
}

type IConfig interface {
	ReadFile(path string)
	GetConfig() (config Config)
}

type Settings struct {
	Inventory struct {
		DefaultInventoryQuantityWarn float64 `json:"default_inventory_quantity_warn" bson:"default_inventory_quantity_warn"`
	} `bson:"inventory" json:"inventory"`
}

func (s *Settings) LoadFromDB(config Config) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Get the "test" collection from the database
	collection := client.Database("waha").Collection("settings")
	err = collection.FindOne(ctx, bson.D{}).Decode(s)
	if err != nil {
		return err
	}

	return nil

}

type Config struct {
	Databases    []Database
	Env          string `mapstructure:"env"`
	JwtSecretKey string `mapstructure:"jwt_secret_key"`
	TimeZone     string `mapstructure:"timezone"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
}
