package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
}

type User struct {
	BaseModel     `bson:",inline"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"` // Custom ID for easier reference
	Username      string             `bson:"username" json:"username"`
	Email         string             `bson:"email" json:"email"`
	DisplayName   string             `bson:"display_name" json:"display_name"`
	ServerAddress string             `bson:"server_address" json:"server_address"` // IP:Port for manual connections
	LastLogin     time.Time          `bson:"last_login" json:"last_login"`
	IsOnline      bool               `bson:"is_online" json:"is_online"`
}

func (user *User) Save(userCol *mongo.Collection) (*mongo.InsertOneResult, error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	details, err := userCol.InsertOne(ctxWithTimeout, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %v", err.Error())
	}
	return details, nil
}

type Crop struct {
	BaseModel       `bson:",inline"`
	Name            string  `bson:"name" json:"name"`
	BasePrice       float64 `bson:"base_price" json:"base_price"`
	GrowthTimeHours int     `bson:"growth_time_hours" json:"growth_time_hours"`
	YieldPerUnit    int     `bson:"yield_per_unit" json:"yield_per_unit"`
	CostPerUnit     float64 `bson:"cost_per_unit" json:"cost_per_unit"`
	Description     string  `bson:"description,omitempty" json:"description,omitempty"`
	IsActive        bool    `bson:"is_active" json:"is_active"`
}
