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

// CompleteReservation handles completing a reservation.
func CompleteReservation(w http.ResponseWriter, r *http.Request) {
    // Extract the reservation ID from the URL parameters
    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "Missing reservation ID in request", http.StatusBadRequest)
        return
    }

    // Convert the ID from string to int
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
        return
    }

    // Call the service to complete the reservation
    err = service.CompleteReservation(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error completing reservation: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Reservation completed successfully"})
}

// GetAvailableVehicles 
func GetAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the request body data
	type TimeRange struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	var timeRange TimeRange

	// Decode the JSON payload from the request body.
	if err := json.NewDecoder(r.Body).Decode(&timeRange); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if start_time and end_time are provided
	if timeRange.StartTime == "" || timeRange.EndTime == "" {
		http.Error(w, "start_time and end_time are required fields in the request body", http.StatusBadRequest)
		return
	}

	// Call the service to get available vehicles based on the provided times
	vehicles, err := service.GetAvailableVehicles(timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving available vehicles: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the available vehicles
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}
