package main

import (
	"fmt"
	_ "golangHotelProject/docs"
	hn "golangHotelProject/internal/delivery/handlers"
	"golangHotelProject/internal/repository"
	"golangHotelProject/internal/repository/db"
	"golangHotelProject/internal/usecase"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Hotel Booking API
// @version 1.0
// @description API server for hotel actions

// @host localhost:8080
// @BasePath /

func withCORS(next http.Handler) http.Handler {
	allowed := map[string]bool{
		"http://localhost:3000": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := db.InitDB(); err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	roomRepo := &repository.PgRoomRepository{DB: db.DB}
	roomUC := usecase.NewRoomUsecase(roomRepo)
	if err := hn.InitDependencies(roomUC); err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	bookingRepo := &repository.PgBookingRepository{DB: db.DB}
	bookingUC := usecase.NewBookingUsecase(bookingRepo)

	if err := hn.InitBookingDependencies(bookingUC); err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	http.HandleFunc("/Create", hn.Create)
	http.HandleFunc("/RemoveRoom", hn.RemoveRoom)
	http.HandleFunc("/Patch", hn.Patch)
	http.HandleFunc("/GetFilteredRooms", hn.GetFilteredRooms)

	http.HandleFunc("/CreateBooking", hn.CreateBooking)
	http.HandleFunc("/ReadBookingByID", hn.ReadBookingByID)
	http.HandleFunc("/PatchBookingByID", hn.PatchBookingByID)
	http.HandleFunc("/RemoveBooking", hn.RemoveBooking)
	http.HandleFunc("/GetFilteredBookings", hn.GetFilteredBookings)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			log.Printf("Error writing health response: %v", err)
		}
	})

	log.Println("server running on http://localhost:8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	log.Println("Health check available at http://localhost:8080/health")

	handler := withCORS(http.DefaultServeMux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Println("the server is not running", err)
	}
}
