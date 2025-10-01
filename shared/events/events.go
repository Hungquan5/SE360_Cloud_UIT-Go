// /shared/events/events.go
package events

type TripRequested struct {
	TripID               string
	OriginLat, OriginLng float64
	DestLat, DestLng     float64
}

type DriverAccepted struct {
	TripID   string
	DriverID string
}

type DriverLocation struct {
	DriverID string
	Lat, Lng float64
	TripID   string // optional when enroute
}

type TripCompleted struct {
	TripID string
}
