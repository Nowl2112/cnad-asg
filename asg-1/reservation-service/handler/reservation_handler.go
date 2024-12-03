package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reservation-service/model"
	"reservation-service/service"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

// AddReservation handles the creation of a new reservation.
func AddReservation(w http.ResponseWriter, r *http.Request) {
	var reservation model.Reservation
	// Decode the JSON payload from the request body.
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to add the reservation.
	err := service.AddReservation(&reservation)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding reservation: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the created reservation.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

// GetReservation handles fetching a reservation by ID.
func GetReservation(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL parameters.
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing reservation ID in request", http.StatusBadRequest)
		return
	}

	// Convert the ID from string to int.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
		return
	}

	// Call the service to get the reservation.
	reservation, err := service.GetReservation(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving reservation: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the reservation details.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}