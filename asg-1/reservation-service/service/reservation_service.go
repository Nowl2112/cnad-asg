package service

import (
	"database/sql"
	"fmt"
	"log"
	"reservation-service/model"
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

// AddReservation
func AddReservation(reservation *model.Reservation) error {
    if reservation.Status == "" {
        reservation.Status = "active" // Default status
    }

    if reservation.StartTime.IsZero() || reservation.EndTime.IsZero() {
        return fmt.Errorf("start_time and end_time must be provided")
    }

    // Ensure end_time is after start_time
    if reservation.EndTime.Before(reservation.StartTime) {
        return fmt.Errorf("end_time must be after start_time")
    }

    // Check for overlapping reservations
    queryConflict := `
        SELECT COUNT(*) 
        FROM reservations 
        WHERE vehicle_id = ? 
          AND status = 'active'
          AND NOT (
              end_time <= ? OR start_time >= ?
          );
    `
    var conflictCount int
    err := db.QueryRow(queryConflict, reservation.VehicleID, reservation.StartTime, reservation.EndTime).Scan(&conflictCount)
    if err != nil {
        return fmt.Errorf("Failed to check vehicle availability: %v", err)
    }

    if conflictCount > 0 {
        return fmt.Errorf("The vehicle is not available for the requested time slot")
    }

    // Calculate the total price (hours difference * vehicle hourly cost)
    var vehicleCost float64
    err = db.QueryRow("SELECT cost FROM vehicles WHERE id = ?", reservation.VehicleID).Scan(&vehicleCost)
    if err != nil {
        return fmt.Errorf("Failed to fetch vehicle cost: %v", err)
    }

    duration := reservation.EndTime.Sub(reservation.StartTime).Hours()
    reservation.TotalPrice = duration * vehicleCost

    // Insert the reservation into the database
    query := "INSERT INTO reservations(user_id, vehicle_id, start_time, end_time, total_price, status) VALUES(?, ?, ?, ?, ?, ?);"
    result, err := db.Exec(query, reservation.UserID, reservation.VehicleID, reservation.StartTime, reservation.EndTime, reservation.TotalPrice, reservation.Status)
    if err != nil {
        return fmt.Errorf("Failed to add reservation: %v", err)
    }

    reservationID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("Failed to retrieve last insert ID: %v", err)
    }

    reservation.ID = int(reservationID)
    return nil
}


// CompleteReservation updates the reservation with an end_time and makes the vehicle available again
func CompleteReservation(reservationID int) error {
    // Update reservation with status to 'completed'
    updateReservationQuery := "UPDATE reservations SET status = 'completed' WHERE id = ?"
    _, err := db.Exec(updateReservationQuery, reservationID)
    if err != nil {
        return fmt.Errorf("Failed to complete reservation: %v", err)
    }

    return nil
}

// GetReservation retrieves a reservation by ID
func GetReservation(id int) (*model.Reservation, error) {
	var reservation model.Reservation
	query := "SELECT id, user_id, vehicle_id, start_time, end_time, total_price, status, created_at, updated_at FROM reservations WHERE id = ?"

	err := db.QueryRow(query, id).Scan(&reservation.ID, &reservation.UserID, &reservation.VehicleID, &reservation.StartTime, &reservation.EndTime, &reservation.TotalPrice, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve reservation: %v", err)
	}

	return &reservation, nil
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

func CalculateEstimatedCost(vehicleID int, startTime, endTime time.Time) (float64, error) {
	// Fetch the vehicle's cost per unit time (e.g., hourly rate)
	var costPerUnit float64
	query := "SELECT cost FROM vehicles WHERE id = ?"
	err := db.QueryRow(query, vehicleID).Scan(&costPerUnit)
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
	totalCost := costPerUnit * duration
	return totalCost, nil
}