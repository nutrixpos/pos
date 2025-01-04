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

type LanguageSettings struct {
	Code     string `json:"code" bson:"code"`
	Language string `json:"language" bson:"language"`
}

// Settings represents the configuration settings structure
type Settings struct {
	Id        string `bson:"id,omitempty" json:"id"`
	Inventory struct {
		DefaultInventoryQuantityWarn float64 `json:"default_inventory_quantity_warn" bson:"default_inventory_quantity_warn"`
	} `bson:"inventory" json:"inventory"`
	Orders         OrderSettings    `bson:"orders" json:"orders"`
	Language       LanguageSettings `bson:"language" json:"language"`
	ReceiptPrinter struct {
		Host string `bson:"host" json:"host"`
	} `bson:"receipt_printer" json:"receipt_printer"`
}
