// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the persistence layer
// and perform operations on the data models of the core module.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/dto"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SalesService contains the configuration and logger for the sales service.
type SalesService struct {
	// Logger is the logger for the sales service.
	Logger logger.ILogger
	// Config is the configuration for the sales service.
	Config config.Config
}

// format 2006-01-02
// GetSalesPerday returns a slice of models.SalesPerDay and the total count of records in the database.
// It takes two parameters, first and rows, which determine the offset and limit of the query.
// It returns an error if the query fails.
func (ss *SalesService) GetSalesPerday(page_number int, page_size int) (salesPerDay []models.SalesPerDay, totalRecords int, err error) {
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

	collection := client.Database(ss.Config.Databases[0].Database).Collection("sales")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"date": -1})

	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"date": 1})
	findOptions.SetLimit(int64(page_size))
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

		if spd.Refunds == nil {
			spd.Refunds = make([]models.SalesPerDayRefund, 0)
		}

		salesPerDay = append(salesPerDay, spd)
	}

	return salesPerDay, totalRecords, nil
}

func (ss *SalesService) AddOrderItemToDayRefund(refund_request dto.OrderItemRefundRequest) error {

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

	material_refunds := make([]models.OrderItemRefundMaterial, 0)
	products_adds := make([]models.OrderItemRefundProductAdd, 0)

	for _, material_refund := range refund_request.MaterialRerunds {
		material_refunds = append(material_refunds, models.OrderItemRefundMaterial{
			MaterialId:         material_refund.MaterialId,
			EntryId:            material_refund.EntryId,
			InventoryReturnQty: material_refund.InventoryReturnQty,
			DisposeQty:         material_refund.DisposeQty,
			WasteQty:           material_refund.WasteQty,
			Comment:            material_refund.Comment,
		})
	}

	for _, product_add := range refund_request.ProductAdd {
		products_adds = append(products_adds, models.OrderItemRefundProductAdd{
			ProductId: product_add.ProductId,
			Quantity:  product_add.Quantity,
			Comment:   product_add.Comment,
		})
	}

	sales_refund := models.SalesPerDayRefund{
		OrderId:         refund_request.OrderId,
		ItemId:          refund_request.ItemId,
		Reason:          refund_request.Reason,
		Amount:          refund_request.RefundValue,
		ProductId:       refund_request.OrderId,
		Destination:     refund_request.Destination,
		MaterialRerunds: material_refunds,
		ProductAdd:      products_adds,
	}

	collection := client.Database(ss.Config.Databases[0].Database).Collection("sales")
	filter := bson.M{"date": time.Now().Format("2006-01-02")}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = collection.InsertOne(ctx, bson.M{"date": time.Now().Format("2006-01-02"), "refunds": []models.SalesPerDayRefund{sales_refund}, "orders": []bson.M{}, "refunds_value": refund_request.RefundValue})
		if err != nil {
			return err
		}
	} else {
		_, err = collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"refunds": sales_refund}, "$inc": bson.M{"refunds_value": refund_request.RefundValue}}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}

// AddOrderToSalesDay adds an order to the "sales" collection in the database.
// It takes two parameters, order and items_cost, which are the order and its associated item costs.
// It returns an error if the query fails.
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

	collection := client.Database(ss.Config.Databases[0].Database).Collection("sales")
	filter := bson.M{"date": time.Now().Format("2006-01-02")}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = collection.InsertOne(ctx, bson.M{"date": time.Now().Format("2006-01-02"), "refunds": []bson.M{}, "orders": []models.SalesPerDayOrder{sales_order}, "costs": sales_order.Order.Cost, "total_sales": sales_order.Order.SalePrice})
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
