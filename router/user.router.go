package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/controller"
	middleware "github.com/hrutik1235/farming-server/midlleware"
	"github.com/hrutik1235/farming-server/types"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func NewUserRoutes(r *gin.RouterGroup, conn *grpc.ClientConn, dbClient *mongo.Database) {
	userController := controller.NewUserController(dbClient)
	group := r.Group("/user")

	group.POST("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})


	group.POST("/register", middleware.ValidateRequest[types.RegisterUser, types.AuthHeader, any](), userController.RegisterUser)

	group.Use(middleware.GateValidateUser())

	group.GET("", userController.GetUserDetails)
}
