package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"userService/database"
	"userService/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	MongoClient      *mongo.Client
	httpServer       *http.Server
	databaseFunction database.DatabaseFunction
	ctx              context.Context
}

// Initialise a Server for the Backend services.
func ServerInit() *Server {

	// Connect to the MongoDatabase.
	mongoClient, ctx, err := database.ConnectDB(utils.EnvData["MONGO_DB_URL"])
	if err != nil {
		log.Fatal("Failed to establish client connection to database: ", err)
	}

	// Checking the mongoDatabase is connected Successfully.
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Failed to ping the mongoDb: ", err)
	}

	// create the http server.
	addr := fmt.Sprintf(":%s", utils.EnvData["GIN_PORT"])

	httpServ := &http.Server{
		Addr: addr,
	}

	// Create database helper.
	dbFunction := database.NewDBHelper(mongoClient)

	return &Server{
		MongoClient:      mongoClient,
		httpServer:       httpServ,
		databaseFunction: dbFunction,
		ctx:              ctx,
	}
}

func (srv *Server) Start() {

	/*

		1. We have to test the database connection.
		2. Inject all the REST API routes into the HTTP Server.
		3. Start the HTTP server.

	*/

	err := srv.MongoClient.Ping(srv.ctx, nil)
	if err != nil {
		fmt.Println("Database connection tested Failed: ", err)
	} else {
		fmt.Println("Database connection Tested Successfully.")
	}

	srv.httpServer.Handler = srv.InjectRoutes()

	severErr := srv.httpServer.ListenAndServe()
	if severErr != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting the HTTP server: ", srv.httpServer.Addr)
		return
	}

}

func (srv *Server) Stop() {
	/*
		1.We have to stop the database connection.

			for this we have to change the database code.
			I think we have to make the struct and pass the struct in the Serve struct with mongoClient and the context.

		2.http server shutdown.
		we have to create the context timeout after that we shutdown the server.
	*/

	err := srv.MongoClient.Disconnect(srv.ctx)
	if err != nil {
		fmt.Println("ailed to disconnect to the database: ", err)
	} else {
		fmt.Println("Successfully Database Disconnect")
	}

	// Passing context with timeout for 5 seconds which is enough time for a query to execute.
	// If the query didn't execute in time, it will return a timeout error which can be sent in the response as an error.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		fmt.Println("Error closing the HTTP Server: ", err)
	} else {
		fmt.Println("Closed the HTTP server Successfully")
	}

}
