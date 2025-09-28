package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestService struct {
	Client        *mongo.Database
	marketService *MarketService
	cropService   *CropService
}

func NewHarvestService(client *mongo.Database) *HarvestService {
	return &HarvestService{
		Client:        client,
		marketService: NewMarketService(client),
		cropService:   NewCropService(client),
	}
}

func (hs *HarvestService) GetPlantedCrop(userId primitive.ObjectID, plantingId primitive.ObjectID) (*models.PlantedCrop, error) {
	var plantedCroop models.PlantedCrop

	err := hs.Client.Collection(utils.PlantedCropsCollection).FindOne(context.TODO(), bson.M{
		"_id":       plantingId,
		"user_id":   userId,
		"is_active": true,
	}).Decode(&plantedCroop)

	if err != nil {
		return nil, err
	}

	return &plantedCroop, nil
}

func (hs *HarvestService) ValidateHarvest(plantedCrop *models.PlantedCrop) error {
	if plantedCrop.IsHarvested {
		return fmt.Errorf("crop is already harvested")
	}

	plantedCrop.GrowthPercentage = hs.cropService.CalculateCurrentGrowth(*plantedCrop)

	fmt.Println("CURRENT GROWTH PERCENTAGE: ", plantedCrop.GrowthPercentage)

	if plantedCrop.GrowthPercentage < 1 {
		return fmt.Errorf("crop is not fully grown")
	}

	return nil
}

func (hs *HarvestService) ValidatePartialHarvest(plantedCrop *models.PlantedCrop, percentage float64) error {
	if plantedCrop.IsHarvested {
		return fmt.Errorf("crop is already harvested")
	}

	if plantedCrop.GrowthPercentage < 0.5 {
		return fmt.Errorf("crop is not fully grown")
	}

	remainingPercentage := 1.0 - percentage
	if remainingPercentage < 0.1 {
		return fmt.Errorf("harvest percentage too high, consider full harvest")
	}

	return nil
}

func (hs *HarvestService) GetHarvestUnits(landUnitIds []string) ([]models.LandUnit, error) {
	if len(landUnitIds) == 0 {
		return []models.LandUnit{}, nil
	}

	// Convert string IDs to ObjectIDs with validation
	objectIDs, err := utils.ConvertObjectIdsFromStringIds(landUnitIds)

	if err != nil {
		return nil, err
	}

	var landUnits []models.LandUnit

	cursor, err := hs.Client.Collection(utils.LandUnitsCollection).Find(
		context.TODO(),
		bson.M{"_id": bson.M{"$in": objectIDs}},
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &landUnits); err != nil {
		return nil, fmt.Errorf("error decoding results: %v", err)
	}

	return landUnits, nil
}

func (hs *HarvestService) FreeLandUnits(ctx context.Context, landUnitIDs []string) error {
	collection := hs.Client.Collection(utils.LandUnitsCollection)

	landObjectIds, err := utils.ConvertObjectIdsFromStringIds(landUnitIDs)
	if err != nil {
		return err
	}

	_, err = collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": landObjectIds}},
		bson.M{"$set": bson.M{
			"is_available": true,
			"updated_at":   time.Now(),
		}},
	)

	return err
}

func (hs *HarvestService) CalculateQualityFactor(platedCrop *models.PlantedCrop, harvestPercentage float64) float64 {
	baseQuantity := platedCrop.GrowthPercentage

	earlyHarvestPenalty := 0.0

	if platedCrop.GrowthPercentage < 1.0 {
		earlyHarvestPenalty = 1.0 - platedCrop.GrowthPercentage
	}

	partialHarvestEffect := 0.0

	if harvestPercentage < 1.0 {
		partialHarvestEffect = 0.1
	}

	quantity := baseQuantity - earlyHarvestPenalty + partialHarvestEffect

	return math.Max(0.1, math.Min(1.0, quantity))
}

func (hs *HarvestService) CalculateHarvestResult(plantedCrop *models.PlantedCrop, harvestPercentage float64) (*models.HarvestResult, error) {

	currentPrice, err := hs.marketService.GetCurrentPrice(context.TODO(), plantedCrop.CropID)

	if err != nil {
		return nil, err
	}

	qualityFactor := hs.CalculateQualityFactor(plantedCrop, harvestPercentage)

	baseYield := plantedCrop.ExpectedYield

	actualYield := int(float64(baseYield) * harvestPercentage * qualityFactor)

	totalValue := float64(actualYield) * currentPrice * qualityFactor

	return &models.HarvestResult{
		PlantingID:        plantedCrop.ID,
		UserID:            plantedCrop.UserID,
		CropID:            plantedCrop.CropID,
		HarvestType:       "full",
		HarvestPercentage: harvestPercentage,
		Quantity:          actualYield,
		QualityFactor:     qualityFactor,
		BasePrice:         currentPrice,
		ActualPrice:       currentPrice * qualityFactor,
		TotalValue:        totalValue,
		HarvestedAt:       time.Now(),
		IsPartial:         false,
	}, nil
}

func (hs *HarvestService) MarkCropAsHarvested(plantindId primitive.ObjectID, harvestResult *models.HarvestResult) error {

	collection := hs.Client.Collection(utils.PlantedCropsCollection)

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": plantindId},
		bson.M{"$set": bson.M{
			"is_harvested":   true,
			"harvested_at":   time.Now(),
			"actual_yield":   harvestResult.Quantity,
			"quality_factor": harvestResult.QualityFactor,
			"updated_at":     time.Now(),
		}},
	)

	return err
}

func (hs *HarvestService) AddToWarehouse(ctx context.Context, userId primitive.ObjectID, harvestResult *models.HarvestResult) error {
	warehouseItem := models.WarehouseItem{
		UserID:        userId,
		CropID:        harvestResult.CropID,
		Quantity:      harvestResult.Quantity,
		BasePrice:     harvestResult.BasePrice,
		CurrentPrice:  harvestResult.ActualPrice,
		StoredAt:      time.Now(),
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour), // 7 days
		QualityFactor: harvestResult.QualityFactor,
		IsExpired:     false,
		Source:        "harvest",
	}

	err := hs.StoreItemInWareHouse(&warehouseItem)

	return err
}

func (hs *HarvestService) StoreItemInWareHouse(warehouseItem *models.WarehouseItem) error {

	collection := hs.Client.Collection(utils.WarehouseCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defaultCapacity := 1000

	var warehouse models.Warehouse

	err := collection.FindOne(ctx, bson.M{"user_id": warehouseItem.UserID}).Decode(&warehouse)

	if err == mongo.ErrNoDocuments {
		newWarehouse := models.Warehouse{
			BaseModel: models.BaseModel{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:        warehouseItem.UserID,
			TotalCapacity: defaultCapacity,
			UsedCapacity:  warehouseItem.Quantity,
			Items:         []models.WarehouseItem{*warehouseItem},
		}

		_, err := collection.InsertOne(ctx, newWarehouse)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if warehouse.UsedCapacity+warehouseItem.Quantity > warehouse.TotalCapacity {
			return fmt.Errorf("warehouse is full")
		}

		update := bson.M{
			"$push": bson.M{"items": warehouseItem},
			"$inc":  bson.M{"used_capacity": warehouseItem.Quantity},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		_, err = collection.UpdateOne(ctx, bson.M{"user_id": warehouseItem.UserID}, update)

		if err != nil {
			return fmt.Errorf("failed to update warehouse: %v", err)
		}
	}
	return nil
}
