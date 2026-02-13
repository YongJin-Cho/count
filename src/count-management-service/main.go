package main

import (
	"count-management-service/adapters/inbound"
	"count-management-service/adapters/outbound"
	"count-management-service/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	processingSvcURL := os.Getenv("PROCESSING_SERVICE_URL")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPass == "" {
		dbPass = "password"
	}
	if dbName == "" {
		dbName = "countdb"
	}
	if processingSvcURL == "" {
		processingSvcURL = "http://count-processing-service:8080"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v. Using in-memory database for demo/fallback.", err)
		// Fallback to sqlite if postgres is not available (for easier testing in some envs, but AGENTS.md says Postgres)
		// For now, let's just fail if it's a real environment.
		// db, _ = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		log.Fatal(err)
	}

	repo := outbound.NewPostgresRepository(db)
	valueClient := outbound.NewValueServiceClient(processingSvcURL)
	service := domain.NewCountItemUseCase(repo, valueClient)
	handler := inbound.NewHTTPHandler(service)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	// UI Endpoints
	ui := r.Group("/ui")
	{
		ui.POST("/count-items", handler.RegisterItemUI)
		ui.GET("/count-items", handler.ListItemsUI)
		ui.DELETE("/count-items/:id", handler.DeleteItemUI)
		ui.PUT("/counts/:id", handler.UpdateItemUI)
		ui.GET("/counts/:id/value", handler.GetItemValueUI)
	}

	// External API Endpoints
	api := r.Group("/api/v1")
	{
		api.GET("/count-items", handler.ListItemsAPI)
		api.POST("/count-items", handler.RegisterItemAPI)
		api.PUT("/count-items/:id", handler.UpdateItemAPI)
		api.DELETE("/count-items/:id", handler.DeleteItemAPI)
	}

	log.Println("Starting count-management-service on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
