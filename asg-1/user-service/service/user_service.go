package service

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"user-service/model"
)

var db *sql.DB

// Initialize the DB connection
func InitDB() {
	var err error
	dsn := "user:Momo9119!@tcp(localhost:3306)/np_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Database connected!")
}

// Register new user
func RegisterUser(user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, fmt.Errorf("failed to hash password: %v", err)
	}

	query := "INSERT INTO users (email, password_hash, phone, membership_tier) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, user.Email, hashedPassword, user.Phone, user.MembershipTier)
	if err != nil {
		return user, fmt.Errorf("failed to register user: %v", err)
	}

	userID, _ := result.LastInsertId()
	user.ID = int(userID)
	user.Password = "" // Do not return the password

	return user, nil
}

// Get user by ID
func GetUser(id int) (model.User, error) {
	var user model.User
	query := "SELECT id, email, phone, membership_tier, created_at, updated_at FROM users WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Phone, &user.MembershipTier, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, fmt.Errorf("failed to retrieve user: %v", err)
	}
	return user, nil
}

// Login user
func Login(email, password string) (bool, error) {
	var hashedPassword string
	var user model.User
	query := "SELECT id, email, password_hash, phone, membership_tier FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Phone, &user.MembershipTier)
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("invalid email or password")
	} else if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %v", err)
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("invalid email or password")
	}
	return true, nil
}

// Update user
func UpdateUser(id int, user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := "UPDATE users SET email = ?, password_hash = ?, phone = ?, updated_at = NOW() WHERE id = ?"
	_, err = db.Exec(query, user.Email, hashedPassword, user.Phone, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// Get rental history
func GetRentalHistory(userID int) ([]model.Reservation, error) {
	var reservations []model.Reservation
	query := `SELECT r.id, r.user_id, r.vehicle_id, r.start_time, r.end_time, r.total_price, r.status, r.created_at, r.updated_at
	          FROM reservations r
	          WHERE r.user_id = ? AND r.end_time < NOW()` // Only past reservations

	rows, err := db.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no past reservations found for this user")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve rental history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reservation model.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.VehicleID, &reservation.StartTime, &reservation.EndTime, &reservation.TotalPrice, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %v", err)
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}
