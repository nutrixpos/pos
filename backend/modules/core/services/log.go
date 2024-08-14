package services

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Log struct {
	Logger logger.ILogger
	Config config.Config
}

func (l *Log) GetComponentLogs(name string) (logs []models.ComponentConsumeLogs, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", l.Config.Databases[0].Host, l.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	filter := bson.M{"type": "component_consume", "name": name}
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

func (l *Log) GetSalesLogs() []models.SalesLogs {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", l.Config.Databases[0].Host, l.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
