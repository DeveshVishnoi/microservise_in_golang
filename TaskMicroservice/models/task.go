package models

var TaskDataBase = "Task-DataBase"
var TaskCollectioName = "task"

type Task struct {
	Id          string `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	EmailID     string `json:"email_id" bson:"email_id"`
}
