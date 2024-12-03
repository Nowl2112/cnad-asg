package model
//reservation struct
type Reservation struct{
	ID			int		`json:"id"`
	UserID      int       `json:"user_id"`
	VehicleID   int       `json:"vehicle_id"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}