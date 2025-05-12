package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	core_models "github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SyncerService struct {
	Config config.Config
	Logger logger.ILogger
}

func (s *SyncerService) CopyToBuffer() error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// connected to db

	db := client.Database(s.Config.Databases[0].Database)
	collection := db.Collection("hubsync")
	var tracker models.Hubsync

	err = collection.FindOne(ctx, bson.M{}).Decode(&tracker)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	lastSyncedTime := time.Unix(tracker.LastSynced, 0)

	cursor, err := db.Collection("logs").Find(ctx, bson.M{
		"$and": []bson.M{
			{"date": bson.M{"$gt": lastSyncedTime}},
			{"$or": []bson.M{
				{"hubsync_status": bson.M{"$exists": false}},
				{"hubsync_status": bson.M{"$eq": false}},
			}},
			{"type": bson.M{"$in": []string{
				core_models.LogTypeMaterialAdd,
				core_models.LogTypeMaterialConsume,
				core_models.LogTypeMaterialInventoryReturn,
				core_models.LogTypeMaterialWaste,
				core_models.LogTypeOrderStart,
				core_models.LogTypeOrderFinish,
				core_models.LogTypeOrderItemRefunded,
			}}},
		},
	})
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var log interface{}
		err = cursor.Decode(&log)
		if err != nil {
			return err
		}

		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			s.Logger.Error("Error decoding raw document: %v", err)
			continue
		}

		logType, ok := raw["type"].(string)
		if !ok {
			s.Logger.Error("Document missing type field")
			continue
		}

		switch logType {
		case core_models.LogTypeOrderStart:
			var db_log core_models.OrderStartLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeOrderFinish:
			var db_log core_models.OrderFinishLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeMaterialConsume:
			var db_log core_models.MaterialConsumeLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeMaterialAdd:
			var db_log core_models.MaterialAddLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeMaterialInventoryReturn:
			var db_log core_models.MaterialInventoryReturnLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeMaterialWaste:
			var db_log core_models.WasteMaterialLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
		case core_models.LogTypeOrderItemRefunded:
			var db_log core_models.OrderItemRefundLog
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}

		default:
			s.Logger.Warning(fmt.Sprintf("Unknown log type: %s", logType))
		}
	}

	tracker_id, _ := primitive.ObjectIDFromHex(tracker.Id)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": tracker_id}, bson.M{
		"$set": bson.M{
			"last_synced":   currentTime.Unix(),
			"sync_progress": 100,
			"sync_status":   "success",
		},
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return nil
}

func (s *SyncerService) Sync() error {
	s.CopyToBuffer()

	return nil
}
