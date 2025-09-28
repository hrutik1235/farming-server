package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/router"
	"github.com/hrutik1235/farming-server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initializeApp(rg *gin.RouterGroup, conn *grpc.ClientConn) {
	client, dbErr := utils.ConnectToDB("mongodb://localhost:27017/")

	if dbErr != nil {
		fmt.Println("Error connecting DB", dbErr.Error())
	}

	db := client.Database("gfarming")

	router.NewUserRoutes(rg, conn, db)
	router.NewCropRoutes(rg, conn, db)
	router.NewHarvestRoutes(rg, conn, db)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	serverGRPC := os.Getenv("SERVER_GRPC")

	conn, err := grpc.NewClient(serverGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic("Failed to connect to gRPC server: " + err.Error())
	}

	router := gin.Default()

	group := router.Group("/api/v1")

	group.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "OK"})
	})

	initializeApp(group, conn)

	router.Run(":" + port)
}
