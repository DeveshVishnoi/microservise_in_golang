package models

type User struct {
	Id         string `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
	EmailId    string `json:"email_id" bson:"email_id"`
	Password   string `json:"password" bson:"password"`
}

var UserDataBase = "User-Movie-DataBase"
var UserCollectioName = "user"

const UserAlreadyExist string = "user is already exist"
const UserNotFound string = "user not found"
