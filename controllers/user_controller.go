package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chamanbetra/user-management-app/models"
	"github.com/chamanbetra/user-management-app/services"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

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

	ctx := context.WithValue(r.Context(), "http_request", r)

	if err := services.CreateUser(ctx, &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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

	user, err := services.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
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

	json.NewEncoder(w).Encode(user)
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

	if err := services.DeleteUser(ctx, email); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
