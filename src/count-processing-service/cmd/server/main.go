package main

import (
	"context"
	"count-processing-service/adapters/inbound"
	"count-processing-service/adapters/outbound"
	"count-processing-service/domain"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := outbound.NewPostgresRepository(db)
	historyRepo := outbound.NewPostgresHistoryRepository(db)
	if err := repo.Init(context.Background()); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	useCase := domain.NewCountValueUseCase(repo, historyRepo)
	handler := inbound.NewCountValueHandler(useCase)

	r := gin.Default()
	handler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
