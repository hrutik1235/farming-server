package models

import "time"

type Trade struct {
	BaseModel    `bson:",inline"`
	TradeID      string    `bson:"trade_id" json:"trade_id"`
	SellerID     string    `bson:"seller_id" json:"seller_id"`
	BuyerID      string    `bson:"buyer_id" json:"buyer_id"`
	CropID       string    `bson:"crop_id" json:"crop_id"`
	Quantity     int       `bson:"quantity" json:"quantity"`
	PricePerUnit float64   `bson:"price_per_unit" json:"price_per_unit"`
	TotalAmount  float64   `bson:"total_amount" json:"total_amount"`
	Status       string    `bson:"status" json:"status"` // PENDING, ACCEPTED, REJECTED, COMPLETED
	ProposedAt   time.Time `bson:"proposed_at" json:"proposed_at"`
	AcceptedAt   time.Time `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`
	CompletedAt  time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	Reason       string    `bson:"reason,omitempty" json:"reason,omitempty"` // For rejection
}

type MarketPrice struct {
	BaseModel     `bson:",inline"`
	CropID        string    `bson:"crop_id" json:"crop_id"`
	CurrentPrice  float64   `bson:"current_price" json:"current_price"`
	BasePrice     float64   `bson:"base_price" json:"base_price"`
	DemandFactor  float64   `bson:"demand_factor" json:"demand_factor"`
	SupplyFactor  float64   `bson:"supply_factor" json:"supply_factor"`
	LastUpdated   time.Time `bson:"last_updated" json:"last_updated"`
	ValidUntil    time.Time `bson:"valid_until" json:"valid_until"`
	ChangePercent float64   `bson:"change_percent" json:"change_percent"`
}

type PriceHistory struct {
	BaseModel `bson:",inline"`
	CropID    string    `bson:"crop_id" json:"crop_id"`
	Price     float64   `bson:"price" json:"price"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Reason    string    `bson:"reason" json:"reason"` // TRADE, MARKET_UPDATE, etc.
}
