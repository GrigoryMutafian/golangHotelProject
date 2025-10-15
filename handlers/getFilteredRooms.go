package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/db"
	hm "golangHotelProject/hotelModel"
	"io"
	"log"
	"net/http"
)

func GetFilteredRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var filter map[string]interface{}
	responses := make(map[string][]int)
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil && err != io.EOF {
		log.Printf("JSON decoding error: %v", err)
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(filter) == 0 {
		rows, err := db.DB.Query(`SELECT id, number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning FROM rooms`)
		if err != nil {
			http.Error(w, "DB query error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var rooms []hm.Room

		for rows.Next() {
			var r hm.Room

			err := rows.Scan(&r.ID, &r.Number, &r.RoomCount, &r.IsOccupied, &r.Floor, &r.SleepingPlaces, &r.RoomType, &r.NeedCleaning)
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
		return
	}

	for column, value := range filter {
		query := fmt.Sprintf("SELECT id FROM rooms WHERE %s = $1", column)
		rows, err := db.DB.Query(query, value)

		if err != nil {
			http.Error(w, "String processing error"+err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var currentID int
			rows.Scan(&currentID)

			strValue := fmt.Sprintf("%v", value)
			responses[strValue] = append(responses[strValue], currentID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responses)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
