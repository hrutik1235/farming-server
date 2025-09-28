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

type UserService struct {
	Client *mongo.Database
}

func NewUserService(client *mongo.Database) *UserService {
	return &UserService{
		Client: client,
	}
}

func (u *UserService) GetUserLand(userId primitive.ObjectID) (*models.Land, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var land models.Land

	err := u.Client.Collection(utils.LandsCollection).FindOne(ctx, models.Land{User: userId}).Decode(&land)

	return &land, err
}

func (u *UserService) GetUserLandUnits(userId primitive.ObjectID) ([]models.LandUnit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var landUnits []models.LandUnit

	cursor, err := u.Client.Collection(utils.LandUnitsCollection).Find(ctx, bson.M{"owner_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &landUnits)

	fmt.Println("LAND UNITS: ", len(landUnits))

	if err != nil {
		return nil, err
	}
	return landUnits, nil
}

func (u *UserService) FindUserByCriteria(filter bson.M) (*models.User, error) {
	var user models.User

	err := u.Client.Collection(utils.UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found")
		}
	}

	return &user, nil
}

func (u *UserService) GetUserById(userId primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := u.Client.Collection(utils.UsersCollection).FindOne(ctx, bson.M{"_id": userId}).Decode(&user)

	return &user, err
}

func (u *UserService) GetUserWallet(userId primitive.ObjectID) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wallet models.Wallet

	err := u.Client.Collection(utils.WalletsCollection).FindOne(ctx, bson.M{"user_id": userId}).Decode(&wallet)

	fmt.Println("===== WALLET: ", wallet)

	return &wallet, err
}

func (u *UserService) DeductPlantingCost(userId primitive.ObjectID, cost float64) error {
	if cost <= 0 {
		return fmt.Errorf("invalid cost amount: must be positive")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := u.Client.Collection(utils.WalletsCollection)

	// First check if user has sufficient balance
	var wallet models.Wallet

	err := collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("wallet not found for user")
		}
		fmt.Printf("Error finding wallet: %v\n", err)
		return err
	}

	fmt.Println(" ===== WALLET: ", wallet)

	if wallet.Balance < cost {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", wallet.Balance, cost)
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"user_id": userId,
			"balance": bson.M{"$gte": cost},
		},
		bson.M{
			"$inc": bson.M{
				"balance":     -cost,
				"total_spent": cost,
			},
			"$set": bson.M{
				"last_updated": time.Now(),
			},
		},
	)

	if err != nil {
		fmt.Printf("Error updating wallet: %v\n", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("unable to deduct cost: insufficient balance or wallet not found")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("wallet was matched but not modified")
	}

	fmt.Printf("Successfully deducted %.2f from user %s\n", cost, userId.Hex())
	return nil
}

func (u *UserService) AllocateLandToUser(userId primitive.ObjectID, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userLand, err := u.GetUserLand(userId)

	_, walletErr := u.Client.Collection(utils.WalletsCollection).InsertOne(ctx, models.Wallet{UserID: userId, Balance: 100})

	if walletErr != nil {
		return walletErr
	}

	if err == nil && userLand != nil {
		return nil
	}

	landUnits := make([]interface{}, 100)

	for i := 0; i < 100; i++ {
		landUnit := models.LandUnit{
			Land:        fmt.Sprintf("%s_land", username),
			OwnerID:     userId,
			SizeUnits:   1,
			IsLeased:    false,
			Position:    i + 1,
			IsAvailable: true,
		}
		landUnits[i] = landUnit
	}

	_, lUnitErr := u.Client.Collection(utils.LandUnitsCollection).InsertMany(ctx, landUnits)

	return lUnitErr
}
