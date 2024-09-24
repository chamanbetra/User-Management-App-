package services

import (
	"errors"

	"github.com/chamanbetra/user-management-app/database"
	"github.com/chamanbetra/user-management-app/models"
)

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	// Here, we should check if the user already exists based on email
	var existingUser models.User
	if err := database.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	// Create the user
	if err := database.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user's details by email
func UpdateUser(email string, updatedUser *models.User) error {
	var user models.User
	if err := database.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Update user fields
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.DOB = updatedUser.DOB

	if err := database.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user from the database by email
func DeleteUser(email string) error {
	if err := database.Where("email = ?", email).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}
