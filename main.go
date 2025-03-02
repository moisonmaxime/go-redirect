package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	redirectServer, err := newRedirectServer("localhost", 8080, "urls.json")

	if err != nil {
		panic(err)
	}

	// Channel for handling graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		redirectServer.start()
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Initiate a graceful shutdown with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redirectServer.stop(shutdownCtx)

	log.Println("Server gracefully stopped")
}
