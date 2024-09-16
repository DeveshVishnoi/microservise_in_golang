package database

import (
	"userService/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseFunction interface {
	CreateUser(models.User) error
	GetUser(string) (*models.User, error)
	UpdateUser(models.User) error
	DeleteUser(emai_id string) error
	GetAllUser() ([]models.User, error)
}

type DBHelper struct {
	MongoDBClient  *mongo.Client
	userCollection *mongo.Collection
}

func NewDBHelper(mongoClient *mongo.Client) DatabaseFunction {
	return &DBHelper{
		userCollection: mongoClient.Database(models.UserDataBase).Collection(models.UserCollectioName),
	}
}
