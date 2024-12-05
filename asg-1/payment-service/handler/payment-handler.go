package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-service/service"
	_ "github.com/go-sql-driver/mysql"
)




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


func HandlePaymentSuccess(w http.ResponseWriter, r *http.Request) {
	// Simulate payment success
	reservationID := 123      
	vehicleID := 456          
	totalCost := 100.00      
	startTime := "2024-12-05T10:00:00Z"
	endTime := "2024-12-05T15:00:00Z"  
	userEmail := "user@example.com"    

	// Send the reservation details email
	err := service.SendReservationEmail(userEmail, reservationID, vehicleID, totalCost, startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment processed and email sent successfully"))
}