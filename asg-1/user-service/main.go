package main

import (
	"log"
	"net/http"
	"user-service/handler"
	"user-service/service"
	"github.com/gorilla/mux"
)

func main() {
	service.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/users/register", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users/login", handler.Login).Methods("POST")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}/history", handler.GetRentalHistory).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
