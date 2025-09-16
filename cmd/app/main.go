package main

import (
	"context"
	"log"
	"os"
	"price-monitoring-server/internal/config"
	"price-monitoring-server/internal/monitoring"
	"price-monitoring-server/internal/repository"
	"price-monitoring-server/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var jsonPath string
var startTime time.Time

func init() {
	jsonPath = "config/config.json"
	startTime = time.Now()
}

func main() {
	// Use a background context that can be cancelled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Database Connection ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL is not set, using default value")
		dbURL = "postgres://user:password@localhost:5433/prices?sslmode=disable"
	}

	dbpool, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	priceRepo := repository.NewPriceRepository(dbpool)
	if err := priceRepo.InitSchema(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v\n", err)
	}

	// --- Config Loading ---
	cfg, err := config.LoadData(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	// --- Service and Monitoring Setup ---
	priceService := service.NewPriceService(priceRepo)

	updates := monitoring.Run(ctx, cfg)

	go monitoring.ConsumeAndAggregate(updates, priceService)

	// --- Gin Router Setup ---
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// API endpoints
	router.GET("/prices", func(c *gin.Context) {
		c.JSON(200, priceService.GetLatestPrices())
	})

	router.GET("/stats", func(c *gin.Context) {
		c.JSON(200, priceService.GetStats())
	})

	router.GET("/history/:productName", func(c *gin.Context) {
		productName := c.Param("productName")
		history, err := priceService.GetHistory(c.Request.Context(), productName)
		if err != nil {
			log.Printf("Error getting history for %s: %v", productName, err)
			c.JSON(500, gin.H{"error": "failed to retrieve price history"})
			return
		}
		c.JSON(200, history)
	})

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
