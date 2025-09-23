package handlers

import (
	"encoding/json"
	"golangHotelProject/db"
	md "golangHotelProject/hotelModel"
	"net/http"
)

func AddRoom(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var newRoom md.Room

	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`INSERT INTO rooms (number, room_count, is_occupied, floor, sleeping_places, room_quality, need_cleaning)
	VALUES($1, $2, $3, $4, $5, $6, $7)`, newRoom.Number, newRoom.RoomCount, newRoom.IsOccupied, newRoom.Floor, newRoom.SleepingPlaces, newRoom.RoomQuality, newRoom.NeedCleaning)

	if err != nil {
		http.Error(w, "Database insertion error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking rows affected: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No rows inserted", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]int{"Number of added Room is": newRoom.Number}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}

}
