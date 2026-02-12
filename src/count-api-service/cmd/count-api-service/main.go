package main

import (
	"context"
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/event"
	"count-api-service/internal/component/collector"
	"count-api-service/internal/component/storage"
	"count-api-service/internal/component/ui"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	// UI Module
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = "internal/component/ui/templates"
	}
	uiHandler, err := ui.NewUIHandler(collectorHandler, store, templateDir, authProvider)
	if err != nil {
		log.Fatalf("Failed to initialize UI handler: %v", err)
	}

	// 2. Orchestration (Routes)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/v1/collect", collectorHandler.CollectCount)
	mux.HandleFunc("/api/v1/counts", collectorHandler.GetCounts)

	// UI Routes
	mux.HandleFunc("/ui/counts", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ui/counts" {
			if r.Method == http.MethodGet {
				uiHandler.GetCountList(w, r)
			} else if r.Method == http.MethodPost {
				uiHandler.CreateCount(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if r.URL.Path == "/ui/counts/new" {
			uiHandler.GetCreateForm(w, r)
		} else {
			// Path with source_id
			if strings.HasSuffix(r.URL.Path, "/increment") {
				uiHandler.IncrementCount(w, r)
			} else if strings.HasSuffix(r.URL.Path, "/decrement") {
				uiHandler.DecrementCount(w, r)
			} else if strings.HasSuffix(r.URL.Path, "/edit") {
				uiHandler.GetEditForm(w, r)
			} else if r.Method == http.MethodGet {
				uiHandler.GetCountRow(w, r)
			} else if r.Method == http.MethodPut {
				uiHandler.UpdateCount(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}
	})

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
