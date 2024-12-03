package model

// User struct
type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"password,omitempty"`
	Phone          string `json:"phone"`
	MembershipTier string `json:"membership_tier"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// Reservation struct
type Reservation struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	VehicleID   int       `json:"vehicle_id"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
