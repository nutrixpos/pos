package models

type Customer struct {
	Id      string `json:"id" bson:"id" mapstructure:"id"`
	Name    string `json:"name" bson:"name" mapstructure:"name"`
	Phone   string `json:"phone" bson:"phone" mapstructure:"phone"`
	Address string `json:"address" bson:"address" mapstructure:"address"`
}
