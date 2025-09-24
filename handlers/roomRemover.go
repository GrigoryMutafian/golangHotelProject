package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/db"
	"net/http"
)

func RemoveRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var roomRemover int

	err := json.NewDecoder(r.Body).Decode(&roomRemover)
	if err != nil {
		http.Error(w, "JSON encoding error"+err.Error(), http.StatusBadRequest)
	}

	result, err := db.DB.Exec(`DELETE FROM rooms WHERE id = $1`, roomRemover)

	if err != nil {
		http.Error(w, "Database manipulating error: "+err.Error(), http.StatusInternalServerError)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking rows affected: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No rows inserted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	removedRoom := fmt.Sprintf("Removed Coffee id: %d", roomRemover)
	err = json.NewEncoder(w).Encode(removedRoom)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
