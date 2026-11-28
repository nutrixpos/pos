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

const (
	PrinterConnectionTypeNetwork = "network"
	PrinterConnectionTypeUSB     = "usb"
)

type PrinterSettings struct {
	// ConnectionType selects how the receipt printer is reached: "network"
	// (default when empty, for backward compatibility) or "usb".
	ConnectionType string `bson:"connection_type" json:"connection_type" mapstructure:"connection_type"`
	// Host is the printer host used for network printers.
	Host string `bson:"host" json:"host" mapstructure:"host"`
	// Name is the OS printer queue name used for USB printers.
	Name string `bson:"name" json:"name" mapstructure:"name"`
}

// IsUSB reports whether the printer should be reached through the OS USB
// printer queue instead of the network.
func (ps PrinterSettings) IsUSB() bool {
	return ps.ConnectionType == PrinterConnectionTypeUSB
}

// PrintingSettings represents settings used to render receipts.
type PrintingSettings struct {
	// BrowserPath is the Chromium-based browser executable used to render
	// receipts. When empty the browser is auto-detected.
	BrowserPath string `bson:"browser_path" json:"browser_path" mapstructure:"browser_path"`
}

// Settings represents the configuration settings structure
type Settings struct {
	Id                    string           `bson:"id,omitempty" json:"id" mapstructure:"id"`
	Inventory             MaterialSettings `bson:"inventory" json:"inventory" mapstructure:"inventory"`
	Orders                OrderSettings    `bson:"orders" json:"orders" mapstructure:"orders"`
	Language              LanguageSettings `bson:"language" json:"language" mapstructure:"language"`
	ClientReceiptPrinter  PrinterSettings  `bson:"client_receipt_printer" json:"client_receipt_printer" mapstructure:"client_receipt_printer"`
	KitchenReceiptPrinter PrinterSettings  `bson:"kitchen_receipt_printer" json:"kitchen_receipt_printer" mapstructure:"kitchen_receipt_printer"`
	Printing              PrintingSettings `bson:"printing" json:"printing" mapstructure:"printing"`
	PaymentSources        []PaymentSource  `bson:"payment_sources" json:"payment_sources" mapstructure:"payment_sources"`
	// ShopMode determines the operational mode: "" (unset/first-run), "kitchen", or "retail"
	ShopMode           string `bson:"shop_mode" json:"shop_mode" mapstructure:"shop_mode"`
	AutoOpenCashDrawer bool   `bson:"auto_open_cash_drawer" json:"auto_open_cash_drawer" mapstructure:"auto_open_cash_drawer"`
}

type PaymentSource struct {
	Name string `bson:"name" json:"name" mapstructure:"name"`
}
