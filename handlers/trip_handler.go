package handlers

import (
	itf "go-backend/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createTrip(c *gin.Context) {
	var req itf.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}
	// Validate the user data

}
