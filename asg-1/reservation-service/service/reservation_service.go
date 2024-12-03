package service
import (
	"database/sql"
	"fmt"
	"log"
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
//AddReservation
func AddReservation(reservation *model.Reservation) error {
    if reservation.Status == "" {
        reservation.Status = "active"
    }

    query := "INSERT INTO reservations(user_id, vehicle_id, start_time, end_time, total_price, status) VALUES(?, ?, ?, ?, ?, ?);"
    result, err := db.Exec(query, reservation.UserID, reservation.VehicleID, reservation.StartTime, reservation.EndTime, reservation.TotalPrice, reservation.Status)
    if err != nil {
        return fmt.Errorf("Failed to add reservation: %v", err)
    }

    reservationID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("Failed to retrieve last insert ID: %v", err)
    }
	updateVehicleQuery := "UPDATE vehicles SET available = false WHERE id = ?"
    _, err = db.Exec(updateVehicleQuery, reservation.VehicleID)
    if err != nil {
        return fmt.Errorf("Failed to update vehicle availability: %v", err)}

    reservation.ID = int(reservationID)
    return nil
	
}

func GetReservation(id int) (*model.Reservation, error) {
	var reservation model.Reservation
	query := "SELECT id, user_id, vehicle_id, start_time, end_time, total_price, status, created_at, updated_at FROM reservations WHERE id = ?"

	err := db.QueryRow(query, id).Scan(&reservation.ID,&reservation.UserID,&reservation.VehicleID,&reservation.StartTime,&reservation.EndTime,&reservation.TotalPrice,&reservation.Status,&reservation.CreatedAt,&reservation.UpdatedAt,)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve reservation: %v", err)
	}

	return &reservation, nil
}
