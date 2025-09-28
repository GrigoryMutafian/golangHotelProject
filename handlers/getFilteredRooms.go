package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/db"
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
	if err != nil {
		http.Error(w, "JSON decoding error"+err.Error(), http.StatusBadRequest)
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
