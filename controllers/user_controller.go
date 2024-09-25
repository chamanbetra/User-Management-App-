package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/chamanbetra/user-management-app/models"
	"github.com/chamanbetra/user-management-app/services"
	"github.com/go-playground/validator"
)

var validate = validator.New()

// BasicAuth middleware for simple authentication
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, hasAuth := r.BasicAuth()

		if !hasAuth || username != "admin" || password != "password" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
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

	if err := services.CreateUser(&user); err != nil {
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

	user, err := services.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles updating a user by email
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string      `json:"email"`
		User  models.User `json:"user"`
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

	if err := services.UpdateUser(email, &requestBody.User); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requestBody.User)
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

	if err := services.DeleteUser(email); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
