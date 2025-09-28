package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HarvestResult struct {
	HarvestID         string             `bson:"harvest_id" json:"harvest_id"`
	PlantingID        primitive.ObjectID `bson:"planting_id" json:"planting_id"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	CropID            primitive.ObjectID `bson:"crop_id" json:"crop_id"`
	HarvestType       string             `bson:"harvest_type" json:"harvest_type"` // "full" or "partial"
	HarvestPercentage float64            `bson:"harvest_percentage" json:"harvest_percentage"`
	Quantity          int                `bson:"quantity" json:"quantity"`
	QualityFactor     float64            `bson:"quality_factor" json:"quality_factor"`
	BasePrice         float64            `bson:"base_price" json:"base_price"`
	ActualPrice       float64            `bson:"actual_price" json:"actual_price"`
	TotalValue        float64            `bson:"total_value" json:"total_value"`
	HarvestedAt       time.Time          `bson:"harvested_at" json:"harvested_at"`
	IsPartial         bool               `bson:"is_partial" json:"is_partial"`
}
