package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func ConvertObjectIdsFromStringIds(ids []string) ([]primitive.ObjectID, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))

	for i, idStr := range ids {
		if !primitive.IsValidObjectID(idStr) {
			return nil, fmt.Errorf("invalid object ID at index %d: %s", i, idStr)
		}

		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("error converting ID %s: %v", idStr, err)
		}
		objectIDs[i] = objID
	}

	return objectIDs, nil
}
