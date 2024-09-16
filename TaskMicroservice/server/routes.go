package server

import "github.com/gin-gonic/gin"

func (srv *Server) InjectRoutes() *gin.Engine {
	router := gin.Default()

	v2 := router.Group("/api/v2")
	{
		task := v2.Group("/tasks")
		task.GET("/:id", srv.GetTask)
		task.GET("/", srv.GetAllTasks)
		task.POST("/add", srv.CreateTask)
		task.DELETE("/:id", srv.DeleteTask)
		task.GET("/user/:email", srv.GetTasksByUserEmail)
		task.GET("/:id/user", srv.GetUserByTaskID)
	}
	return router
}
