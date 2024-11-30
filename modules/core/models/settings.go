package models

// OrderQueueSettings represents the configuration settings for an order queue
type OrderQueueSettings struct {
	Prefix string `json:"prefix" bson:"prefix"`
	Next   uint32 `json:"next" bson:"next"`
}

// OrderSettings represents the configuration settings for orders
type OrderSettings struct {
	Queues []OrderQueueSettings `json:"queues" bson:"queues"`
}

// Settings represents the configuration settings structure
type Settings struct {
	Id        string `bson:"id,omitempty" json:"id"`
	Inventory struct {
		DefaultInventoryQuantityWarn float64 `json:"default_inventory_quantity_warn" bson:"default_inventory_quantity_warn"`
	} `bson:"inventory" json:"inventory"`
	Orders OrderSettings `bson:"orders" json:"orders"`
}
