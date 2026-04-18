package common

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var singleDBInstance *mongo.Client
var lock = &sync.Mutex{}

func GetDatabaseClient(logger logger.ILogger, conf *config.Config) (*mongo.Client, error) {

	if singleDBInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleDBInstance == nil {
			logger.Info("Creating DB single instance now.")
			clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", conf.Databases[0].Host, conf.Databases[0].Port))

			deadline := 5 * time.Second
			if conf.Env == "dev" {
				deadline = 1000 * time.Second
			}

			ctx, cancel := context.WithTimeout(context.Background(), deadline)
			defer cancel()

			client, err := mongo.Connect(ctx, clientOptions)
			if err != nil {
				return nil, err
			}
			singleDBInstance = client
		}
	}

	return singleDBInstance, nil
}
