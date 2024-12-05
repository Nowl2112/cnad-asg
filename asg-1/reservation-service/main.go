package main

import (
	"log"
	"net/http"
	"reservation-service/handler"
	"reservation-service/service"

	"fmt"
	"github.com/gorilla/mux"
)

func main() {
	dsn := "user:Momo9119!@tcp(localhost:3306)/np_db" 
	if err := service.InitDB(dsn); err != nil {
		log.Fatal(err)
	}


	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/reservations", handler.AddReservation).Methods("POST")
	router.HandleFunc("/reservations/{id}", handler.GetReservation).Methods("GET")
	router.HandleFunc("/reservations/{id}/complete", handler.CompleteReservation).Methods("PUT")
	router.HandleFunc("/reservations/{id}", handler.UpdateReservation).Methods("PUT")
	router.HandleFunc("/reservations/{id}/cancel", handler.CancelReservation).Methods("PUT")
	router.HandleFunc("/reservations/estimate", handler.CalculateReservationCostHandler).Methods("POST")

	// Start the server
	fmt.Println("Reservation service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
