package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chamanbetra/user-management-app/models"
	"github.com/chamanbetra/user-management-app/services"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var validate = validator.New()

// BasicAuth middleware for simple authentication
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, password, hasAuth := r.BasicAuth()

		if !hasAuth {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "http_request", r)

		user, err := services.GetUserByEmail(ctx, email)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid email", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Unauthorized: Invalid password", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if authentication is successful
		next.ServeHTTP(w, r)
	})
}

func sendVerificationEmail(toEmail string, token string) error {
	from := mail.NewEmail("UserManagement", "flatmates.web@gmail.com")
	subject := "Please verify your email address"
	to := mail.NewEmail("New User", toEmail)

	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	plainTextContent := fmt.Sprintf("Please click the following link to verify your email: %s", verificationLink)
	htmlContent := fmt.Sprintf("<p>Please click the following link to verify your email:</p><a href=\"%s\">Verify Email</a>", verificationLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	sendgridAPIKey := os.Getenv("SENDGRID_APIKEY")
	if sendgridAPIKey == "" {
		return fmt.Errorf("SendGrid API key not found in environment variables")
	}
	client := sendgrid.NewSendClient(sendgridAPIKey)

	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send verification email: %v\n", err)
		return err
	}

	log.Printf("Email sent, status code: %d\n", response.StatusCode)
	return nil
}

// CreateUser handles user creation
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		// Return validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	user.VerificationToken = uuid.NewString()
	user.Verified = false
	user.Token_GeneratedTime = time.Now()

	ctx := context.WithValue(r.Context(), "http_request", r)

	if err := services.CreateUser(ctx, &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if err := sendVerificationEmail(user.Email, user.VerificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

func IsUserVerified(ctx context.Context, email string) (bool, error) {
	user, err := services.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if !user.Verified {
		return false, nil
	}

	return true, nil
}

// GetUser handles fetching a user by email
func GetUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := requestBody.Email
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "http_request", r)

	verified, err := IsUserVerified(ctx, email)
	if err != nil || !verified {
		http.Error(w, "User is not verified or not found", http.StatusForbidden)
		return
	}

	user, err := services.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	json.NewEncoder(w).Encode(userResponse)
}

// UpdateUser handles updating a user by email
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		DOB       string `json:"dob"`
	}

	username, _, _ := r.BasicAuth()

	current_email := username

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := requestBody.Email
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "http_request", r)

	verified, err := IsUserVerified(ctx, current_email)
	if err != nil || !verified {
		http.Error(w, "User is not verified or not found", http.StatusForbidden)
		return
	}

	user, err := services.GetUserByEmail(ctx, current_email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if requestBody.FirstName != "" {
		user.FirstName = requestBody.FirstName
	}

	if requestBody.LastName != "" {
		user.LastName = requestBody.LastName
	}

	if requestBody.DOB != "" {
		user.DOB = requestBody.DOB
	}

	if requestBody.Email != "" {
		user.Email = requestBody.Email
	}

	if err := services.UpdateUser(ctx, current_email, user); err != nil {
		log.Println(err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	json.NewEncoder(w).Encode(userResponse)
}

// DeleteUser handles deleting a user by email
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := requestBody.Email
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	ctx := context.WithValue(r.Context(), "http_request", r)

	verified, err := IsUserVerified(ctx, email)
	if err != nil || !verified {
		http.Error(w, "User is not verified or not found", http.StatusForbidden)
		return
	}

	if err := services.DeleteUser(ctx, email); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "There seems to be some error with the Token", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "http_request", r)

	username, err := services.ValidateToken(ctx, token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if err := services.VerifyUserByEmail(ctx, username); err != nil {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}

}
