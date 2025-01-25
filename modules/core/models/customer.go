package models

type Customer struct {
	Id      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Phone   string `json:"phone" bson:"phone"`
	Address string `json:"address" bson:"address"`
}
