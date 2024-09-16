package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"userService/models"

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

func (dbHelper *DBHelper) GetUser(email_id string) (*models.User, error) {

	var user models.User

	filter := bson.M{"email_id": email_id}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := dbHelper.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No emailId register of this email")
			return &user, err
		}
		fmt.Println("Failed to read user databases on the emailID")
	}

	return &user, nil
}

func (dbHelper *DBHelper) CreateUser(user models.User) error {

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	data, err := dbHelper.GetUser(user.EmailId)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			result, err := dbHelper.userCollection.InsertOne(ctx, user)
			if err != nil {
				fmt.Println("failed to create user")
				return err
			}
			fmt.Println("Successfully User Insertd : ", result.InsertedID)
			return nil

		} else {

			// An error other than "no documents" occurred
			fmt.Println("failed to create user")
			return err
		}
	}

	if data != nil {
		fmt.Println("User is already exist")
		return errors.New(models.UserAlreadyExist)
	}

	return nil
}

func (dbHelper *DBHelper) UpdateUser(userInfo models.User) error {

	// Filter to find the user by email_id
	filter := bson.M{"email_id": userInfo.EmailId}

	// Create an update document with the fields you want to update
	update := bson.M{
		"$set": bson.M{
			"name":       userInfo.Name,
			"password":   userInfo.Password,
			"department": userInfo.Department,
			"updated_at": time.Now(), // Assuming you have an updated_at field
		},
	}

	// Set a context with a timeout for the operation
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Perform the update operation
	result, err := dbHelper.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("Failed to update user:", err)
		return err
	}

	if result.MatchedCount == 0 {
		fmt.Println("No user found with the given email_id.")
		return errors.New(models.UserNotFound)
	}

	fmt.Println("Successfully updated user:", userInfo.EmailId)
	return nil
}

func (dbHelper *DBHelper) DeleteUser(email string) error {

	filter := bson.M{"email_id": email}

	// Set a context with a timeout for the operation
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := dbHelper.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Failed to delete the user : ", err)
		return err
	}
	if result.DeletedCount == 0 {
		fmt.Println("No user found with the given email_id.")
		return errors.New(models.UserNotFound)
	}

	return nil

}

func (dbHelper *DBHelper) GetAllUser() ([]models.User, error) {
	var users []models.User

	filter := bson.M{}

	// Set a context with a timeout for the operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a cursor from the Find operation
	cursor, err := dbHelper.userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check if there were any errors during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
