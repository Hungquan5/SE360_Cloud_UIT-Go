package driver
// services/driver/main.go
package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "strconv"

  "github.com/go-chi/chi/v5"
  "github.com/redis/go-redis/v9"
)

var (
  rdb   *redis.Client
  ctx   = context.Background()
  GEO_KEY = "drivers:geo"
)

func main() {
  rdb = redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
  r := chi.NewRouter()

  // PUT /drivers/{id}/location
  r.Put("/drivers/{id}/location", func(w http.ResponseWriter, req *http.Request) {
    id := chi.URLParam(req, "id")
    latStr := req.URL.Query().Get("lat")
    lngStr := req.URL.Query().Get("lng")
    lat, _ := strconv.ParseFloat(latStr, 64)
    lng, _ := strconv.ParseFloat(lngStr, 64)

    if err := rdb.GeoAdd(ctx, GEO_KEY, &redis.GeoLocation{
      Name: id, Latitude: lat, Longitude: lng,
    }).Err(); err != nil {
      http.Error(w, err.Error(), 500); return
    }
    w.WriteHeader(http.StatusNoContent)
  })

  // GET /drivers/search?lat=&lng=&radius=3000
  r.Get("/drivers/search", func(w http.ResponseWriter, req *http.Request) {
    lat,_ := strconv.ParseFloat(req.URL.Query().Get("lat"),64)
    lng,_ := strconv.ParseFloat(req.URL.Query().Get("lng"),64)
    radius := req.URL.Query().Get("radius")
    if radius == "" { radius = "3000" } // meters

    res, err := rdb.GeoRadius(ctx, GEO_KEY, lng, lat, &redis.GeoRadiusQuery{
      Radius: mustFloat(radius), Unit: "m", WithDist: true, Count: 10, Sort: "ASC",
    }).Result()
    if err != nil { http.Error(w, err.Error(), 500); return }
    writeJSON(w, res)
  })

  log.Println("driver-svc on :8080")
  http.ListenAndServe(":8080", r)
}

func mustFloat(s string) float64 { v,_ := strconv.ParseFloat(s,64); return v }
