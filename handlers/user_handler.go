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
	if err := c.ShouldBindJSON(&itf.CreateUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Name:  itf.CreateUserRequest.Name,
		Email: itf.CreateUserRequest.Email,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userJob := models.UserJob{
		UserID:         user.ID,
		JobTitle:       itf.CreateUserRequest.Job,
		JobDescription: itf.CreateUserRequest.JobDesc,
	}

	if err := db.DB.Create(&userJob).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user
	var retrievedUser models.User
	db.DB.Last(&retrievedUser)

	response := itf.CreateUserResponse{
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
