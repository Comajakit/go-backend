package utils

import (
	"fmt"

	db "go-backend/database"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func GenerateUUID() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
func ValidatePassword(username string, password string) (bool, error) {
	hashedPassword, err := db.GetHashedPassword(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == nil {
		// Password is valid
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(err)
		// Password is invalid
		return false, nil
	} else {
		// An error occurred during password comparison
		return false, err
	}

}
