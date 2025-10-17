package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/db"
	"net/http"
)

func PatchRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	replacementColumn := map[int]map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&replacementColumn)
	if err != nil {
		http.Error(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
		return
	}
	for id, columns := range replacementColumn {
		for columnName, columnValue := range columns {
			query := fmt.Sprintf("UPDATE rooms SET %s = $1 WHERE id = $2", columnName)
			result, err := db.DB.Exec(query, columnValue, id)

			if err != nil {
				http.Error(w, "database manipulating error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				http.Error(w, "error checking rows affected: "+err.Error(), http.StatusInternalServerError)
			}

			if rowsAffected == 0 {
				http.Error(w, "no rows updated"+err.Error(), http.StatusNotFound)
				return
			}
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "rooms " + "updated"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
