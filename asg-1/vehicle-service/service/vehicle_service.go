package service

import (
	"database/sql"
	"fmt"
	"log"
	"vehicle-service/model"
	"time"
)

// DB instance 
var db *sql.DB

// Initialize the database connection
func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Error pinging the database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

// AddVehicle 
func AddVehicle(vehicle *model.Vehicle) error {
	query := "INSERT INTO vehicles (license_plate, model, charge_level, cleanliness,location, cost) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, vehicle.LicensePlate, vehicle.Model, vehicle.ChargeLevel, vehicle.Cleanliness,  vehicle.Location, vehicle.Cost)
	if err != nil {
		return fmt.Errorf("Failed to add vehicle: %v", err)
	}

	vehicleID, _ := result.LastInsertId()
	vehicle.ID = int(vehicleID)
	return nil
}


// GetVehicle retrieves a vehicle by ID
func GetVehicle(id int) (*model.Vehicle, error) {
	var vehicle model.Vehicle
	query := "SELECT id, license_plate, model, charge_level, cleanliness,location, cost, created_at, updated_at FROM vehicles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Model, &vehicle.ChargeLevel, &vehicle.Cleanliness, &vehicle.Location, &vehicle.Cost, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Vehicle not found")
	} else if err != nil {
		return nil, fmt.Errorf("Failed to retrieve vehicle: %v", err)
	}
	return &vehicle, nil
}

// UpdateVehicle updates an existing vehicle
func UpdateVehicle(id int, vehicle *model.Vehicle) error {
	query := "UPDATE vehicles SET license_plate = ?, model = ?, charge_level = ?, cleanliness = ?, location = ?, cost = ?, updated_at = NOW() WHERE id = ?"
	_, err := db.Exec(query, vehicle.LicensePlate, vehicle.Model, vehicle.ChargeLevel, vehicle.Cleanliness, vehicle.Location, vehicle.Cost, id)
	if err != nil {
		return fmt.Errorf("Failed to update vehicle: %v", err)
	}
	return nil
}

// DeleteVehicle deletes a vehicle by ID
func DeleteVehicle(id int) error {
	query := "DELETE FROM vehicles WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete vehicle: %v", err)
	}
	return nil
}

func CalculateEstimatedCost(vehicleID int, startTime, endTime time.Time) (float64, error) {
	// Fetch the vehicle's cost per unit time (e.g., hourly rate)
	var costPerCar float64
	query := "SELECT cost FROM vehicles WHERE id = ?"
	err := db.QueryRow(query, vehicleID).Scan(&costPerCar)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("Vehicle not found")
	} else if err != nil {
		return 0, fmt.Errorf("Failed to retrieve vehicle cost: %v", err)
	}

	// Calculate the duration in hours
	duration := endTime.Sub(startTime).Hours()
	if duration < 0 {
		return 0, fmt.Errorf("End time cannot be before start time")
	}

	// Calculate total cost
	totalCost := costPerCar * duration
	return totalCost, nil
}

func GetAvailableVehicles(startTime, endTime string) ([]model.Vehicle, error) {
	if startTime == "" || endTime == "" {
		return nil, fmt.Errorf("start_time and end_time must be provided")
	}

	query := `
		SELECT id, license_plate, model, charge_level, cleanliness, location, cost
		FROM vehicles 
		WHERE id NOT IN (
			SELECT vehicle_id 
			FROM reservations 
			WHERE status = 'active' 
			  AND NOT (
				  end_time <= ? OR start_time >= ?
			  )
		);
	`

	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve available vehicles: %v", err)
	}
	defer rows.Close()

	var vehicles []model.Vehicle
	for rows.Next() {
		var vehicle model.Vehicle
		// Adjust the Scan to match the fields in your Vehicle model
		if err := rows.Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Model, &vehicle.ChargeLevel, &vehicle.Cleanliness, &vehicle.Location, &vehicle.Cost); err != nil {
			return nil, fmt.Errorf("Error scanning vehicle row: %v", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating through vehicle rows: %v", err)
	}

	return vehicles, nil
}
