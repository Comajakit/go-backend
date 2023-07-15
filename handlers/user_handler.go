package handlers

import (
	"net/http"

	"go-backend/database/models"
	itf "go-backend/interfaces"
	util "go-backend/utils"

	db "go-backend/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	hashedPassword, err := util.HashPassword(req.Password)
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

func UserLogin(c *gin.Context) {
	var req itf.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	result, err := util.ValidatePassword(req.Username, req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result {
		// Generate a new access token
		accessToken := uuid.New().String()

		// Create session
		session := sessions.Default(c)
		session.Set("token", "test")
		session.Set(accessToken, req.Username)

		// Set session expiration
		// expiration := 3 * 60 * 60 // Default session expiration is 3 hours
		// if req.Forever {
		// 	expiration = 30 * 24 * 60 * 60 // Set session to expire in 30 days (forever)
		// }

		// session.Options(sessions.Options{
		// 	MaxAge:   expiration,
		// 	HttpOnly: true,
		// 	Secure:   false, // Set to true if using HTTPS
		// })
		// rm := session.Get("token")
		// fmt.Println(rm)
		// rm2 := session.Get(accessToken)
		// fmt.Println("this is token: " + rm2.(string))

		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save session"})
			return
		}

		// Return success response with the access token
		c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
	} else {
		// Return failure response
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
	}
}

// func UserLogin(c *gin.Context) {
// 	var req itf.UserLoginRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
// 		return
// 	}

// 	result, err := util.ValidatePassword(req.Username, req.Password)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if result {
// 		// Create session
// 		session := sessions.Default(c)
// 		session.Set(req.Username, req.Username)

// 		// Set session expiration
// 		expiration := 3 * 60 * 60 // Default session expiration is 3 hours
// 		if req.Forever {
// 			expiration = 30 * 24 * 60 * 60 // Set session to expire in 30 days (forever)
// 		}

// 		session.Options(sessions.Options{
// 			MaxAge:   expiration,
// 			HttpOnly: true,
// 			Secure:   true, // Set to true if using HTTPS
// 		})

// 		if err := session.Save(); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save session"})
// 			return
// 		}

// 		// Return success response
// 		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
// 	} else {
// 		// Return failure response
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
// 	}

// }

func ProtectedRoute(c *gin.Context) {
	// Check if the user is authenticated
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Proceed with handling the protected resource
	c.JSON(http.StatusOK, gin.H{"message": "Protected resource"})
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
