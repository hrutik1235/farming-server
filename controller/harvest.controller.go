package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/service"
	"github.com/hrutik1235/farming-server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestController struct {
	service *service.HarvestService
}

func NewHarvestController(dbClient *mongo.Database) *HarvestController {
	return &HarvestController{
		service: service.NewHarvestService(dbClient),
	}
}

func (hc *HarvestController) HarvestCrop(c *gin.Context) {

	userId := c.GetHeader("user_id")
	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	plantingId := c.Param("plantid")
	plantingObjectId, _ := primitive.ObjectIDFromHex(plantingId)

	var harvestResult *models.HarvestResult

	plantedCrop, err := hc.service.GetPlantedCrop(userObjectId, plantingObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	if err := hc.service.ValidateHarvest(plantedCrop); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHttpError(c, err.Error(), http.StatusBadRequest))
		return
	}

	harvestResult, err = hc.service.CalculateHarvestResult(plantedCrop, 1.0)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	if err := hc.service.MarkCropAsHarvested(plantingObjectId, harvestResult); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	if err := hc.service.FreeLandUnits(context.TODO(), plantedCrop.LandUnitIDs); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	if err := hc.service.AddToWarehouse(context.TODO(), userObjectId, harvestResult); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": harvestResult,
	})
}
