package server

import "github.com/gin-gonic/gin"

func (srv *Server) InjectRoutes() *gin.Engine {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users")
		user.POST("/add", srv.CreateUser)
		user.PUT("/add", srv.UpdateUser)
		user.DELETE("/delete/:email", srv.DeleteUser)
		user.GET("/:email", srv.GetUser)
		user.GET("/all", srv.GetAllUser)
	}

	return router
}
