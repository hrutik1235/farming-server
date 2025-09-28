package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/controller"
	middleware "github.com/hrutik1235/farming-server/midlleware"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func NewHarvestRoutes(r *gin.RouterGroup, conn *grpc.ClientConn, dbClient *mongo.Database) {
	harvestController := controller.NewHarvestController(dbClient)

	group := r.Group("/harvest")

	group.Use(middleware.GateValidateUser())
	group.POST("/:plantid", middleware.ValidateRequest[any, any, any](), harvestController.HarvestCrop)
}
