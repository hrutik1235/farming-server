package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/kafkaconn"
	"github.com/hrutik1235/farming-server/models"
	"github.com/hrutik1235/farming-server/service"
	"github.com/hrutik1235/farming-server/types"
	"github.com/hrutik1235/farming-server/utils"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	service  *service.UserService
	dbClient *mongo.Database
}

func NewUserController(dbClient *mongo.Database) *UserController {
	return &UserController{
		service:  service.NewUserService(dbClient),
		dbClient: dbClient,
	}
}

func (userController *UserController) GetUserDetails(c *gin.Context) {
	userId := c.GetHeader("user_id")

	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	user, err := userController.service.GetUserById(userObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, "User Error", http.StatusInternalServerError))
		return
	}

	wallet, err := userController.service.GetUserWallet(userObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, "Wallet Error", http.StatusInternalServerError))
		return
	}

	landUnits, err := userController.service.GetUserLandUnits(userObjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, "Land Error", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"wallet": wallet,
			"land":   landUnits,
		},
	})

}

func (userController *UserController) RegisterUser(c *gin.Context) {
	body := c.MustGet("body").(types.RegisterUser)

	user := models.User{
		Username:      body.Username,
		DisplayName:   body.Name,
		Email:         body.Email,
		ServerAddress: "localhost:8080",
	}

	prevuser, _ := userController.service.FindUserByCriteria(bson.M{
		"username": body.Username,
		"email":    body.Email,
	})


	if prevuser != nil {
		c.JSON(200, gin.H{
			"message": "User already exists",
			"data":    prevuser.ID,
		})
		return
	}

	savedUser, err := user.Save(userController.dbClient.Collection("users"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	landErr := userController.service.AllocateLandToUser(savedUser.InsertedID.(primitive.ObjectID), user.Username)

	if landErr != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, landErr.Error(), http.StatusInternalServerError))
		return
	}

	kconfig := kafkaconn.NewKafka([]string{"localhost:9092"})

	defer kconfig.Close()

	if err := kconfig.CreateTopic("register"); err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	message := kafka.Message{
		Key:   []byte("register"),
		Value: []byte(fmt.Sprint(body)),
	}

	if err := kconfig.WriteMessage(c, "register", message); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHttpError(c, err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("User registered with name %s", body.Name),
		"data":    savedUser.InsertedID,
	})
}
