package models

// SalesPerDayOrder represents an order and its associated costs for a specific day.
type SalesPerDayOrder struct {
	Order Order      `json:"order" bson:"order,inline"`
	Costs []ItemCost `json:"costs" bson:"costs"`
}

type SalesPerDayRefund struct {
	OrderId string `json:"order_id" bson:"order_id"`
	ItemId  string `json:"item_id" bson:"item_id"`
	Reason  string `json:"reason" bson:"reason"`
}

// SalesPerDay aggregates sales data for a specific day, including total costs and sales.
type SalesPerDay struct {
	Id         string              `json:"id" bson:"id,omitempty"`
	Date       string              `json:"date" bson:"date"`
	Orders     []SalesPerDayOrder  `json:"orders" bson:"orders"`
	Refunds    []SalesPerDayRefund `json:"refunds" bson:"refunds"`
	Costs      float64             `json:"costs" bson:"costs"`
	TotalSales float64             `json:"total_sales" bson:"total_sales"`
}
