package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"events/internal/adapter/broker"
	httpHandler "events/internal/adapter/handler"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

const (
	KafkaServer = "localhost:9092"
	KafkaTopic  = "orders-v1-topic"
)

func listenEvents(evChan chan kafka.Event) {
	for ev := range evChan {
		log.Printf("Producer event: %v", ev)
	}
}

func main() {
	kafkaURL := os.Getenv("KAFKA_BROKERS")
	if _, err := url.Parse(kafkaURL); err != nil || kafkaURL == "" {
		log.Fatalf("Invalid kafka brokers URL: %s\n", kafkaURL)
	}

	provider, err := broker.NewKafkaEventProducer(kafkaURL)
	if err != nil {
		log.Fatalf("Create kafka producer error: %v", err)
	}
	defer provider.Close()

	router := gin.Default()

	apiGroup := router.Group("/api")

	{
		eventsHandler := httpHandler.NewEventsHandler(provider)
		eventsHandler.RegisterRoutes(apiGroup)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8082"),
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
