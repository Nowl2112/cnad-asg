package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-service/handler"
	"payment-service/service"

	"github.com/gorilla/mux"
)

func main() {
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
	router.HandleFunc("/payment/estimate", handler.CalculateReservationCostHandler).Methods("POST")
	router.HandleFunc("/create-payment-intent", handler.HandleCreatePaymentIntent).Methods("POST")
	router.HandleFunc("/create-payment-for-reservation", handler.HandleCreatePaymentForReservation).Methods("POST")

	// Start the server
	fmt.Println("Payment service is running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}
