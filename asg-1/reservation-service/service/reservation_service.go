package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"reservation-service/model"
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

// AddReservation handles adding a new reservation with a null end_time
func AddReservation(reservation *model.Reservation) error {
    if reservation.Status == "" {
        reservation.Status = "active"  // Default status
    }

    // Set start_time to the current time if not provided
    if reservation.StartTime.IsZero() {
        reservation.StartTime = time.Now()
    }

    // Set end_time to NULL for new reservations
    query := "INSERT INTO reservations(user_id, vehicle_id, start_time, end_time, total_price, status) VALUES(?, ?, ?, NULL, ?, ?);"
    result, err := db.Exec(query, reservation.UserID, reservation.VehicleID, reservation.StartTime, reservation.TotalPrice, reservation.Status)
    if err != nil {
        return fmt.Errorf("Failed to add reservation: %v", err)
    }

    reservationID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("Failed to retrieve last insert ID: %v", err)
    }

    // Update vehicle availability to false
    updateVehicleQuery := "UPDATE vehicles SET available = false WHERE id = ?"
    _, err = db.Exec(updateVehicleQuery, reservation.VehicleID)
    if err != nil {
        return fmt.Errorf("Failed to update vehicle availability: %v", err)
    }

    reservation.ID = int(reservationID)
    return nil
}

// CompleteReservation updates the reservation with an end_time and makes the vehicle available again
func CompleteReservation(reservationID int) error {
    // Get the current time as the end_time
    endTime := time.Now()

    // Update reservation with end_time and status to 'completed'
    updateReservationQuery := "UPDATE reservations SET end_time = ?, status = 'completed' WHERE id = ?"
    _, err := db.Exec(updateReservationQuery, endTime, reservationID)
    if err != nil {
        return fmt.Errorf("Failed to complete reservation: %v", err)
    }

    // Get the vehicle ID from the reservation
    var vehicleID int
    err = db.QueryRow("SELECT vehicle_id FROM reservations WHERE id = ?", reservationID).Scan(&vehicleID)
    if err != nil {
        return fmt.Errorf("Failed to retrieve vehicle ID for reservation: %v", err)
    }

    // Update vehicle availability to true
    updateVehicleQuery := "UPDATE vehicles SET available = true WHERE id = ?"
    _, err = db.Exec(updateVehicleQuery, vehicleID)
    if err != nil {
        return fmt.Errorf("Failed to update vehicle availability: %v", err)
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
