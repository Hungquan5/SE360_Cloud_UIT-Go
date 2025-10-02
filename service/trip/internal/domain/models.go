package domain

import "time"

type Trip struct {
	ID                string    `json:"id"`
	PassengerID       string    `json:"passenger_id"`
	DriverID          *string   `json:"driver_id,omitempty"`
	OriginLat         float64   `json:"origin_lat"`
	OriginLng         float64   `json:"origin_lng"`
	DestLat           float64   `json:"dest_lat"`
	DestLng           float64   `json:"dest_lng"`
	FareEstimateCents int       `json:"fare_estimate_cents"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type DriverLocation struct {
	DriverID string  `json:"driverId"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	TripID   string  `json:"tripId,omitempty"`
}

type DriverAccepted struct {
	TripID   string `json:"tripId"`
	DriverID string `json:"driverId"`
}

type TripRequested struct {
	TripID    string  `json:"tripId"`
	OriginLat float64 `json:"originLat"`
	OriginLng float64 `json:"originLng"`
	DestLat   float64 `json:"destLat"`
	DestLng   float64 `json:"destLng"`
}
