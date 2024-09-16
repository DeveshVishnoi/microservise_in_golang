package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func forwardRequest(c *gin.Context, serviceURL string) {

	fmt.Println("C ", c)
	req, err := http.NewRequest(c.Request.Method, serviceURL, c.Request.Body)
	if err != nil {
		log.Println("Error creating request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forwarding request"})
		return
	}

	fmt.Println("Request : ", c.Request.Body)
	// Copy headers from the original request
	req.Header = c.Request.Header

	// Use http client to send the request to the target microservice
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error forwarding request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forwarding request"})
		return
	}
	defer resp.Body.Close()

	fmt.Println("Rep.Body", resp.Body, resp.Request.Body)
	// Copy the response back to the client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	// Set the status code and return the response body
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Route for UserService
	router.Any("/api/v1/*path", func(c *gin.Context) {
		serviceURL := "http://localhost:9091" + c.Request.URL.Path

		fmt.Println("Data Py Params : ", c.Request.URL)
		forwardRequest(c, serviceURL)
	})

	// Route for TaskService
	router.Any("/api/v2/*path", func(c *gin.Context) {
		serviceURL := "http://localhost:9092" + c.Request.URL.Path // Point to TaskService
		fmt.Println("Data Py Params : ", c.Request.URL)
		forwardRequest(c, serviceURL)
	})

	return router
}

func main() {
	router := SetupRouter()
	router.Run(":9000") // API Gateway running on port 8080
}
