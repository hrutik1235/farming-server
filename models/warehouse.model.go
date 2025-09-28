package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WarehouseItem struct {
	BaseModel     `bson:",inline"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CropID        primitive.ObjectID `bson:"crop_id" json:"crop_id"`
	Quantity      int                `bson:"quantity" json:"quantity"`
	BasePrice     float64            `bson:"base_price" json:"base_price"`
	CurrentPrice  float64            `bson:"current_price" json:"current_price"`
	StoredAt      time.Time          `bson:"stored_at" json:"stored_at"`
	ExpiresAt     time.Time          `bson:"expires_at" json:"expires_at"`
	QualityFactor float64            `bson:"quality_factor" json:"quality_factor"` // 0.0 to 1.0
	IsExpired     bool               `bson:"is_expired" json:"is_expired"`
	Source        string             `bson:"source" json:"source"` // HARVEST, TRADE, etc.
}

type Warehouse struct {
	BaseModel     `bson:",inline"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	TotalCapacity int                `bson:"total_capacity" json:"total_capacity"`
	UsedCapacity  int                `bson:"used_capacity" json:"used_capacity"`

	// Embedded items (or could be separate collection)
	Items []WarehouseItem `bson:"items,omitempty" json:"items,omitempty"`
}
