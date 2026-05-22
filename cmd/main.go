package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/jendrix/worldcup-stats-service/config"
	v1 "github.com/jendrix/worldcup-stats-service/internal/handler/v1"
	"github.com/jendrix/worldcup-stats-service/internal/middleware"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

func main() {
	// Load .env file if it exists (ignored in production)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Wire up dependencies: repository → service → handler
	confederationRepo := repository.NewConfederationRepository(pool)
	confederationSvc := service.NewConfederationService(confederationRepo)
	confederationHandlerV1 := v1.NewConfederationHandler(confederationSvc)
	nationalTeamRepo := repository.NewNationalTeamRepository(pool)
	nationalTeamSvc := service.NewNationalTeamService(nationalTeamRepo)
	nationalTeamHandlerV1 := v1.NewNationalTeamHandler(nationalTeamSvc)

	// Set up Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes with header-based versioning
	api := router.Group("/api", middleware.Versioning())
	{
		// v1 routes
		v1Group := api.Group("", middleware.RequireVersion(1))
		confederationHandlerV1.RegisterRoutes(v1Group)
		nationalTeamHandlerV1.RegisterRoutes(v1Group)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
