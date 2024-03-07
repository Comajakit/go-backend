package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

func NameFromToken(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	fmt.Println(authHeader)

	// Check if the Authorization header is empty
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
		return
	}

	// Extract the token from the Authorization header
	// Assuming the token is in the format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	token := tokenParts[1]

	// Use the default session to retrieve the access token
	session := sessions.Default(c)
	// Retrieve the session name from the session data

	// Retrieve the access token from the session
	accessTokenInterface := session.Get(token)

	// Check if the token exists in the session
	if accessTokenInterface == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in session"})
		return
	}

	// Perform type assertion to retrieve the username as a string
	accessTokenFromSession, ok := accessTokenInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid access token type"})
		return
	}

	// Respond with HTTP OK and JSON body containing the username
	c.JSON(http.StatusOK, gin.H{"username": accessTokenFromSession})
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

		// Create the default session
		session := sessions.Default(c)
		session.Options(sessions.Options{Path: "/", MaxAge: 3600}) // Adjust options as needed

		// Set the custom session name as part of the session data
		session.Set("session_name", "accesstokensession")
		// Store the access token in the session
		session.Set(accessToken, req.Username)
		fmt.Println(accessToken)

		// Save the session
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
	token := c.Request.Header
	username := session.Get(token.Get("token"))
	sessionID := session
	fmt.Println(sessionID)
	rm := session.Get("token")
	fmt.Println(rm)
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
