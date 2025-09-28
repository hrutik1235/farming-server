package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/controller"
	middleware "github.com/hrutik1235/farming-server/midlleware"
	"github.com/hrutik1235/farming-server/types"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func NewCropRoutes(r *gin.RouterGroup, conn *grpc.ClientConn, dbClient *mongo.Database) {
	cropController := controller.NewCropController(dbClient)
	group := r.Group("/crop")

	group.Use(middleware.GateValidateUser())
	group.POST("", middleware.ValidateRequest[types.CreateCrop, any, any](), cropController.CreateCrop)
	group.GET("", cropController.GetAllCrops)
	group.POST("/plant/:cropid", middleware.ValidateRequest[types.PlantCrop, any, any](), cropController.PlantCrop)
	group.GET("/plant", cropController.GetAllPlantedCrops)
}
