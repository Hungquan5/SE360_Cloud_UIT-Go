package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/driver-svc/internal/domain"
)

const (
	geoKeyDrivers    = "geo:drivers"
	hashKeyDriverFmt = "driver:%s" // fields: status, updated_at
)

type Redis struct{ R *redis.Client }

func NewRedis(addr string) (*Redis, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		opt = &redis.Options{Addr: addr}
	}
	r := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := r.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{R: r}, nil
}

func (s *Redis) SetDriverStatus(ctx context.Context, driverID string, st domain.DriverStatus) error {
	key := fmt.Sprintf(hashKeyDriverFmt, driverID)
	return s.R.HSet(ctx, key, map[string]any{
		"status":     st,
		"updated_at": time.Now().UTC().Format(time.RFC3339Nano),
	}).Err()
}

func (s *Redis) UpdateLocation(ctx context.Context, driverID string, lat, lon float64) error {
	return s.R.GeoAdd(ctx, geoKeyDrivers, &redis.GeoLocation{
		Name:      driverID,
		Longitude: lon,
		Latitude:  lat,
	}).Err()
}
func (s *Redis) Nearby(ctx context.Context, lat, lon, radiusMeters float64, limit int) ([]domain.Nearby, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	q := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude: lon,
			Latitude:  lat,
			Radius:    radiusMeters,
			Sort:      "ASC",
			Count:     limit,
		},
		WithCoord: true,
		WithDist:  true,
	}
	locs, err := s.R.GeoSearchLocation(ctx, geoKeyDrivers, q).Result()
	if err != nil {
		return nil, err
	}

	out := make([]domain.Nearby, 0, len(locs))
	for _, l := range locs {
		out = append(out, domain.Nearby{DriverID: l.Name, Lat: l.Latitude, Lon: l.Longitude, Meters: l.Dist})
	}
	return out, nil
}

func (s *Redis) Ready(ctx context.Context) error {
	return s.R.Ping(ctx).Err()
}

func ParseLimit(v string, def int) int {
	if v == "" {
		return def
	}
	n, err := strconvAtoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

// tiny local helper (no extra imports in handlers)
func strconvAtoi(s string) (int, error) {
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}
