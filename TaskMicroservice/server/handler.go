package server

import (
	"net/http"
	"taskService/models"
	"taskService/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTask handles the creation of a new task.
func (srv *Server) CreateTask(c *gin.Context) {
	var taskInfo models.Task

	taskInfo.Id = uuid.New().String()
	// Decode request body into taskInfo struct
	if err := c.ShouldBindJSON(&taskInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate email format
	if err := utils.ValidateEmail(taskInfo.EmailID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid email format", "details": err.Error()})
		return
	}

	// Check if user exists in UserMicroservice
	_, isUserExist, err := CheckUserExists(taskInfo.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error checking user existence", "details": err.Error()})
		return
	}
	if !isUserExist {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User does not exist"})
		return
	}

	// Insert the task into the database
	if err := srv.databaseFunction.CreateTask(taskInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create task", "details": err.Error()})
		return
	}

	// Successfully created the task
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Task created successfully",
		"task":    taskInfo,
	})
}

// GetUserByTaskID retrieves the user associated with a task.
func (srv *Server) GetUserByTaskID(c *gin.Context) {
	taskID := c.Param("id")

	task, err := srv.databaseFunction.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Task not found", "details": err.Error()})
		return
	}

	user, _, err := CheckUserExists(task.EmailID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
}

// GetTasksByUserEmail retrieves tasks for a specific user email.
func (srv *Server) GetTasksByUserEmail(c *gin.Context) {
	emailID := c.Param("email")

	tasks, err := srv.databaseFunction.GetTasksByUserEmail(emailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "tasks": tasks})
}

// GetTask retrieves a specific task by its ID.
func (srv *Server) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	task, err := srv.databaseFunction.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Task not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "task": task})
}

// DeleteTask deletes a specific task by its ID.
func (srv *Server) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	if err := srv.databaseFunction.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Task not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Task deleted successfully"})
}

// GetAllTasks retrieves all tasks.
func (srv *Server) GetAllTasks(c *gin.Context) {
	tasks, err := srv.databaseFunction.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "tasks": tasks})
}
