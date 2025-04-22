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

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogService represents the logging service in the core module.
// It holds the logger and configuration required for logging operations.
type LogService struct {
	// Logger is used to log messages with different levels of severity.
	Logger logger.ILogger
	// Config holds the configuration settings for the logging service.
	Config config.Config
}

// GetMaterialLogs gets all logs for a given component_id.
func (l *LogService) GetMaterialLogs(component_id string, page_number, page_size int) (logs []models.ComponentConsumeLogs, total_records int64, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", l.Config.Databases[0].Host, l.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, total_records, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, total_records, err
	}

	// Connected successfully
	fmt.Println("Connected to MongoDB!")

	filter := bson.M{"type": bson.M{
		"$in": []string{models.LogTypeMaterialConsume, models.LogTypeMaterialInventoryReturn},
	}, "component_id": component_id}

	skip := (page_number - 1) * page_size
	if page_number == 1 {
		skip = 0
	}

	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(page_size)).SetSort(bson.M{"date": -1})

	collection := client.Database(l.Config.Databases[0].Database).Collection("logs")
	total_records, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, 0, err
	}

	cur, err := client.Database(l.Config.Databases[0].Database).Collection("logs").Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error(err.Error())
		return logs, total_records, err
	}
	defer cur.Close(ctx)

	if err = cur.All(ctx, &logs); err != nil {
		l.Logger.Error(err.Error())
		return logs, total_records, err
	}

	return logs, total_records, err
}

// GetSalesLogs gets all logs for a given component_id.
func (l *LogService) GetSalesLogs() []models.SalesLogs {
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

	// find all documents from db of logs collection filter on type = order_finish
	collection := client.Database(l.Config.Databases[0].Database).Collection("logs")
	filter := bson.M{"type": models.LogTypeOrderFinish}
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

func (l *LogService) GetOrderItemsRefundLogs(order_items_ids [][]string) (logs []models.OrderItemRefundLog, err error) {
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

	// find all documents from db of logs collection filter on type = order_finish
	collection := client.Database(l.Config.Databases[0].Database).Collection("logs")
	filter := bson.M{"type": models.LogTypeOrderItemRefunded}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		l.Logger.Error(err.Error())
	}
	defer cursor.Close(ctx)

	refund_logs := []models.OrderItemRefundLog{}
	if err = cursor.All(ctx, &refund_logs); err != nil {
		l.Logger.Error(err.Error())
	}

	return logs, err
}
