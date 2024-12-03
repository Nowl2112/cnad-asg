package main

import (
	"fmt"
	"log"
	"net/http"
	"vehicle-service/service"
	"vehicle-service/handler"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	dsn := "user:Momo9119!@tcp(localhost:3306)/np_db" // Replace with your DB credentials
	if err := service.InitDB(dsn); err != nil {
		log.Fatal(err)
	}

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/vehicles", handler.AddVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{id}", handler.GetVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{id}", handler.UpdateVehicle).Methods("PUT")
	router.HandleFunc("/vehicles/{id}", handler.DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/available",handler.GetAvailable).Methods("GET")

	// Start server
	fmt.Println("Vehicle service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
