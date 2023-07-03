package handlers

import (
	"fmt"
	"net/http"

	db "go-backend/database"
	"go-backend/database/models"
	itf "go-backend/interfaces"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateTrip(c *gin.Context) {
	var req itf.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	token := c.GetHeader("token")
	fmt.Println(token)
	session := sessions.Default(c)
	session.Set("token", "hehehe")
	username := session.Get(token)
	test := session.Get("token")
	fmt.Println(username)
	fmt.Println(test)

	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// trip := models.Trip{
	// 	TripName: req.TripName,
	// 	TripCode: "test",
	// 	OwnerId: user.ID,

	// }

	// Validate the user data

}

func SetValueHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Set a value in the session
	session.Set("myKey", "myValue")

	// Save the session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Value set in the session"})
}

func GetValueHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Retrieve the value from the session
	value := session.Get("token")

	// Check if the value exists in the session
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Value not found in the session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
