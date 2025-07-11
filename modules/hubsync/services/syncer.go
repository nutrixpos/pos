package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	Info   models.Hubsync
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
				core_models.LogTypeSalesPerDayOrder,
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
			var db_log core_models.LogOrderStart
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}

			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}

		case core_models.LogTypeOrderFinish:
			var db_log core_models.LogOrderFinish
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}

			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeMaterialConsume:
			var db_log core_models.LogMaterialConsume
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}

			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeMaterialAdd:
			var db_log core_models.LogMaterialAdd
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}

			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeMaterialInventoryReturn:
			var db_log core_models.LogMaterialInventoryReturn
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeMaterialWaste:
			var db_log core_models.LogWasteMaterial
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeSalesPerDayOrder:
			var db_log core_models.LogSalesPerDayOrder
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
			}
		case core_models.LogTypeOrderItemRefunded:
			var db_log core_models.LogOrderItemRefund
			if err := cursor.Decode(&db_log); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			_, err = db.Collection("hubsync_buffer").InsertOne(ctx, db_log)
			if err != nil {
				return fmt.Errorf("error inserting log into hubsync_buffer: %v", err)
			}
			_, err = db.Collection("logs").UpdateOne(ctx, bson.M{"id": db_log.Id}, bson.M{
				"$set": bson.M{
					"hubsync_status": true,
				},
			})
			if err != nil {
				return fmt.Errorf("error updating log in logs collection: %v", err)
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

func (s *SyncerService) UploadToServer(host string) error {
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

	collection := client.Database(s.Config.Databases[0].Database).Collection("hubsync_buffer")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	db_logs := make([]interface{}, 0)

	for cursor.Next(ctx) {
		var db_log core_models.Log
		if err := cursor.Decode(&db_log); err != nil {
			return fmt.Errorf("error decoding user log: %v", err)
		}

		if db_log.Type == core_models.LogTypeSalesPerDayOrder {
			var db_log_sales_order core_models.LogSalesPerDayOrder
			if err := cursor.Decode(&db_log_sales_order); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			db_logs = append(db_logs, db_log_sales_order)
		}

		if db_log.Type == core_models.LogTypeOrderItemRefunded {
			var db_log_order_item_refund core_models.LogOrderItemRefund
			if err := cursor.Decode(&db_log_order_item_refund); err != nil {
				return fmt.Errorf("error decoding user log: %v", err)
			}
			db_logs = append(db_logs, db_log_order_item_refund)
		}

	}

	if len(db_logs) == 0 {
		return nil
	}

	body := struct {
		Data []interface{} `json:"data"`
	}{Data: db_logs}

	json_body, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/api/logs?tenant_id=1", host), bytes.NewBuffer(json_body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Info.Settings.Token))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	http_client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := http_client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error sending request: %v", resp.Status)
	}

	if resp.StatusCode == http.StatusCreated {
		_, err = collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("error deleting logs from db: %v", err)
		}

	}

	s.Logger.Info(fmt.Sprintf("Uploaded %s logs to hub successfully", strconv.Itoa(len(db_logs))))

	return nil
}

func (s *SyncerService) Sync() error {

	settings_svc := SettingsSvc{
		Config: s.Config,
		Logger: s.Logger,
	}
	hubsync, err := settings_svc.Get()
	s.Info = hubsync

	if err != nil {
		return err
	}

	if !hubsync.Settings.Enabled {
		s.Logger.Info("Sync is disabled")
		return nil
	}

	s.Logger.Info("Starting sync...")
	err = s.CopyToBuffer()
	if err != nil {
		return err
	}
	err = s.UploadToServer(hubsync.Settings.ServerHost)
	if err != nil {
		return err
	}

	s.Logger.Info("Sync completed successfully")

	return nil
}
