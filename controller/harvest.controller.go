package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/service"
	"github.com/hrutik1235/farming-server/utils"
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
	// var harvestResult models.HarvestResult

	// userId := c.GetHeader("user_id")
	// userIdObjectId, _ := primitive.ObjectIDFromHex(userId)

	// plantingId := c.Params.ByName("plantingid")

	// plantingIdObjectId, _ := primitive.ObjectIDFromHex(plantingId)

	// plantedCrop, err := hc.service.GetPlantedCrop(userIdObjectId, plantingIdObjectId)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
	// 	return
	// }

	landUnitIds := []string{"68d89ec1681d9197594dffd9", "68d89ec1681d9197594dffda", "68d89ec1681d9197594dffdb"}

	landUnits, err := hc.service.GetHarvestUnits(landUnitIds)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": landUnits,
	})

}
