package domain

import "time"

type DriverStatus string

const (
	StatusOnline  DriverStatus = "online"
	StatusOffline DriverStatus = "offline"
)

type UpdateLocationRequest struct {
	Lat    float64      `json:"lat"`
	Lon    float64      `json:"lon"`
	At     time.Time    `json:"at"`
	Status DriverStatus `json:"status"` // optional; default online
}

type Nearby struct {
	DriverID string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Meters   float64 `json:"meters"`
}
