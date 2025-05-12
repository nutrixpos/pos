package models

type Hubsync struct {
	Id           string  `json:"id" bson:"_id"`
	LastSynced   int64   `json:"last_synced" bson:"last_synced"`
	SyncInterval int64   `json:"sync_interval" bson:"sync_interval"`
	SyncProgress float64 `json:"sync_progress" bson:"sync_progress"`
	BufferSize   int64   `json:"buffer_size" bson:"buffer_size"`
}
