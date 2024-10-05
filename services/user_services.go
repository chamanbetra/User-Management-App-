package services

import (
	"context"
	"errors"
	"time"

	"github.com/chamanbetra/user-management-app/database"
	"github.com/chamanbetra/user-management-app/models"
)

// CreateUser creates a new user in the database
func CreateUser(ctx context.Context, user *models.User) error {
	// Here, we should check if the user already exists based on email
	var existingUser models.User
	if err := database.DB.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	// Create the user
	if err := database.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user's details by email
func UpdateUser(ctx context.Context, email string, updatedUser *models.User) error {
	var user models.User
	if err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Update user fields
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.DOB = updatedUser.DOB
	user.Email = updatedUser.Email

	if err := database.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user from the database by email
func DeleteUser(ctx context.Context, email string) error {
	if err := database.DB.WithContext(ctx).Where("email = ?", email).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}

func ValidateToken(ctx context.Context, token string) (string, error) {
	var user models.User

	if err := database.DB.WithContext(ctx).Where("verification_token = ?", token).First(&user).Error; err != nil {
		return "", errors.New("invalid token")
	}

	if time.Since(user.Token_GeneratedTime) > 5*time.Minute {
		return "", errors.New("token expired")
	}

	return user.Email, nil

}

func VerifyUserByEmail(ctx context.Context, email string) error {
	var user models.User

	if err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.Verified = true

	if err := database.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return errors.New("failed to update user verification status")
	}

	return nil

}
