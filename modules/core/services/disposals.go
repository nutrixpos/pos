package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DisposalService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

func (rs *DisposalService) GetDisposal(disposal_id string) (disposal interface{}, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", rs.Config.Databases[0].Host, rs.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if rs.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return disposal, err
	}
	// connected to db

	collection := client.Database(rs.Config.Databases[0].Database).Collection("disposals")
	err = collection.FindOne(ctx, bson.M{"id": disposal_id}).Decode(&disposal)
	if err != nil {
		return disposal, err
	}

	return disposal, nil
}

func (ds *DisposalService) GetDisposals(params GetDisposalsParameters) (disposals []interface{}, totalRecords int64, err error) {

	disposals = make([]interface{}, 0)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ds.Config.Databases[0].Host, ds.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return disposals, 0, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return disposals, 0, err
	}

	// Connected successfully
	ds.Logger.Info("Connected to MongoDB!")
	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetLimit(int64(params.PageSize))
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))

	if params.DisposalIdContains != "" {
		filter["id"] = bson.M{
			"$regex": fmt.Sprintf("(?i).*%s.*", params.DisposalIdContains),
		}
	}

	positiveStateFilters := []string{}
	negativeStateFilter := []string{}

	for _, state := range params.FilterState {
		if state[0] == '!' {
			negativeStateFilter = append(negativeStateFilter, state[1:])
		} else {
			positiveStateFilters = append(positiveStateFilters, state)
		}
	}

	stateFilters := []bson.M{}

	if len(positiveStateFilters) > 0 {
		stateFilters = append(stateFilters, bson.M{"state": bson.M{"$in": positiveStateFilters}})
	}

	if len(negativeStateFilter) > 0 {
		stateFilters = append(stateFilters, bson.M{"state": bson.M{"$nin": negativeStateFilter}})
	}

	if len(stateFilters) > 0 {
		filter["$and"] = stateFilters
	}

	totalRecords, err = client.Database(ds.Config.Databases[0].Database).Collection("disposals").CountDocuments(ctx, bson.D{})
	if err != nil {
		ds.Logger.Error(err.Error())
		return
	}

	cursor, err := client.Database(ds.Config.Databases[0].Database).Collection("disposals").Find(context.Background(), filter, findOptions)
	if err != nil {
		return disposals, 0, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var disposal interface{}
		if err := cursor.Decode(&disposal); err != nil {
			return disposals, 0, err
		}

		disposals = append(disposals, disposal)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return disposals, 0, err
	}

	return disposals, totalRecords, nil
}

// UpdateDisposal updates a disposal in the database.
func (cs *DisposalService) UpdateDisposal(id string, disposal interface{}) (updatedDisposal interface{}, err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return updatedDisposal, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return updatedDisposal, err
	}

	// Connected successfully

	data := bson.M{}

	var idisposal interface{}
	idisposal = disposal

	if material_disposal, ok := idisposal.(models.MaterialDisposal); ok && material_disposal.Type == models.TypeDisposalMaterial {

		data["entry_id"] = material_disposal.EntryId
		data["material_id"] = material_disposal.MaterialId
		data["quantity"] = material_disposal.Quantity
		data["type"] = material_disposal.Type
		data["order_id"] = material_disposal.OrderId

		collection := client.Database(cs.Config.Databases[0].Database).Collection("disposals")
		_, err = collection.UpdateOne(ctx, bson.M{"id": material_disposal.Id}, bson.M{"$set": data})

		var db_material_disposal models.MaterialDisposal
		err = client.Database(cs.Config.Databases[0].Database).Collection("settings").FindOne(ctx, bson.M{}).Decode(&db_material_disposal)
		if err != nil {
			return db_material_disposal, err
		}
	}

	if product_disposal, ok := idisposal.(models.ProductDisposal); ok && product_disposal.Type == models.TypeDisposalProduct {

		data["item"] = product_disposal.Item
		data["quantity"] = product_disposal.Quantity
		data["type"] = product_disposal.Type
		data["order_id"] = product_disposal.OrderId

		collection := client.Database(cs.Config.Databases[0].Database).Collection("disposals")
		_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": data})

		var db_orderitem_disposal models.ProductDisposal
		err = client.Database(cs.Config.Databases[0].Database).Collection("settings").FindOne(ctx, bson.M{}).Decode(&db_orderitem_disposal)
		if err != nil {
			return db_orderitem_disposal, err
		}
	}

	return updatedDisposal, err
}

// DeleteDisposal deletes a category from the database.
func (cs *DisposalService) DeleteDisposal(disposal_id string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", cs.Config.Databases[0].Host, cs.Config.Databases[0].Port))

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
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

	collection := client.Database(cs.Config.Databases[0].Database).Collection("disposals")
	_, err = collection.DeleteOne(ctx, bson.M{"id": disposal_id})

	return err
}

func (ds *DisposalService) AddMaterialDisposal(disposal models.MaterialDisposal) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ds.Config.Databases[0].Host, ds.Config.Databases[0].Port))

	timeout := 1000 * time.Second

	if ds.Config.Env == "dev" {
		timeout = 5 * time.Minute
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	disposal.Id = primitive.NewObjectID().Hex()

	// Insert the DBComponent struct into the "materials" collection
	collection := client.Database(ds.Config.Databases[0].Database).Collection("disposals")
	_, err = collection.InsertOne(ctx, disposal)
	if err != nil {
		ds.Logger.Error(err.Error())
		return err
	}

	disposal_add_log := models.DisposalMaterialAddLog{
		Log: models.Log{
			Id:   primitive.NewObjectID().Hex(),
			Type: models.LogTypeDisposalAdd,
			Date: time.Now(),
		},
		Disposal: disposal,
	}

	_, err = client.Database(ds.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, disposal_add_log)
	if err != nil {
		ds.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (ds *DisposalService) AddProductDisposal(disposal models.ProductDisposal) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ds.Config.Databases[0].Host, ds.Config.Databases[0].Port))

	timeout := 1000 * time.Second

	if ds.Config.Env == "dev" {
		timeout = 5 * time.Minute
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	disposal.Id = primitive.NewObjectID().Hex()

	// Insert the DBComponent struct into the "materials" collection
	collection := client.Database(ds.Config.Databases[0].Database).Collection("disposals")
	_, err = collection.InsertOne(ctx, disposal)
	if err != nil {
		ds.Logger.Error(err.Error())
		return err
	}

	disposal_add_log := models.DisposalProductAddLog{
		Log: models.Log{
			Id:   primitive.NewObjectID().Hex(),
			Type: models.LogTypeDisposalAdd,
			Date: time.Now(),
		},
		Disposal: disposal,
	}

	_, err = client.Database(ds.Config.Databases[0].Database).Collection("logs").InsertOne(ctx, disposal_add_log)
	if err != nil {
		ds.Logger.Error(err.Error())
		return err
	}

	return nil
}
