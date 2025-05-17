package handler

import (
	"log"
	"net/http"

	"github.com/KirillKotovsky/location_user_api/internal/model"
	"github.com/KirillKotovsky/location_user_api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleLocationUpdate(c *gin.Context) {
	var req model.LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Handling location update for user: %s", req.Username)
	if err := h.service.UpdateLocation(c, req); err != nil {
		log.Printf("Failed to update location: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) HandleNearbyUsers(c *gin.Context) {
	var req model.NearbyUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Invalid query parameters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Searching for users near lat: %.6f, lon: %.6f, radius: %.2f",
		req.Latitude, req.Longitude, req.Radius)
	users, err := h.service.FindNearbyUsers(c, req)
	if err != nil {
		log.Printf("Failed to find nearby users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
