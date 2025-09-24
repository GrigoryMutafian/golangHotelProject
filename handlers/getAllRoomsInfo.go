package handlers

import (
	"encoding/json"
	"golangHotelProject/db"
	hm "golangHotelProject/hotelModel"
	"net/http"
)

func GetAllRoomsInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query(`SELECT id, number, room_count, is_occupied, floor, sleeping_places, room_quality, need_cleaning FROM rooms`)
	if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var rooms []hm.Room

	for rows.Next() {
		var r hm.Room

		err := rows.Scan(&r.ID, &r.Number, &r.RoomCount, &r.IsOccupied, &r.Floor, &r.SleepingPlaces, &r.RoomQuality, &r.NeedCleaning)
		if err != nil {
			http.Error(w, "String processing error", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, r)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error as a result of the request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
