package model

type EstimateRequest struct {
	VehicleID int    `json:"vehicle_id"`
	UserID    int    `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}


type EstimateResponse struct {
	TotalCost float64 `json:"total_cost"`
}

type Item struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}