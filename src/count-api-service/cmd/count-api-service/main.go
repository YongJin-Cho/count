package main

import (
	"context"
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/event"
	"count-api-service/internal/component/collector"
	"count-api-service/internal/component/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting Count API Service...")

	// 1. Initialization
	// Event Module (Internal Bus)
	bus := event.NewEventBus()

	// Storage Module
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "counts.log"
	}
	store := storage.NewFileStorage(storagePath)
	store.Start(bus)

	// Auth Module
	authProvider := auth.NewAuthProvider()

	// Collector Module
	collectorHandler := collector.NewCollectorHandler(authProvider, bus, store)

	// 2. Orchestration (Routes)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/v1/collect", collectorHandler.CollectCount)
	mux.HandleFunc("/api/v1/counts", collectorHandler.GetCounts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// 3. Lifecycle (Startup and Graceful Shutdown)
	go func() {
		log.Printf("HTTP Server listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Service stopped.")
}
