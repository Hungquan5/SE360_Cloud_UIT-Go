package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DriverProfile struct {
	UserID       string    `json:"user_id"`
	VehiclePlate string    `json:"vehicle_plate"`
	VehicleModel string    `json:"vehicle_model"`
	Approved     bool      `json:"approved"`
	Online       bool      `json:"online"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
