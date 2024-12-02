package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// User struct
	type User struct {
		ID             int    `json:"id"`
		Email          string `json:"email"`
		Password       string `json:"password,omitempty"` 
		Phone          string `json:"phone"`
		MembershipTier string `json:"membership_tier"`
		CreatedAt      string `json:"created_at,omitempty"`
		UpdatedAt      string `json:"updated_at,omitempty"`
	}

	var db *sql.DB

	func initDB() {
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
	func registerUser(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Insert user 
		query := "INSERT INTO users (email, password_hash, phone, membership_tier) VALUES (?, ?, ?, ?)"
		result, err := db.Exec(query, user.Email, hashedPassword, user.Phone, user.MembershipTier)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		userID, _ := result.LastInsertId()
		user.ID = int(userID)

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}

	// Get user by ID
	func getUser(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var user User
		query := "SELECT id, email, phone, membership_tier, created_at, updated_at FROM users WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Phone, &user.MembershipTier, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}

	// Login user
	func login(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var hashedPassword string
		var user User
		query := "SELECT id, email, password_hash, phone, membership_tier FROM users WHERE email = ?"
		err := db.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Phone, &user.MembershipTier)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	}
	func updateUser(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
	2
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	
		query := "UPDATE users SET email = ?, password_hash = ?, phone = ?, updated_at = NOW() WHERE id = ?"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
	
		_, err = db.Exec(query, user.Email, hashedPassword, user.Phone, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
	}

	func main() {
		initDB()
		defer db.Close()

		router := mux.NewRouter()

		router.HandleFunc("/users", registerUser).Methods("POST")
		router.HandleFunc("/users/{id}", getUser).Methods("GET")
		router.HandleFunc("/login", login).Methods("POST")
		router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
		fmt.Println("User service is running on port 8080")
		log.Fatal(http.ListenAndServe(":8080", router))
	}
