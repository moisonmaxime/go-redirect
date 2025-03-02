package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	host := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	filename := os.Args[3]

	redirectServer, err := newRedirectServer(host, port, filename)

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
