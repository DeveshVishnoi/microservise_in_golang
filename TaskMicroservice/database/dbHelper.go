package database

import (
	"taskService/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseFunction interface {
	CreateTask(models.Task) error
	GetTasksByUserEmail(string) ([]models.Task, error)
	DeleteTask(string) error
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(string) (models.Task, error)
}

type DBHelper struct {
	MongoDBClient  *mongo.Client
	taskCollection *mongo.Collection
}

func NewDBHelper(client *mongo.Client) DatabaseFunction {

	return &DBHelper{
		MongoDBClient:  client,
		taskCollection: client.Database(models.TaskDataBase).Collection(models.TaskCollectioName),
	}
}
