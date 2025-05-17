package models

// SalesPerDayOrder represents an order and its associated costs for a specific day.
type SalesPerDayOrder struct {
	Order Order      `json:"order" bson:"order,inline" mapstructure:"order"`
	Costs []ItemCost `json:"costs" bson:"costs" mapstructure:"costs"`
}

type OrderItemRefundMaterial struct {
	MaterialId         string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	EntryId            string  `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
	InventoryReturnQty float64 `json:"inventory_return_qty" bson:"inventory_return_qty" mapstructure:"inventory_return_qty"`
	DisposeQty         float64 `json:"dispose_qty" bson:"dispose_qty" mapstructure:"dispose_qty"`
	WasteQty           float64 `json:"waste_qty" bson:"waste_qty" mapstructure:"waste_qty"`
	CostPerUnit        float64 `json:"cost_per_unit" bson:"cost_per_unit" mapstructure:"cost_per_unit"`
	Comment            string  `json:"comment" bson:"comment" mapstructure:"comment"`
}

type OrderItemRefundProductAdd struct {
	ProductId string  `json:"product_id" mapstructure:"product_id"`
	Quantity  float64 `json:"quantity" mapstructure:"quantity"`
	Comment   string  `json:"comment" bson:"comment" mapstructure:"comment"`
}

type SalesPerDayRefund struct {
	OrderId         string                      `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	ItemId          string                      `json:"order_item_id" bson:"order_item_id" mapstructure:"order_item_id"`
	ProductId       string                      `json:"product_id" bson:"product_id" mapstructure:"product_id"`
	Reason          string                      `json:"reason" bons:"reason" mapstructure:"reason"`
	Amount          float64                     `json:"amount" bson:"amount" mapstructure:"amount"`
	ItemCost        float64                     `json:"item_cost" bson:"item_cost" mapstructure:"item_cost"`
	Destination     string                      `json:"destination" bson:"destination" mapstructure:"destination"`
	MaterialRerunds []OrderItemRefundMaterial   `json:"material_refunds" bson:"material_refunds" mapstructure:"material_refunds"`
	ProductAdd      []OrderItemRefundProductAdd `json:"products_add" bson:"products_add" mapstructure:"products_add"`
}

// SalesPerDay aggregates sales data for a specific day, including total costs and sales.
type SalesPerDay struct {
	Id           string              `json:"id" bson:"id,omitempty" mapstructure:"id"`
	Date         string              `json:"date" bson:"date" mapstructure:"date"`
	Orders       []SalesPerDayOrder  `json:"orders" bson:"orders" mapstructure:"orders"`
	Refunds      []SalesPerDayRefund `json:"refunds" bson:"refunds" mapstructure:"refunds"`
	Costs        float64             `json:"costs" bson:"costs" mapstructure:"costs"`
	TotalSales   float64             `json:"total_sales" bson:"total_sales" mapstructure:"total_sales"`
	RefundsValue float64             `json:"refunds_value" bson:"refunds_value" mapstructure:"refunds_value"`
}
