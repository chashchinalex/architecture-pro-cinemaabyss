package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"proxy/handlers"
	"strconv"
	"syscall"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

const (
	moviesMigrationPercentEnv = "MOVIES_MIGRATION_PERCENT"
)

func main() {
	migrationParcentRaw := os.Getenv(moviesMigrationPercentEnv)
	migrationPercent, err := strconv.Atoi(migrationParcentRaw)
	if err != nil {
		log.Fatalf("Parse '%s' error. Integer expected\n", moviesMigrationPercentEnv)
	}

	monolithURL := os.Getenv("MONOLITH_URL")
	if _, err = url.Parse(monolithURL); err != nil || monolithURL == "" {
		log.Fatalf("Invalid monolith URL: %s\n", monolithURL)
	}
	moviesURL := os.Getenv("MOVIES_SERVICE_URL")
	if _, err = url.Parse(moviesURL); err != nil || moviesURL == "" {
		log.Fatalf("Invalid movies URL: %s\n", moviesURL)
	}

	router := gin.Default()

	apiGroup := router.Group("/api")

	{
		proxyHandler := handlers.NewMoviesHandler(monolithURL, moviesURL, migrationPercent)
		proxyHandler.RegisterRoutes(apiGroup)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8000"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return ":" + value
}
