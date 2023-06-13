package handlers

import (
	"fmt"
	"net/http"
	db "pokemon/database"
	"pokemon/database/models"
	itf "pokemon/interfaces"

	"github.com/gin-gonic/gin"
)

func CreateUserHandlerGin(c *gin.Context) {
	fmt.Println("Calling Create User Gin")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the user to the database or perform any necessary operations
	// fmt.Println(user)

	db.DB.Create(&user)

	// Retrieve the user
	var retrievedUser models.User
	db.DB.Last(&retrievedUser)

	response := itf.UserResponse{
		ID:    retrievedUser.ID,
		Name:  retrievedUser.Name,
		Email: retrievedUser.Email,
	}

	// Return the created user as a response
	c.JSON(http.StatusOK, response)
}

func DeleteRecentUserHandlerGin(c *gin.Context) {
	// Retrieve the latest user ID
	var latestUser models.User
	err := db.DB.Order("id desc").First(&latestUser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the latest user directly by ID
	err = db.DB.Delete(&models.User{}, latestUser.ID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a response indicating successful deletion
	response := map[string]interface{}{
		"message": "Latest user deleted successfully",
	}

	// Set the Content-Type header to application/json
	c.JSON(http.StatusOK, response)
}
