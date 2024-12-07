package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reservation-service/model"
	"reservation-service/service"
	"strconv"
	"github.com/gorilla/mux"
	"time"
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

// UpdateReservation updates an existing active reservation.
func UpdateReservation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "Missing reservation ID in request", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
        return
    }

    var reservation model.Reservation
    if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    updateReservation := &model.UpdateReservation{
        ID:         id,
        StartTime:  reservation.StartTime,
        EndTime:    reservation.EndTime,
        Status:     reservation.Status,
        TotalPrice: reservation.TotalPrice,
    }

    err = service.UpdateReservation(id, updateReservation)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error updating reservation: %v", err), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully"})
}

// CancelReservation cancels an existing active reservation.
func CancelReservation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "Missing reservation ID in request", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
        return
    }

    err = service.CancelReservation(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error cancelling reservation: %v", err), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Reservation cancelled successfully"})
}

func CalculateReservationCostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request body into the EstimateRequest model
	var request model.EstimateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request fields
	if request.VehicleID == 0 || request.UserID == 0 {
		http.Error(w, "Missing vehicle_id or user_id in request", http.StatusBadRequest)
		return
	}

	// Convert the Unix timestamp (int64) to time.Time using time.Unix
	startTime := time.Unix(request.StartTime, 0)  // Convert Unix timestamp to time.Time
	endTime := time.Unix(request.EndTime, 0)      // Convert Unix timestamp to time.Time

	// Convert back the time.Time to Unix timestamp (int64)
	startUnix := startTime.Unix()  // Get Unix timestamp in int64
	endUnix := endTime.Unix()      // Get Unix timestamp in int64

	// Call the service to calculate the estimated cost with Unix timestamps
	totalCost, err := service.CalculateEstimatedCost(request.VehicleID, request.UserID, startUnix, endUnix)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate cost: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the calculated total cost
	response := model.EstimateResponse{TotalCost: totalCost}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}
