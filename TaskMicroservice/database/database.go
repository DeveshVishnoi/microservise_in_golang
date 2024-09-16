package database

import (
	"context"
	"fmt"
	"log"
	"taskService/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ConnectDB(address string) (*mongo.Client, context.Context, error) {

	// Create an empty background context to keep alive till the application is running.
	ctx := context.Background()

	mongoConn := options.Client().ApplyURI(address)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err == nil {
		log.Printf("Connected to mongo database : %s\n", address)
		return mongoClient, ctx, nil
	} else {
		return nil, nil, nil
	}

}

func (dbHelper *DBHelper) CreateTask(taskInfo models.Task) error {

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := dbHelper.taskCollection.InsertOne(ctx, taskInfo)
	if err != nil {
		fmt.Println("failed to create task")
		return err
	}
	fmt.Println("Successfully task Insertd : ", result.InsertedID)
	return nil

}

func (dbHelper *DBHelper) GetTasksByUserEmail(emailID string) ([]models.Task, error) {
	var tasks []models.Task

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email_id": emailID}
	cursor, err := dbHelper.taskCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func (dbHelper *DBHelper) DeleteTask(taskID string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": taskID}
	_, err := dbHelper.taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
func (dbHelper *DBHelper) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	cursor, err := dbHelper.taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func (dbHelper *DBHelper) GetTaskByID(taskID string) (models.Task, error) {
	var task models.Task

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": taskID}
	err := dbHelper.taskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return task, err
	}

	return task, nil
}
