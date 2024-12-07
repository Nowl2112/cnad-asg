package handler

import (
	"fmt"
	"net/http"
	"vehicle-service/service"
	"vehicle-service/model"
	"github.com/gorilla/mux"
	"strconv"
	"encoding/json" 
	_ "github.com/go-sql-driver/mysql"

)

// AddVehicle handler
func AddVehicle(w http.ResponseWriter, r *http.Request) {
	var vehicle model.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := service.AddVehicle(&vehicle)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

// GetVehicle handler
func GetVehicle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	vehicle, err := service.GetVehicle(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

// UpdateVehicle handler
func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	var vehicle model.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = service.UpdateVehicle(id, &vehicle)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle updated successfully"})
}

// DeleteVehicle handler
func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	err = service.DeleteVehicle(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle deleted successfully"})
}

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
