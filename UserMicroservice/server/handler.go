package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"userService/models"
	"userService/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (srv *Server) CreateUser(c *gin.Context) {

	var userInfo models.User

	err := json.NewDecoder(c.Request.Body).Decode(&userInfo)
	if err != nil {
		fmt.Println("Error getting decode the request data into struct format: ", err)
		return
	}

	fmt.Println("User Data : ", userInfo)

	if userInfo.EmailId == "" {
		c.JSON(http.StatusBadRequest, "email Id can not be empty")
		return
	}

	if userInfo.Password == "" {
		c.JSON(http.StatusBadRequest, "password can not be empty")
		return
	}

	if userInfo.Name == "" {
		c.JSON(http.StatusBadRequest, "Name can not be empty")
		return
	}

	if userInfo.Department == "" {
		userInfo.Department = "GLOBAL"
	}
	userInfo.Id = uuid.New().String()

	err = srv.databaseFunction.CreateUser(userInfo)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		if err.Error() == models.UserAlreadyExist {
			c.JSON(http.StatusConflict, gin.H{"Message": models.UserAlreadyExist})
			return
		}

		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"Message": "Success"})
		return
	}
}

func (srv *Server) UpdateUser(c *gin.Context) {

	var userInfo models.User

	err := json.NewDecoder(c.Request.Body).Decode(&userInfo)
	if err != nil {
		fmt.Println("Error getting decode the request data into struct format: ", err)
		return
	}

	fmt.Println("User Data : ", userInfo)

	if userInfo.EmailId == "" {
		c.JSON(http.StatusBadRequest, "email Id can not be empty")
		return
	}

	if userInfo.Password == "" {
		c.JSON(http.StatusBadRequest, "password can not be empty")
		return
	}

	if userInfo.Name == "" {
		c.JSON(http.StatusBadRequest, "Name can not be empty")
		return
	}

	if userInfo.Department == "" {
		userInfo.Department = "GLOBAL"
	}
	userInfo.Id = uuid.New().String()

	err = srv.databaseFunction.UpdateUser(userInfo)
	if err != nil {
		if err.Error() == models.UserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"Message": models.UserNotFound,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Got Some error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully Update!!!",
	})

}

func (srv *Server) DeleteUser(c *gin.Context) {

	email_id := c.Param("email")

	// Validate email format
	if err := utils.ValidateEmail(email_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := srv.databaseFunction.DeleteUser(email_id)
	if err != nil {
		if err.Error() == models.UserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"Message": models.UserNotFound,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Got Some error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (srv *Server) GetAllUser(c *gin.Context) {

	users, err := srv.databaseFunction.GetAllUser()
	if err != nil {
		fmt.Println("got some error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": users,
	})
	return

}

func (srv *Server) GetUser(c *gin.Context) {

	email_id := c.Param("email")

	// Validate email format
	if err := utils.ValidateEmail(email_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := srv.databaseFunction.GetUser(email_id)
	if err != nil {
		if err.Error() == models.UserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"Message": models.UserNotFound,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Got Some error",
			})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"Message": user,
	})

}
