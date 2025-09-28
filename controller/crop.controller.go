package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/service"
	"github.com/hrutik1235/farming-server/types"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CropController struct {
	service     *service.CropService
	userService *service.UserService
	dbClient    *mongo.Database
}

func NewCropController(dbClient *mongo.Database) *CropController {
	return &CropController{
		service:     service.NewCropService(dbClient),
		userService: service.NewUserService(dbClient),
		dbClient:    dbClient,
	}
}

func (cropController *CropController) CreateCrop(c *gin.Context) {
	body := c.MustGet("body").(types.CreateCrop)

	crop := models.Crop{
		Name:            body.Name,
		Description:     body.Description,
		BasePrice:       body.BasePrice,
		GrowthTimeHours: body.GrowthTimeHours,
		YieldPerUnit:    body.YieldPerUnit,
		CostPerUnit:     body.CostPerUnit,
		IsActive:        true,
	}

	existingCrop, _ := cropController.service.GetCropByName(crop.Name)

	if existingCrop.Name == crop.Name {
		c.JSON(http.StatusBadRequest, utils.NewHttpError(c, "Crop already exists", http.StatusBadRequest))
		return
	}

	_, err := cropController.dbClient.Collection(utils.CropsCollection).InsertOne(context.TODO(), crop)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, utils.NewHttpError(c, "Crop created successfully", http.StatusCreated))
}

func (cropController *CropController) GetAllCrops(c *gin.Context) {

	crops, err := cropController.service.GetAllCrops(context.TODO())

	fmt.Println(len(crops))

	fmt.Println(crops)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": crops,
	})
}

func (cropController *CropController) GetAllPlantedCrops(c *gin.Context) {
	userId := c.GetHeader("user_id")

	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	plantedCrops, err := cropController.service.GetUserPlantedCrops(context.TODO(), userObjectId, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plantedCrops,
	})
}

func (cropController *CropController) PlantCrop(c *gin.Context) {
	userId := c.GetHeader("user_id")
	landUnits := c.MustGet("body").(types.PlantCrop)
	cropId := c.Param("cropid")

	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	cropObjectId, _ := primitive.ObjectIDFromHex(cropId)

	availableLand, err := cropController.service.GetAvailableLandUnits(context.TODO(), userObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	if len(availableLand) < landUnits.LandUnits {
		c.JSON(http.StatusBadRequest, utils.NewHttpError(c, "Not enough land units", http.StatusBadRequest))
		return
	}

	crop, err := cropController.service.GetCropById(cropObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	wallet, err := cropController.userService.GetUserWallet(userObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, "Wallet Error", http.StatusInternalServerError))
		return
	}

	totalCost := crop.CostPerUnit * float64(landUnits.LandUnits)

	if wallet.Balance < totalCost {
		c.JSON(http.StatusBadRequest, utils.NewHttpError(c, "Not enough balance", http.StatusBadRequest))
		return
	}

	if err := cropController.userService.DeductPlantingCost(userObjectId, totalCost); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	landUnitIds := cropController.service.GetLandUnitIDs(availableLand[:landUnits.LandUnits])

	fmt.Println("Land Unit IDs: ", landUnitIds)

	if err := cropController.service.MarkLandUnitsOccupied(context.TODO(), landUnitIds); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	plantedCrop, planterr := cropController.service.CreatePlantedCrop(context.TODO(), userObjectId, crop, landUnitIds, landUnits.LandUnits, totalCost)

	if planterr != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plantedCrop,
	})
}
