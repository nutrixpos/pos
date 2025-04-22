package models

const (
	TypeDisposalMaterial = "disposal_material"
	TypeDisposalProduct  = "disposal_product"
)

// Disposal used to store a returned order item or material which is too unique and modified to return to normal inventory state
type Disposal struct {
	Id       string  `json:"id" bson:"id"`
	OrderId  string  `json:"order_id" bson:"order_id"`
	Type     string  `json:"type" bson:"type"` // material or order_item
	Quantity float64 `json:"quantity" bson:"quantity"`
	Comment  string  `json:"comment" bson:"comment"`
}

type MaterialDisposal struct {
	Disposal   `json:",inline"`
	MaterialId string `json:"material_id" bson:"material_id"`
	EntryId    string `json:"entry_id" bson:"entry_id"`
}

type ProductDisposal struct {
	Disposal `json:",inline"`
	Item     OrderItem `json:"item" bson:"item"`
}
