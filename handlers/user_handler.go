package handlers

import (
	"net/http"
	db "pokemon/database"
	"pokemon/database/models"
	itf "pokemon/interfaces"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	// Parse the request body
	var req itf.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}
	// Validate the user data
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required fields"})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword
	// Create user in database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created user
	c.JSON(http.StatusCreated, user)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func validatePassword(username string, password string) {

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
