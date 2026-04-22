package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username" mapstructure:"username"`
	Email        string             `json:"email" bson:"email" mapstructure:"email"`
	PasswordHash string             `json:"-" bson:"password_hash" mapstructure:"password_hash"`
	Roles        []string           `json:"roles" bson:"roles" mapstructure:"roles"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Roles:    u.Roles,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterRequest struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}