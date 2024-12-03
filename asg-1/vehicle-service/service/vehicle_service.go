package service

import (
	"database/sql"
	"fmt"
	"log"
	"vehicle-service/model"
)

// DB instance (to be initialized in the main function)
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

// AddVehicle adds a new vehicle to the database
func AddVehicle(vehicle *model.Vehicle) error {
	query := "INSERT INTO vehicles (license_plate, model, charge_level, cleanliness, available, location, cost) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, vehicle.LicensePlate, vehicle.Model, vehicle.ChargeLevel, vehicle.Cleanliness, vehicle.Available, vehicle.Location, vehicle.Cost)
	if err != nil {
		return fmt.Errorf("Failed to add vehicle: %v", err)
	}

	vehicleID, _ := result.LastInsertId()
	vehicle.ID = int(vehicleID)
	return nil
}

// GetAvailable retrieves all available vehicles
func GetAvailable() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	query := "SELECT id, license_plate, model, charge_level, cleanliness, available, location, cost, created_at, updated_at FROM vehicles WHERE available = true"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve available vehicles: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle model.Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Model, &vehicle.ChargeLevel, &vehicle.Cleanliness, &vehicle.Available, &vehicle.Location, &vehicle.Cost, &vehicle.CreatedAt, &vehicle.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Failed to scan vehicle: %v", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating through available vehicles: %v", err)
	}

	if len(vehicles) == 0 {
		return nil, fmt.Errorf("No available vehicles")
	}

	return vehicles, nil
}

// GetVehicle retrieves a vehicle by ID
func GetVehicle(id int) (*model.Vehicle, error) {
	var vehicle model.Vehicle
	query := "SELECT id, license_plate, model, charge_level, cleanliness, available, location, cost, created_at, updated_at FROM vehicles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Model, &vehicle.ChargeLevel, &vehicle.Cleanliness, &vehicle.Available, &vehicle.Location, &vehicle.Cost, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Vehicle not found")
	} else if err != nil {
		return nil, fmt.Errorf("Failed to retrieve vehicle: %v", err)
	}
	return &vehicle, nil
}

// UpdateVehicle updates an existing vehicle
func UpdateVehicle(id int, vehicle *model.Vehicle) error {
	query := "UPDATE vehicles SET license_plate = ?, model = ?, charge_level = ?, cleanliness = ?, available = ?, location = ?, cost = ?, updated_at = NOW() WHERE id = ?"
	_, err := db.Exec(query, vehicle.LicensePlate, vehicle.Model, vehicle.ChargeLevel, vehicle.Cleanliness, vehicle.Available, vehicle.Location, vehicle.Cost, id)
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
