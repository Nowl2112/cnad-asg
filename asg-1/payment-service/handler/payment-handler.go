package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-service/service"
	_ "github.com/go-sql-driver/mysql"
)

func SendReservationEmail(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body from the request
	var requestData struct {
		ReservationID int     `json:"reservation_id"`
		UserEmail     string  `json:"user_email"`
		CarPlate      string  `json:"CarPlate"`
		TotalCost     float64 `json:"total_cost"`
	}

	// Decode the request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Call SendReservationEmail with the extracted data
	err := service.SendReservationEmail(requestData.UserEmail, requestData.ReservationID, requestData.CarPlate, requestData.TotalCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email sent successfully"))
}

func HandleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req struct {
		Items []service.Item `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create PaymentIntent using the service
	pi, err := service.CreatePaymentIntent(req.Items)
	if err != nil {
		http.Error(w, "Failed to create payment intent", http.StatusInternalServerError)
		return
	}

	// Respond with the PaymentIntent client secret
	response := struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	}

	writeJSON(w, response)
}

// writeJSON writes the response as JSON
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func HandleCreatePaymentForReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Parse reservation ID from query or JSON body
	var req struct {
		ReservationID int `json:"reservation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ReservationID == 0 {
		http.Error(w, "Missing reservation_id in request", http.StatusBadRequest)
		return
	}

	// Create PaymentIntent using the reservation ID
	pi, err := service.CreatePaymentIntentForReservation(req.ReservationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create payment intent: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the PaymentIntent client secret
	response := struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	}

	writeJSON(w, response)
}

