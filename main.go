package main

import (
	"fmt"
	"golangHotelProject/db"
	hn "golangHotelProject/handlers"
	"log"
	"net/http"
)

func main() {
	if err := db.InitDB(); err != nil {
		fmt.Println(err)
		return
	}

	defer db.DB.Close()

	http.HandleFunc("/AddRoom", hn.AddRoom)
	http.HandleFunc("/GetAllRoomsInfo", hn.GetAllRoomsInfo)
	http.HandleFunc("/RemoveRoom", hn.RemoveRoom)
	http.HandleFunc("/PatchRoom", hn.PatchRoom)
	http.HandleFunc("/GetFilteredRooms", hn.GetFilteredRooms)
	log.Println("server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("the server is not running", err)
	}
}
