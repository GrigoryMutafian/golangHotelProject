package main

import (
	"fmt"
	hn "golangHotelProject/internal/delivery/handlers"
	"golangHotelProject/internal/repository"
	"golangHotelProject/internal/repository/db"
	"golangHotelProject/internal/usecase"
	"log"
	"net/http"
)

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
	log.Println("server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("the server is not running", err)
	}
}
