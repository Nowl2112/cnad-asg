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
	CreatedAt   string    `json:"created_at, omitempty"`
	UpdatedAt   string    `json:"updated_at, omitempty"`
	CarPlate   string    `json:"car_plate,omitempty"`

}
type UpdateReservation struct{
	ID			int       `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
	TotalPrice  float64   `json:"total_price"`
}


type EstimateRequest struct {
	VehicleID int    `json:"vehicle_id"`
	UserID    int    `json:"user_id"`
	StartTime int64  `json:"start_time"`  
	EndTime   int64  `json:"end_time"`   
}

type EstimateResponse struct {
	TotalCost float64 `json:"total_cost"`
}
