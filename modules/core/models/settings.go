package models

// OrderQueueSettings represents the configuration settings for an order queue
type OrderQueueSettings struct {
	Prefix string `json:"prefix" bson:"prefix" mapstructure:"prefix"`
	Next   uint32 `json:"next" bson:"next" mapstructure:"next"`
}

// OrderSettings represents the configuration settings for orders
type OrderSettings struct {
	Queues                       []OrderQueueSettings `json:"queues" bson:"queues" mapstructure:"queues"`
	DefaultCostCalculationMethod string               `json:"default_cost_calculation_method" bson:"default_cost_calculation_method" mapstructure:"default_cost_calculation_method"`
}

type LanguageSettings struct {
	Code     string `json:"code" bson:"code" mapstructure:"code"`
	Language string `json:"language" bson:"language" mapstructure:"language"`
}

// MaterialSettings represents settings associated with a material, such as stock alert threshold.
type MaterialSettings struct {
	StockAlertTreshold float64 `json:"stock_alert_treshold" bson:"stock_alert_treshold" mapstructure:"stock_alert_treshold"`
}

type PrinterSettings struct {
	Host string `bson:"host" json:"host" mapstructure:"host"`
}

// Settings represents the configuration settings structure
type Settings struct {
	Id             string           `bson:"id,omitempty" json:"id" mapstructure:"id"`
	Inventory      MaterialSettings `bson:"inventory" json:"inventory" mapstructure:"inventory"`
	Orders         OrderSettings    `bson:"orders" json:"orders" mapstructure:"orders"`
	Language       LanguageSettings `bson:"language" json:"language" mapstructure:"language"`
	ReceiptPrinter PrinterSettings  `bson:"receipt_printer" json:"receipt_printer" mapstructure:"receipt_printer"`
}
