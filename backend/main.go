package main

import (
	"context"
	"forum/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Initialize the database
	if err := application.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer application.CloseDB()

	// Server setup
	srv := &http.Server{
		Addr:    ":8080",
		Handler: application.Router(),
	}

	// WaitGroup for WebSocket server shutdown
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Server started on http://localhost%s", srv.Addr)
	<-stop

	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	wg.Wait()
	log.Println("Server gracefully stopped")
}
