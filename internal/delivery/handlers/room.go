package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	md "golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"net/http"
	"strconv"
)

var roomUC *usecase.RoomUsecase

func InitDependencies(uc *usecase.RoomUsecase) error {
	if uc == nil {
		return fmt.Errorf("nil usecase")
	}
	roomUC = uc
	return nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var NewRoom md.Room

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&NewRoom)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := roomUC.AddRoom(r.Context(), NewRoom); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]int{"Number of added Room is": NewRoom.Number}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}

}

func Patch(w http.ResponseWriter, r *http.Request) {

	if roomUC == nil {
		http.Error(w, "handler not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var patch dto.RoomPatch

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&patch)
	if err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = roomUC.PatchRoom(r.Context(), id, patch); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "rooms updated"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var romovingRoomID int

	err := json.NewDecoder(r.Body).Decode(&romovingRoomID)
	if err != nil {
		http.Error(w, "JSON encoding error"+err.Error(), http.StatusBadRequest)
	}

	if err = roomUC.RemoveRoom(r.Context(), romovingRoomID); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	removedRoom := fmt.Sprintf("Removed Coffee id: %d", romovingRoomID)
	err = json.NewEncoder(w).Encode(removedRoom)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetFilteredRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var filter map[string]interface{} //string - columns interface{} - values, getting ids with same parametrs
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(filter) == 0 {
		rooms, err := roomUC.GetList(r.Context())

		if err != nil {
			switch {
			case usecase.IsValidationErr(err):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case usecase.IsConflictErr(err):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			}
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
	responses, err := roomUC.GetFilteredRooms(r.Context(), filter)
	if err != nil {
		switch {
		case usecase.IsValidationErr(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := "value[ids]"
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(responses)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
