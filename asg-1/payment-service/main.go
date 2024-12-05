package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-service/handler"
	"payment-service/service"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins (can be restricted to specific origins)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	dsn := "user:Momo9119!@tcp(localhost:3306)/np_db"
	if err := service.InitDB(dsn); err != nil {
		log.Fatal(err)
	}

	stripeKey := os.Getenv("STRIPE_API_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_API_KEY not set in environment")
	}

	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/create-payment-intent", handler.HandleCreatePaymentIntent).Methods("POST")
	router.HandleFunc("/create-payment-for-reservation", handler.HandleCreatePaymentForReservation).Methods("POST")
	
	
	handler := corsHandler.Handler(router)

	// Start the server
	fmt.Println("Payment service is running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", handler))
}
