package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CropService struct {
	Client *mongo.Database

	userService *UserService
}

func NewCropService(client *mongo.Database) *CropService {
	return &CropService{
		Client:      client,
		userService: NewUserService(client),
	}
}

func (c *CropService) GetCropById(cropId primitive.ObjectID) (*models.Crop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var crop models.Crop

	err := c.Client.Collection(utils.CropsCollection).FindOne(ctx, bson.M{"id": cropId}).Decode(&crop)

	if err != nil {
		return nil, err
	}

	return &crop, nil
}

func (c *CropService) GetCropByName(cropName string) (*models.Crop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var crop models.Crop

	fmt.Println("CROP NAME", cropName)

	err := c.Client.Collection(utils.CropsCollection).FindOne(ctx, bson.M{"name": cropName}).Decode(&crop)

	fmt.Println("FOUNDED CROP: ", crop)

	return &crop, err
}

func (c *CropService) GetPlantedCrops(userId primitive.ObjectID, plantingId primitive.ObjectID) ([]models.PlantedCrop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.Client.Collection(utils.PlantedCropsCollection).Find(ctx, bson.M{"user_id": userId, "crop_id": plantingId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var plantedCrops []models.PlantedCrop

	err = cur.All(ctx, &plantedCrops)

	if err != nil {
		return nil, err
	}

	return plantedCrops, nil
}

func (cs *CropService) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	collection := cs.Client.Collection(utils.CropsCollection)

	var crops []models.Crop
	cursor, err := collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &crops)
	return crops, err
}

func (cs *CropService) GetAvailableLandUnits(ctx context.Context, userId primitive.ObjectID) ([]models.LandUnit, error) {
	collection := cs.Client.Collection(utils.LandUnitsCollection)

	fmt.Println("USER ID: ", userId)

	var landUnits []models.LandUnit
	cursor, err := collection.Find(ctx, bson.M{"owner_id": userId, "is_available": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &landUnits)

	fmt.Println("LAND UNITS: ", len(landUnits))
	return landUnits, err
}

func (p *CropService) GetLandUnitIDs(landUnits []models.LandUnit) []string {
	ids := make([]string, len(landUnits))
	for i, unit := range landUnits {
		ids[i] = unit.ID.Hex()
	}
	return ids
}

func (p *CropService) MarkLandUnitsOccupied(ctx context.Context, landUnitIDs []string) error {
	collection := p.Client.Collection(utils.LandUnitsCollection)

	_, err := collection.UpdateMany(
		ctx,
		bson.M{"land_unit_id": bson.M{"$in": landUnitIDs}},
		bson.M{
			"$set": bson.M{
				"is_available": false,
				"updated_at":   time.Now(),
			},
		},
	)

	return err
}

func (p *CropService) CreatePlantedCrop(ctx context.Context, userID primitive.ObjectID, crop *models.Crop, landUnitIDs []string, landUnits int, totalCost float64) (*models.PlantedCrop, error) {
	collection := p.Client.Collection(utils.PlantedCropsCollection)

	plantedAt := time.Now()
	expectedHarvest := plantedAt.Add(time.Duration(crop.GrowthTimeHours) * time.Hour)

	plantedCrop := &models.PlantedCrop{
		BaseModel: models.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: plantedAt,
			UpdatedAt: plantedAt,
			IsActive:  true,
		},
		UserID:            userID,
		CropID:            crop.ID,
		LandUnitIDs:       landUnitIDs, // Store all land unit IDs used
		QuantityPlanted:   landUnits,
		PlantedAt:         plantedAt,
		ExpectedHarvestAt: expectedHarvest,
		GrowthPercentage:  0.0,
		IsHarvested:       false,
		TotalCost:         totalCost,
		ExpectedYield:     crop.YieldPerUnit * landUnits,
	}

	_, err := collection.InsertOne(ctx, plantedCrop)
	if err != nil {
		return nil, err
	}

	return plantedCrop, nil
}

func (p *CropService) GetUserPlantedCrops(ctx context.Context, userID primitive.ObjectID, activeOnly bool) ([]models.PlantedCrop, error) {
	collection := p.Client.Collection(utils.PlantedCropsCollection)

	filter := bson.M{"user_id": userID, "is_active": activeOnly}
	if activeOnly {
		filter["is_harvested"] = false
	}

	var plantedCrops []models.PlantedCrop
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &plantedCrops)
	return plantedCrops, err
}
