package main

import (
	"os"
	"os/signal"
	"syscall"
	"userService/server"
	"userService/utils"
)

func main() {

	// Load Env Data.
	utils.EnvConfig()

	// Make a channel to capture OS signals.
	exitSignalChannel := make(chan os.Signal, 1)

	// This channel subscribes to SIGINT (Interrupt) and SIGTERM (Terminate) signals for gracefully shutdown of the server.
	signal.Notify(exitSignalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	srv := server.ServerInit()

	//  Gin Router and HTTP Server.
	go srv.Start()

	<-exitSignalChannel

	srv.Stop()

}
