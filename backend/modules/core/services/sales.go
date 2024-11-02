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

type SalesService struct {
	Logger logger.ILogger
	Config config.Config
}

func (ss *SalesService) AddOrderToSalesDay(order models.Order, totalcost float64, totalsale float64) error {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
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

	// Connected successfully

	collection := client.Database("waha").Collection("sales")
	filter := bson.M{"date": time.Now().Format("2006-01-02")}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = collection.InsertOne(ctx, bson.M{"date": time.Now().Format("2006-01-02"), "orders": []models.Order{order}, "costs": totalcost, "total_sales": totalsale})
		if err != nil {
			return err
		}
	} else {
		_, err = collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"orders": order}, "$inc": bson.M{"costs": totalcost, "total_sales": totalsale}}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}
