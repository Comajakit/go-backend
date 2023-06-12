package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "pokemon/database"
	"pokemon/database/models"
	itf "pokemon/interfaces"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Calling Created User")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the user to the database or perform any necessary operations
	fmt.Println(user)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteRecentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the latest user ID
	var latestUser models.User
	err := db.DB.Order("id desc").First(&latestUser).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the latest user directly by ID
	err = db.DB.Delete(&models.User{}, latestUser.ID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response indicating successful deletion
	response := map[string]interface{}{
		"message": "Latest user deleted successfully",
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response map as JSON and write it to the response body
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
