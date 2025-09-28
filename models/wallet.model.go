package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	BaseModel     `bson:",inline"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Balance       float64            `bson:"balance" json:"balance"`
	TotalEarnings float64            `bson:"total_earnings" json:"total_earnings"`
	TotalSpent    float64            `bson:"total_spent" json:"total_spent"`
	LastUpdated   time.Time          `bson:"last_updated" json:"last_updated"`
}

type Transaction struct {
	BaseModel     `bson:",inline"`
	TransactionID string    `bson:"transaction_id" json:"transaction_id"`
	UserID        string    `bson:"user_id" json:"user_id"`
	Type          string    `bson:"type" json:"type"` // INCOME, EXPENSE
	Amount        float64   `bson:"amount" json:"amount"`
	Description   string    `bson:"description" json:"description"`
	Category      string    `bson:"category" json:"category"`                             // CROP_SALE, LEASE_INCOME, PLANTING_COST, etc.
	ReferenceID   string    `bson:"reference_id,omitempty" json:"reference_id,omitempty"` // TradeID, LeaseID, etc.
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
	BalanceAfter  float64   `bson:"balance_after" json:"balance_after"`
}
