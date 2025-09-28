package service

import (
	"context"
	"fmt"

	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestService struct {
	Client *mongo.Database
}

func NewHarvestService(client *mongo.Database) *HarvestService {
	return &HarvestService{
		Client: client,
	}
}

func (hs *HarvestService) GetPlantedCrop(userId primitive.ObjectID, plantingId primitive.ObjectID) (*models.PlantedCrop, error) {
	var plantedCroop models.PlantedCrop

	err := hs.Client.Collection(utils.PlantedCropsCollection).FindOne(context.TODO(), bson.M{
		"planting_id": plantingId,
		"user_id":     userId,
		"is_active":   true,
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

	if plantedCrop.GrowthPercentage < 100 {
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
	var landUnits []models.LandUnit

	cursor, err := hs.Client.Collection(utils.LandUnitsCollection).Find(context.TODO(), bson.M{"_id": bson.M{"$in": landUnitIds}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &landUnits)

	return landUnits, err
}

// func (hs *HarvestService) freeLandUnits(ctx context.Context, landUnitIDs []string) error {
//     collection := hs.Client.Collection(utils.LandUnitsCollection)

//     _, err := collection.UpdateMany(
//         ctx,
//         bson.M{"_id": bson.M{"$in": landUnitIDs}},
//         bson.M{"$set": bson.M{
//             "is_available": true,
//             "updated_at":   time.Now(),
//         }},
//     )

//     return err
// }

// func (hs *HarvestService) CalculateHarvestResult(plantedCrop *models.PlantedCrop, harvestPercentage float64) (*models.HarvestResult, error) {
// 	currentPrice, err := hs.
// }
