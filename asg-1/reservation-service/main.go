package main

import (
	"log"
	"net/http"
	"reservation-service/handler"
	"reservation-service/service"
	"github.com/rs/cors" 
	"fmt"
	"github.com/gorilla/mux"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, 
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

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
	handler := corsHandler.Handler(router)

	// Start the server
	fmt.Println("Reservation service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
