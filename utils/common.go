package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HttpError[T any] struct {
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
	Path       string `json:"path"`
	Error      T      `json:"error"`
}

// Generic constructor
func NewHttpError[T any](c *gin.Context, message T, status int) *HttpError[T] {
	return &HttpError[T]{
		StatusCode: status,
		Timestamp:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		Path:       c.FullPath(),
		Error:      message,
	}
}

const (
	UsersCollection           = "users"
	CropsCollection           = "crops"
	LandsCollection           = "lands"
	LandUnitsCollection       = "land_units"
	PlantedCropsCollection    = "planted_crops"
	LeasesCollection          = "leases"
	WarehouseCollection       = "warehouse"
	TradesCollection          = "trades"
	MarketPricesCollection    = "market_prices"
	WalletsCollection         = "wallets"
	TransactionsCollection    = "transactions"
	PeerConnectionsCollection = "peer_connections"
	EventsCollection          = "events"
)

const (
	InitialLandUnitSize = 100
)
