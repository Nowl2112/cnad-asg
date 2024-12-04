package model

// Vehicle struct represents a vehicle in the car-sharing system
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

