package model

import "time"

type Reservation struct{
	ID			int       `json:"id"`
	UserID      int       `json:"user_id"`
	VehicleID   int       `json:"vehicle_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type Vehicle struct {
	ID           int     `json:"id"`
	LicensePlate string  `json:"license_plate"`
	Model        string  `json:"model"`
	ChargeLevel  float64 `json:"charge_level"`
	Cleanliness  string  `json:"cleanliness"`
	Location     string  `json:"location"`
	Cost         float64 `json:"cost"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

type EstimateRequest struct {
	VehicleID int       `json:"vehicle_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}

type EstimateResponse struct {
	TotalCost float64 `json:"total_cost"`
}