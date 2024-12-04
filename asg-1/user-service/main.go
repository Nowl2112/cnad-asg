package main

import (
    "log"
    "net/http"
    "user-service/handler"
    "user-service/service"
	"github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	dsn := "user:Momo9119!@tcp(localhost:3306)/np_db"
	service.InitDB(dsn)

	// Set up the router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users/register", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", handler.Login).Methods("POST")
	router.HandleFunc("/user/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}/rental-history", handler.GetRentalHistory).Methods("GET")
	router.HandleFunc("/verify", handler.VerifyEmail).Methods("GET")

	// Add CORS headers
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins; replace "*" with specific origin if needed
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allowed headers
	)(router)

	// Start the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}