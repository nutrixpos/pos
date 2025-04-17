package dto

const (
	DTOOrderItemRefundDestination_Inventory = "inventory"
	DTOOrderItemRefundDestination_Disposals = "disposals"
	DTOOrderItemRefundDestination_Waste     = "waste"
	DTOOrderItemRefundDestination_Custom    = "custom"
)

type OrderItemRefundRequest struct {
	OrderId         string                         `json:"order_id" bson:"order_id"`
	ProductId       string                         `json:"product_id" bson:"product_id"`
	RefundValue     float64                        `json:"refund_value"`
	Destination     string                         `json:"destination" bson:"destination"`
	MaterialRerunds []OrderItemRefundMaterialDTO   `json:"material_refunds"`
	ProductAdd      []OrderItemRefundProductAddDTO `json:"products_add"`
}

type OrderItemRefundMaterialDTO struct {
	MaterialId         string  `json:"material_id"`
	EntryId            string  `json:"entry_id"`
	InventoryReturnQty float64 `json:"inventory_return_qty"`
	DisposeQty         float64 `json:"dispose_return_qty"`
	WasteQty           float64 `json:"waste_return_qty"`
	Comment            string  `json:"comment" bson:"comment"`
}

type OrderItemRefundProductAddDTO struct {
	ProductId string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Comment   string  `json:"comment" bson:"comment"`
}
