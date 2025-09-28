package models

import "time"

type Lease struct {
	BaseModel      `bson:",inline"`
	LeaseID        string    `bson:"lease_id" json:"lease_id"`
	LandownerID    string    `bson:"landowner_id" json:"landowner_id"`
	TenantID       string    `bson:"tenant_id" json:"tenant_id"`
	LandUnitIDs    []string  `bson:"land_unit_ids" json:"land_unit_ids"`
	TotalLandUnits int       `bson:"total_land_units" json:"total_land_units"`
	Status         string    `bson:"status" json:"status"` // PENDING, ACTIVE, REJECTED, EXPIRED, CANCELLED
	DurationDays   int       `bson:"duration_days" json:"duration_days"`
	LeasePrice     float64   `bson:"lease_price" json:"lease_price"`
	RequestedAt    time.Time `bson:"requested_at" json:"requested_at"`
	AcceptedAt     time.Time `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`
	StartTime      time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime        time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Reason         string    `bson:"reason,omitempty" json:"reason,omitempty"` // For rejection/cancellation
}

type LeaseActivity struct {
	BaseModel    `bson:",inline"`
	LeaseID      string    `bson:"lease_id" json:"lease_id"`
	ActivityType string    `bson:"activity_type" json:"activity_type"` // PLANT, HARVEST, WATER, etc.
	TenantID     string    `bson:"tenant_id" json:"tenant_id"`
	CropID       string    `bson:"crop_id,omitempty" json:"crop_id,omitempty"`
	Description  string    `bson:"description" json:"description"`
	Timestamp    time.Time `bson:"timestamp" json:"timestamp"`
}
