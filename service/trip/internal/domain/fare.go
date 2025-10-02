package domain

import "math"

// rough fare: base + (distance_km * per_km)
func FareEstimateCents(lat1, lng1, lat2, lng2 float64) int {
	distKm := haversine(lat1, lng1, lat2, lng2)
	base := 15000.0 // 15k VND base (example)
	perKm := 9000.0 // 9k VND per km
	total := base + distKm*perKm
	return int(total)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
