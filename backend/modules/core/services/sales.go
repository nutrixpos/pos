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

type SalesService struct {
	Logger logger.ILogger
	Config config.Config
}

// format 2006-01-02
func (ss *SalesService) GetSalesPerday(first int, rows int) (salesPerDay []models.SalesPerDay, totalRecords int, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		ss.Logger.Error(err.Error())
		return salesPerDay, totalRecords, err
	}

	collection := client.Database("waha").Collection("sales")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"date": -1})
	findOptions.SetSkip(int64(first))
	findOptions.SetLimit(int64(rows))
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		ss.Logger.Error(err.Error())
		return salesPerDay, 0, err
	}
	totalRecords = int(count)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		ss.Logger.Error(err.Error())
		return salesPerDay, totalRecords, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var spd models.SalesPerDay
		if err := cursor.Decode(&spd); err != nil {
			return salesPerDay, totalRecords, err
		}

		salesPerDay = append(salesPerDay, spd)
	}

	return salesPerDay, totalRecords, nil
}

func (ss *SalesService) AddOrderToSalesDay(order models.Order, items_cost []models.ItemCost) error {

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

	sales_order := models.SalesPerDayOrder{
		Order: order,
		Costs: items_cost,
	}

	collection := client.Database("waha").Collection("sales")
	filter := bson.M{"date": time.Now().Format("2006-01-02")}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = collection.InsertOne(ctx, bson.M{"date": time.Now().Format("2006-01-02"), "orders": []models.SalesPerDayOrder{sales_order}, "costs": sales_order.Order.Cost, "total_sales": sales_order.Order.SalePrice})
		if err != nil {
			return err
		}
	} else {
		_, err = collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"orders": sales_order}, "$inc": bson.M{"costs": sales_order.Order.Cost, "total_sales": sales_order.Order.SalePrice}}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}
