package service

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"log"
	"gopkg.in/gomail.v2"

)

const (
	SMTPHost     = "smtp.gmail.com" // Use your email provider's SMTP server
	SMTPPort     = 587             // SMTP port
	SMTPUsername = "kotaro.da.kat@gmail.com"
	SMTPPassword = "mkin ajob zriq oifi" // Use an app-specific password if using Gmail
)
type Item struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}
// DB instance 
var db *sql.DB

// Initialize the database connection
func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Error pinging the database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}


func init() {
    stripeKey := os.Getenv("STRIPE_API_KEY")
    if stripeKey == "" {
        log.Fatalf("Stripe secret key is not set in the environment variables")
    }
    log.Printf("Stripe Key: %s", stripeKey[:8]) // Only show the first few characters for debugging
    stripe.Key = stripeKey
}


  type Reservation struct {
	ReservationID int
	VehicleID     int
	UserID        int
	StartTime     string
	EndTime       string
	TotalCost     float64
}

func FetchReservationDetails(reservationID int) (*Reservation, error) {
	var reservation Reservation
	query := `
        SELECT r.id, r.vehicle_id, r.user_id, r.start_time, r.end_time, r.estimated_cost
        FROM reservations r
        WHERE r.id = ?`
	err := db.QueryRow(query, reservationID).Scan(
		&reservation.ReservationID,
		&reservation.VehicleID,
		&reservation.UserID,
		&reservation.StartTime,
		&reservation.EndTime,
		&reservation.TotalCost,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Reservation not found")
	} else if err != nil {
		return nil, fmt.Errorf("Failed to fetch reservation: %v", err)
	}
	return &reservation, nil
}

func CreatePaymentIntentForReservation(reservationID int) (*stripe.PaymentIntent, error) {
    // Fetch reservation details
    reservation, err := FetchReservationDetails(reservationID)
    if err != nil {
        log.Printf("Failed to fetch reservation: %v", err)
        return nil, err
    }

    // Log reservation details and total cost for debugging
    log.Printf("Fetched reservation: %+v, Total cost: %.2f", reservation, reservation.TotalCost)

    // Check if the reservation total cost is valid (greater than 0)
    if reservation.TotalCost <= 0 {
        log.Printf("Invalid reservation total cost: %.2f. Amount must be greater than zero.", reservation.TotalCost)
        return nil, fmt.Errorf("Invalid reservation total cost: %.2f. Amount must be greater than zero.", reservation.TotalCost)
    }

    // Convert the total cost to cents for Stripe (as Stripe uses smallest currency units)
    amount := int64(reservation.TotalCost * 100)
    log.Printf("Amount in cents: %d", amount)

    // Check if the amount is below the minimum charge allowed (50 cents for USD)
    if amount < 50 {
        log.Printf("Amount is below the minimum (50 cents), adjusting to 50 cents.")
        amount = 50 // Set to 50 cents if less
    }

    // Create PaymentIntent parameters
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(amount),
        Currency: stripe.String(string(stripe.CurrencyUSD)),
        AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
            Enabled: stripe.Bool(true),
        },
        Metadata: map[string]string{
            "reservation_id": fmt.Sprintf("%d", reservationID),
        },
    }

    // Create the PaymentIntent
    pi, err := paymentintent.New(params)
    if err != nil {
        log.Printf("Failed to create payment intent: %v", err)
        return nil, err
    }

    return pi, nil
}


func CalculateOrderAmount(items []Item) int64 {
	var total int64
	for _, item := range items {
		total += item.Amount
	}
	return total
}

func CreatePaymentIntent(items []Item) (*stripe.PaymentIntent, error) {
	amount := CalculateOrderAmount(items)

	// Create PaymentIntent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Create the PaymentIntent
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Failed to create payment intent: %v", err)
		return nil, err
	}

	return pi, nil
}

func SendReservationEmail(toEmail string, reservationID int, vehicleplate string, totalCost float64) error {
	// Create the email content
	subject := "Your Reservation Details"
	body := fmt.Sprintf(`
		Hello,

		Thank you for your reservation. Here are your reservation details:

		Reservation ID: %d
		Vehicle plate: %d
		Total Cost: $%.2f

		We hope you enjoy your experience!

		Best regards,
		Your Company Name
	`, reservationID, vehicleplate, totalCost)

	// Create the email message
	message := gomail.NewMessage()
	message.SetHeader("From", SMTPUsername)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	// Configure the SMTP client
	dialer := gomail.NewDialer(SMTPHost, SMTPPort, SMTPUsername, SMTPPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}