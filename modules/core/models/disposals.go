package models

const (
	TypeDisposalMaterial = "disposal_material"
	TypeDisposalProduct  = "disposal_product"
)

// Disposal used to store a returned order item or material which is too unique and modified to return to normal inventory state
type Disposal struct {
	Id       string  `json:"id" bson:"id" mapstructure:"id"`
	OrderId  string  `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	Type     string  `json:"type" bson:"type" mapstructure:"type"` // material or order_item
	Quantity float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Comment  string  `json:"comment" bson:"comment" mapstructure:"comment"`
}

type MaterialDisposal struct {
	Disposal   `json:",inline" mapstructure:",squash"`
	MaterialId string `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	EntryId    string `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
}

type ProductDisposal struct {
	Disposal `json:",inline" mapstructure:",squash"`
	Item     OrderItem `json:"item" bson:"item" mapstructure:"item"`
}
