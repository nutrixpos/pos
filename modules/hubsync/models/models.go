package models

type Settings struct {
	Enabled       bool   `bson:"enabled" json:"enabled" mapstructure:"enabled"`
	ServerHost    string `bson:"server_host" json:"server_host" mapstructure:"server_host"`
	Token         string `bson:"token" json:"token" mapstructure:"token"`
	SyncInterval  int64  `json:"sync_interval" bson:"sync_interval"`
	BufferSize    int64  `json:"buffer_size" bson:"buffer_size"`
	SyncSales     bool   `json:"sync_sales" bson:"sync_sales"`
	SyncInventory bool   `json:"sync_inventory" bson:"sync_inventory"`
}

type Hubsync struct {
	Id           string   `json:"id" bson:"_id"`
	Settings     Settings `json:"settings" bson:"settings"`
	LastSynced   int64    `json:"last_synced" bson:"last_synced"`
	SyncProgress float64  `json:"sync_progress" bson:"sync_progress"`
}
