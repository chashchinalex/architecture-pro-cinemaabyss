package handlers

import (
	models "events/core/domain"
	"events/core/ports"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	defaultUserKafkaTopic    = "users-topic"
	defaultPaymentKafkaTopic = "payment-topic"
	defaultMovieKafkaTopic   = "movie-topic"
)

var (
	emptyKey = []byte{}
)

type EventsHandler struct {
	producer ports.EventProducer
}

func NewEventsHandler(producer ports.EventProducer) *EventsHandler {
	return &EventsHandler{
		producer: producer,
	}
}

func (h *EventsHandler) RegisterRoutes(router *gin.RouterGroup) {
	movies := router.Group("/events")
	{
		movies.POST("/users", h.CreateUser)
	}
}

func (h *EventsHandler) CreateUser(c *gin.Context) {
	var userEvent models.UserEvent

	if err := c.ShouldBindJSON(&userEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.producer.SendEvent(c.Request.Context(), defaultUserKafkaTopic, emptyKey, userEvent.EncodeToBinary()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Event produced: CreateUser")
}
