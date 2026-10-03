package common

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/FerretDB/FerretDB/ferretdb"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var singleDBInstance *mongo.Client
var ferretCancel context.CancelFunc
var ferretDone chan struct{}
var lock = &sync.Mutex{}

func GetDatabaseClient(logger logger.ILogger, conf *config.Config) (*mongo.Client, error) {

	if singleDBInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleDBInstance == nil {
			logger.Info("Creating DB single instance now.")

			dbConf := conf.Databases[0]

			var clientOptions *options.ClientOptions

			switch dbConf.Type {
			case "ferret":
				uri, err := startEmbeddedFerret(logger, &dbConf)
				if err != nil {
					return nil, err
				}
				clientOptions = options.Client().ApplyURI(uri)
			default:
				uri := dbConf.URI
				if uri == "" {
					uri = fmt.Sprintf("mongodb://%s:%v", dbConf.Host, dbConf.Port)
				}
				clientOptions = options.Client().ApplyURI(uri)
			}

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

			if dbConf.Type == "ferret" {
				if err := waitForFerretReady(ctx, client); err != nil {
					return nil, err
				}
			}

			singleDBInstance = client
		}
	}

	return singleDBInstance, nil
}

// startEmbeddedFerret starts an in-process FerretDB instance backed by a
// local SQLite file. It returns the MongoDB URI the driver should connect to.
// No external database service is required.
func startEmbeddedFerret(logger logger.ILogger, dbConf *config.Database) (string, error) {
	logger.Info("Starting embedded FerretDB (sqlite) instance.")

	filePath := dbConf.FilePath
	if filePath == "" {
		filePath = "./data/db"
	}

	if err := os.MkdirAll(filePath, 0o755); err != nil {
		return "", fmt.Errorf("failed to create database directory: %w", err)
	}

	f, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: "127.0.0.1:0",
		},
		Handler:   "sqlite",
		SQLiteURL: "file:" + filePath + "/",
	})
	if err != nil {
		return "", fmt.Errorf("failed to construct FerretDB instance: %w", err)
	}

	ferretCtx, cancel := context.WithCancel(context.Background())
	ferretCancel = cancel
	ferretDone = make(chan struct{})

	go func() {
		defer close(ferretDone)
		if err := f.Run(ferretCtx); err != nil && ferretCtx.Err() == nil {
			logger.Error("embedded ferretdb run error: " + err.Error())
		}
	}()

	return f.MongoDBURI(), nil
}

// CloseDatabase gracefully shuts down the database layer: it disconnects the
// mongo client first, then stops the embedded FerretDB instance (if any) so
// its SQLite backend is flushed cleanly.
func CloseDatabase() {
	lock.Lock()
	defer lock.Unlock()

	if singleDBInstance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = singleDBInstance.Disconnect(ctx)
		singleDBInstance = nil
	}

	if ferretCancel != nil {
		done := ferretDone
		ferretCancel()
		ferretCancel = nil

		if done != nil {
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
		}
	}
	ferretDone = nil

}

// waitForFerretReady pings the embedded instance until it answers or the
// context deadline expires.
func waitForFerretReady(ctx context.Context, client *mongo.Client) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err := client.Ping(pingCtx, nil)
		cancel()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("embedded ferretdb did not become ready: %w", err)
		case <-ticker.C:
		}
	}
}
