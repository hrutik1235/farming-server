package service

import (
	"context"
	"errors"
	"time"

	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MarketService struct {
	client *mongo.Database
}

func NewMarketService(client *mongo.Database) *MarketService {
	return &MarketService{client: client}
}

// GetCurrentPrice with lazy initialization - creates market price if not exists
func (ms *MarketService) GetCurrentPrice(ctx context.Context, cropID primitive.ObjectID) (float64, error) {
	collection := ms.client.Collection(utils.MarketPricesCollection)

	var marketPrice models.MarketPrice
	err := collection.FindOne(ctx, bson.M{
		"_id":         cropID,
		"is_active":   true,
		"valid_until": bson.M{"$gt": time.Now()},
	}).Decode(&marketPrice)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Market price doesn't exist - create it lazily
			return ms.CreateAndGetMarketPrice(ctx, cropID)
		}
		return 0, err
	}

	return marketPrice.CurrentPrice, nil
}

// Create market price for crop and return the price
func (ms *MarketService) CreateAndGetMarketPrice(ctx context.Context, cropID primitive.ObjectID) (float64, error) {
	// Get crop base price first
	basePrice, err := ms.GetBasePriceFromCrop(ctx, cropID)
	if err != nil {
		return 0, err
	}

	// Create new market price
	marketPrice := models.MarketPrice{
		BaseModel: models.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive:  true,
		},
		CropID:       cropID,
		CurrentPrice: basePrice,
		BasePrice:    basePrice,
		DemandFactor: 1.0, // Neutral demand
		SupplyFactor: 1.0, // Neutral supply
		LastUpdated:  time.Now(),
		ValidUntil:   time.Now().Add(24 * time.Hour), // Valid for 24 hours
	}

	// Save to database
	collection := ms.client.Collection(utils.MarketPricesCollection)
	_, err = collection.InsertOne(ctx, marketPrice)
	if err != nil {
		return 0, err
	}

	return marketPrice.CurrentPrice, nil
}

// Get base price from crop definition
func (ms *MarketService) GetBasePriceFromCrop(ctx context.Context, cropID primitive.ObjectID) (float64, error) {
	collection := ms.client.Collection(utils.CropsCollection)

	var crop models.Crop
	err := collection.FindOne(ctx, bson.M{
		"_id":       cropID,
		"is_active": true,
	}).Decode(&crop)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, errors.New("crop not found")
		}
		return 0, err
	}

	return crop.BasePrice, nil
}
