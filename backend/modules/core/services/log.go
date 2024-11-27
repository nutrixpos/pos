// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the database and
// external services. They are used to implement the HTTP handlers in the
// handlers package.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
	"github.com/elmawardy/nutrix/backend/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Log represents the logging service in the core module.
// It holds the logger and configuration required for logging operations.
type Log struct {
	// Logger is used to log messages with different levels of severity.
	Logger logger.ILogger
	// Config holds the configuration settings for the logging service.
	Config config.Config
}

// GetComponentLogs gets all logs for a given component_id.
func (l *Log) GetComponentLogs(component_id string) (logs []models.ComponentConsumeLogs, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", l.Config.Databases[0].Host, l.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, err
	}

	// Connected successfully
	fmt.Println("Connected to MongoDB!")

	filter := bson.M{"type": "component_consume", "component_id": component_id}
	cur, err := client.Database("waha").Collection("logs").Find(ctx, filter)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, err
	}
	defer cur.Close(ctx)

	if err = cur.All(ctx, &logs); err != nil {
		l.Logger.Error(err.Error())
		return logs, err
	}

	return logs, nil
}

// GetSalesLogs gets all logs for a given component_id.
func (l *Log) GetSalesLogs() []models.SalesLogs {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", l.Config.Databases[0].Host, l.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.Logger.Error(err.Error())
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		l.Logger.Error(err.Error())
	}

	// Connected successfully
	fmt.Println("Connected to MongoDB!")

	// find all documents from db of logs collection filter on type = order_finished
	collection := client.Database("waha").Collection("logs")
	filter := bson.M{"type": "order_finish"}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		l.Logger.Error(err.Error())
	}
	defer cursor.Close(ctx)

	sales_logs := []models.SalesLogs{}
	if err = cursor.All(ctx, &sales_logs); err != nil {
		l.Logger.Error(err.Error())
	}

	return sales_logs
}
