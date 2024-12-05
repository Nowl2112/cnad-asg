package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/model"
	"user-service/service"
	"github.com/gorilla/mux"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    createdUser, err := service.RegisterUser(user)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdUser)
}

// Get user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) // Convert string ID to int
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Login user
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	valid, user, err := service.Login(credentials.Email, credentials.Password)
	if err != nil || !valid {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Respond with a success message and user ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user_id": user.ID,  // Return the user_id to store in the frontend
	})
}


// Update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) // Convert string ID to int
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = service.UpdateUser(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func GetRentalHistoryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userID, err := strconv.Atoi(params["id"]) // Use mux.Vars
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    history, err := service.GetRentalHistoryWithVehicle(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(history)
}


func VerifyEmail(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "Invalid verification token", http.StatusBadRequest)
        return
    }

    db := service.GetDB() // Get the database connection from the service package

    query := `UPDATE users 
              SET is_verified = true, verification_token = NULL, token_expiry = NULL 
              WHERE verification_token = ? AND token_expiry > NOW()`
    result, err := db.Exec(query, token)
    if err != nil {
        http.Error(w, "Failed to verify email", http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Invalid or expired token", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
}
