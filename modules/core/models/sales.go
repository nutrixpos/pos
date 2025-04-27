package models

// SalesPerDayOrder represents an order and its associated costs for a specific day.
type SalesPerDayOrder struct {
	Order Order      `json:"order" bson:"order,inline"`
	Costs []ItemCost `json:"costs" bson:"costs"`
}

type OrderItemRefundMaterial struct {
	MaterialId         string  `json:"material_id" bson:"material_id"`
	EntryId            string  `json:"entry_id" bson:"entry_id"`
	InventoryReturnQty float64 `json:"inventory_return_qty" bson:"inventory_return_qty"`
	DisposeQty         float64 `json:"dispose_qty" bson:"dispose_qty"`
	WasteQty           float64 `json:"waste_qty" bson:"waste_qty"`
	CostPerUnit        float64 `json:"cost_per_unit" bson:"cost_per_unit"`
	Comment            string  `json:"comment" bson:"comment"`
}

type OrderItemRefundProductAdd struct {
	ProductId string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Comment   string  `json:"comment" bson:"comment"`
}

type SalesPerDayRefund struct {
	OrderId         string                      `json:"order_id" bson:"order_id"`
	ItemId          string                      `json:"order_item_id" bson:"order_item_id"`
	ProductId       string                      `json:"product_id" bson:"product_id"`
	Reason          string                      `json:"reason" bons:"reason"`
	Amount          float64                     `json:"amount" bson:"amount"`
	ItemCost        float64                     `json:"item_cost" bson:"item_cost"`
	Destination     string                      `json:"destination" bson:"destination"`
	MaterialRerunds []OrderItemRefundMaterial   `json:"material_refunds" bson:"material_refunds"`
	ProductAdd      []OrderItemRefundProductAdd `json:"products_add" bson:"products_add"`
}

// SalesPerDay aggregates sales data for a specific day, including total costs and sales.
type SalesPerDay struct {
	Id           string              `json:"id" bson:"id,omitempty"`
	Date         string              `json:"date" bson:"date"`
	Orders       []SalesPerDayOrder  `json:"orders" bson:"orders"`
	Refunds      []SalesPerDayRefund `json:"refunds" bson:"refunds"`
	Costs        float64             `json:"costs" bson:"costs"`
	TotalSales   float64             `json:"total_sales" bson:"total_sales"`
	RefundsValue float64             `json:"refunds_value" bson:"refunds_value"`
}
