package service

import (
	"database/sql"
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"user-service/model"
	"crypto/rand"
	"encoding/hex"
	"time"
	"net/smtp"
)

var db *sql.DB


// InitDB initializes the database connection
func InitDB(dsn string) {
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
    fmt.Println("Database connected!")
}

// GetDB returns the database connection
func GetDB() *sql.DB {
    return db
}
// Register new user
func RegisterUser(user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, fmt.Errorf("failed to hash password: %v", err)
	}

	token, err := GenerateToken()
	if err != nil {
		return user, fmt.Errorf("failed to generate verification token: %v", err)
	}

	expiry := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	query := `INSERT INTO users (email, password_hash, phone, membership_tier, verification_token, token_expiry)
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Email, hashedPassword, user.Phone, user.MembershipTier, token, expiry)
	if err != nil {
		return user, fmt.Errorf("failed to register user: %v", err)
	}

	userID, _ := result.LastInsertId()
	user.ID = int(userID)
	user.VerificationToken = token
	user.TokenExpiry = expiry.Format("2006-01-02 15:04:05")

	// Send verification email
	go sendVerificationEmail(user.Email, token)

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
//login
func Login(email, password string) (bool, model.User, error) {
	var hashedPassword string
	var user model.User
	query := "SELECT id, email, password_hash, phone, membership_tier FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Phone, &user.MembershipTier)
	if err == sql.ErrNoRows {
		return false, user, fmt.Errorf("invalid email or password")
	} else if err != nil {
		return false, user, fmt.Errorf("failed to retrieve user: %v", err)
	}
	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, user, fmt.Errorf("invalid email or password")
	}

	return true, user, nil  // Return the user details
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


// GenerateToken creates a random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Register new user with email verification

func sendVerificationEmail(email, token string) {
	// Configure your SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	username := "kotaro.da.kat@gmail.com" // Your Gmail address
	password := "mkin ajob zriq oifi"      // Your App Password (NOT your Gmail account password)

	// Create the email message
	from := username // Use the authenticated username as the sender
	subject := "Email Verification"
	body := fmt.Sprintf("Please verify your email by clicking the link: http://localhost:8080/verify?token=%s", token)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, email, subject, body)

	// Send the email
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg)); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
	} else {
		fmt.Println("Verification email sent successfully!")
	}
}
// Get user by email
func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	query := "SELECT id, email, phone, membership_tier, created_at, updated_at FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Phone, &user.MembershipTier, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, fmt.Errorf("failed to retrieve user: %v", err)
	}
	return user, nil
}
