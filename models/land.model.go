package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LandUnit struct {
	BaseModel   `bson:",inline"`
	Land        string             `bson:"land" json:"land"`
	OwnerID     primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	LesseeID    primitive.ObjectID `bson:"lessee_id,omitempty" json:"lessee_id,omitempty"`
	SizeUnits   int                `bson:"size_units" json:"size_units"`
	IsLeased    bool               `bson:"is_leased" json:"is_leased"`
	IsAvailable bool               `bson:"is_available" json:"is_available"`
	Position    int                `bson:"position" json:"position"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"` // Optional: for future expansion
}

type Land struct {
	BaseModel `bson:",inline"`
	User      primitive.ObjectID `bson:"user" json:"user"`
}

type PlantedCrop struct {
	BaseModel `bson:",inline"`
	// PlantingID        primitive.ObjectID `bson:"planting_id" json:"planting_id"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	CropID            primitive.ObjectID `bson:"crop_id" json:"crop_id"`
	LandUnitIDs       []string           `bson:"land_unit_ids" json:"land_unit_ids"`
	QuantityPlanted   int                `bson:"quantity_planted" json:"quantity_planted"`
	PlantedAt         time.Time          `bson:"planted_at" json:"planted_at"`
	ExpectedHarvestAt time.Time          `bson:"expected_harvest_at" json:"expected_harvest_at"`
	GrowthPercentage  float64            `bson:"growth_percentage" json:"growth_percentage"` // 0.0 to 1.0
	IsHarvested       bool               `bson:"is_harvested" json:"is_harvested"`
	HarvestedAt       time.Time          `bson:"harvested_at,omitempty" json:"harvested_at,omitempty"`
	QualityFactor     float64            `bson:"quality_factor" json:"quality_factor"` // 0.0 to 1.0
	TotalCost         float64            `bson:"total_cost" json:"total_cost"`
	ExpectedYield     int                `bson:"expected_yield" json:"expected_yield"`

	// For partial harvest tracking
	PartialHarvests []PartialHarvest `bson:"partial_harvests,omitempty" json:"partial_harvests,omitempty"`
}

type PartialHarvest struct {
	HarvestID   primitive.ObjectID `bson:"harvest_id" json:"harvest_id"`
	Percentage  float64            `bson:"percentage" json:"percentage"` // 0.0 to 1.0
	Quantity    int                `bson:"quantity" json:"quantity"`
	HarvestedAt time.Time          `bson:"harvested_at" json:"harvested_at"`
	Quality     float64            `bson:"quality" json:"quality"`
}
