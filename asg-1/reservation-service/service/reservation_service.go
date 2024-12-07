package service

import (
	"database/sql"
	"fmt"
	"log"
	"reservation-service/model"
	"time"
    
)

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

    reservation.StartTime = reservation.StartTime.In(time.Local)
    reservation.EndTime = reservation.EndTime.In(time.Local)
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

    // Calculate the total price 
    var vehicleCost float64
    err = db.QueryRow("SELECT cost FROM vehicles WHERE id = ?", reservation.VehicleID).Scan(&vehicleCost)
    if err != nil {
        return fmt.Errorf("Failed to fetch vehicle cost: %v", err)
    }

    duration := reservation.EndTime.Sub(reservation.StartTime).Hours()

    // Fetch the user membership tier
    var membershipTier string
    query := "SELECT membership_tier FROM users WHERE id = ?"
    err = db.QueryRow(query, reservation.UserID).Scan(&membershipTier)
    if err == sql.ErrNoRows {
        return fmt.Errorf("User not found")
    } else if err != nil {
        return fmt.Errorf("Failed to retrieve user membership tier: %v", err)
    }

    var totalCost float64
    switch membershipTier {
    case "Basic":
        totalCost = duration * vehicleCost
    case "Premium":
        totalCost = duration * vehicleCost * 0.90 
    case "VIP":
        totalCost = duration * vehicleCost * 0.75 
    default:
        return fmt.Errorf("Unknown membership tier: %s", membershipTier)
    }

    reservation.TotalPrice = totalCost

    // Insert the reservation into the database
    query = "INSERT INTO reservations(user_id, vehicle_id, start_time, end_time, total_price, status) VALUES(?, ?, ?, ?, ?, ?);"
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
    query := `
        SELECT 
            r.id, r.user_id, r.vehicle_id, r.start_time, r.end_time, r.total_price, r.status,v.license_plate
        FROM reservations r
        JOIN vehicles v ON r.vehicle_id = v.id
        WHERE r.id = ?`
    
    var startTime []byte
    var endTime []byte
    err := db.QueryRow(query, id).Scan(&reservation.ID, &reservation.UserID, &reservation.VehicleID, &startTime, &endTime, &reservation.TotalPrice, &reservation.Status,&reservation.CarPlate)
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve reservation: %v", err)
    }

    // Convert byte slices to time.Time (local time)
    reservation.StartTime, err = time.ParseInLocation("2006-01-02 15:04:05", string(startTime), time.Local)
    if err != nil {
        return nil, fmt.Errorf("Failed to parse start_time: %v", err)
    }

    reservation.EndTime, err = time.ParseInLocation("2006-01-02 15:04:05", string(endTime), time.Local)
    if err != nil {
        return nil, fmt.Errorf("Failed to parse end_time: %v", err)
    }

    return &reservation, nil
}

func UpdateReservation(id int, updated *model.UpdateReservation) error {
    // Convert updated start and end times to local time
    updated.StartTime = updated.StartTime.In(time.Local)
    updated.EndTime = updated.EndTime.In(time.Local)

    // Validate that start_time and end_time are not zero or empty
    if updated.StartTime.IsZero() || updated.EndTime.IsZero() {
        return fmt.Errorf("Start time and end time cannot be zero or empty")
    }

    // Fetch the current reservation from the database
    var current model.Reservation
    var startTime []byte
    var endTime []byte
    query := "SELECT id, start_time, end_time, status, user_id, vehicle_id FROM reservations WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&current.ID, &startTime, &endTime, &current.Status, &current.UserID, &current.VehicleID)

    if err != nil {
        return fmt.Errorf("Reservation not found or retrieval error: %v", err)
    }

    if current.Status != "Active" {
        return fmt.Errorf("Only active reservations can be updated")
    }

    // Check for conflicts with existing reservations
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
    err = db.QueryRow(queryConflict, current.VehicleID, updated.StartTime, updated.EndTime).Scan(&conflictCount)
    if err != nil {
        return fmt.Errorf("Failed to check vehicle availability: %v", err)
    }

    if conflictCount > 0 {
        return fmt.Errorf("The vehicle is not available for the requested time slot")
    }

    // Calculate the updated total price based on the new times and user membership tier
    var vehicleCost float64
    err = db.QueryRow("SELECT cost FROM vehicles WHERE id = ?", current.VehicleID).Scan(&vehicleCost)
    if err != nil {
        return fmt.Errorf("Failed to fetch vehicle cost: %v", err)
    }

    duration := updated.EndTime.Sub(updated.StartTime).Hours()

    // Fetch the user's membership tier
    var membershipTier string
    query = "SELECT membership_tier FROM users WHERE id = ?"
    err = db.QueryRow(query, current.UserID).Scan(&membershipTier)
    if err == sql.ErrNoRows {
        return fmt.Errorf("User not found")
    } else if err != nil {
        return fmt.Errorf("Failed to retrieve user membership tier: %v", err)
    }

    var totalCost float64
    switch membershipTier {
    case "Basic":
        totalCost = duration * vehicleCost
    case "Premium":
        totalCost = duration * vehicleCost * 0.90
    case "VIP":
        totalCost = duration * vehicleCost * 0.75
    default:
        return fmt.Errorf("Unknown membership tier: %s", membershipTier)
    }

    // Update the reservation in the database
    query = `
        UPDATE reservations 
        SET start_time = ?, end_time = ?, total_price = ?, updated_at = NOW(), status = 'active' 
        WHERE id = ?`
    _, err = db.Exec(query, updated.StartTime.Format("2006-01-02 15:04:05"), updated.EndTime.Format("2006-01-02 15:04:05"), totalCost, id)
    if err != nil {
        return fmt.Errorf("Failed to update reservation: %v", err)
    }

    return nil
}

// CancelReservation cancels an active reservation if conditions are met.
func CancelReservation(id int) error {
	var current model.Reservation
    var startTime []byte
	query := "SELECT id, start_time, status FROM reservations WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&current.ID, &startTime, &current.Status)
	if err != nil {
		return fmt.Errorf("Reservation not found or retrieval error: %v", err)
	}

	if current.Status != "Active" {
		return fmt.Errorf("Only active reservations can be cancelled")
	}



	query = "UPDATE reservations SET status = 'Cancelled', updated_at = NOW() WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to cancel reservation: %v", err)
	}

	return nil
}

func CalculateEstimatedCost(vehicleID, userID int, startTime, endTime int64) (float64, error) {
	// Convert Unix timestamps to time.Time (local time)
	start := time.Unix(startTime, 0).In(time.Local)
	end := time.Unix(endTime, 0).In(time.Local)

	// Fetch the vehicle's cost per unit time (e.g., hourly rate)
	var costPerUnit float64
	query := "SELECT cost FROM vehicles WHERE id = ?"
	err := db.QueryRow(query, vehicleID).Scan(&costPerUnit)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("Vehicle not found")
	} else if err != nil {
		return 0, fmt.Errorf("Failed to retrieve vehicle cost: %v", err)
	}

	// Fetch the user's membership tier
	var membershipTier string
	query = "SELECT membership_tier FROM users WHERE id = ?"
	err = db.QueryRow(query, userID).Scan(&membershipTier)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("User not found")
	} else if err != nil {
		return 0, fmt.Errorf("Failed to retrieve user membership tier: %v", err)
	}

	// Calculate the duration in hours
	duration := end.Sub(start).Hours()
	if duration < 0 {
		return 0, fmt.Errorf("End time cannot be before start time")
	}

	// Calculate total cost before discount
	totalCost := costPerUnit * duration

	// Apply discount based on membership tier
	switch membershipTier {
	case "Basic":
		// No discount
	case "Premium":
		totalCost *= 0.90 // 10% discount
	case "VIP":
		totalCost *= 0.75 // 25% discount
	default:
		return 0, fmt.Errorf("Unknown membership tier: %s", membershipTier)
	}

	return totalCost, nil
}
