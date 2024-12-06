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

    var startTime []byte
    var endTime []byte
    err := db.QueryRow(query, id).Scan(&reservation.ID, &reservation.UserID, &reservation.VehicleID, &startTime, &endTime, &reservation.TotalPrice, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve reservation: %v", err)
    }

    // Convert byte slices to time.Time
    reservation.StartTime, err = time.Parse("2006-01-02 15:04:05", string(startTime))
    if err != nil {
        return nil, fmt.Errorf("Failed to parse start_time: %v", err)
    }

    reservation.EndTime, err = time.Parse("2006-01-02 15:04:05", string(endTime))
    if err != nil {
        return nil, fmt.Errorf("Failed to parse end_time: %v", err)
    }

    return &reservation, nil
}

// UpdateReservation updates an existing reservation if conditions are met.
func UpdateReservation(id int, updated *model.Reservation) error {
    var current model.Reservation
    query := "SELECT id, start_time, status FROM reservations WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&current.ID, &current.StartTime, &current.Status)
    if err != nil {
        return fmt.Errorf("Reservation not found or retrieval error: %v", err)
    }

    if current.Status != "active" {
        return fmt.Errorf("Only active reservations can be updated")
    }

    if time.Until(current.StartTime) < 5*time.Minute {
        return fmt.Errorf("Reservations can only be updated at least 5 minutes before the start time")
    }

    if updated.StartTime.Before(time.Now().Add(5 * time.Minute)) || updated.EndTime.Before(updated.StartTime) {
        return fmt.Errorf("Invalid new start or end time")
    }

    query = "UPDATE reservations SET start_time = ?, end_time = ?, updated_at = NOW() WHERE id = ?"
    _, err = db.Exec(query, updated.StartTime, updated.EndTime, id)
    if err != nil {
        return fmt.Errorf("Failed to update reservation: %v", err)
    }

    return nil
}

// CancelReservation cancels an active reservation if conditions are met.
func CancelReservation(id int) error {
    var current model.Reservation
    query := "SELECT id, start_time, status FROM reservations WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&current.ID, &current.StartTime, &current.Status)
    if err != nil {
        return fmt.Errorf("Reservation not found or retrieval error: %v", err)
    }

    if current.Status != "active" {
        return fmt.Errorf("Only active reservations can be cancelled")
    }

    if time.Until(current.StartTime) < 5*time.Minute {
        return fmt.Errorf("Reservations can only be cancelled at least 5 minutes before the start time")
    }

    query = "UPDATE reservations SET status = 'cancelled', updated_at = NOW() WHERE id = ?"
    _, err = db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("Failed to cancel reservation: %v", err)
    }

    return nil
}

func CalculateEstimatedCost(vehicleID, userID int, startTime, endTime int64) (float64, error) {
	// Convert Unix timestamps to time.Time
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

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
