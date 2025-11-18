package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"net/http"
	"strconv"
)

var bookingUC *usecase.BookingUsecase

func InitBookingDependencies(uc *usecase.BookingUsecase) error {
	if uc == nil {
		return fmt.Errorf("nil usecase")
	}
	bookingUC = uc
	return nil
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var NewBooking model.Booking

	err := json.NewDecoder(r.Body).Decode(&NewBooking)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := bookingUC.CreateBooking(r.Context(), NewBooking); err != nil {
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
	response := "Booking created"
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
}

func ReadBookingByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id input is clear", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "pars error", http.StatusBadRequest)
		return
	}

	book, err := bookingUC.ReadByIDUsecase(r.Context(), idInt)
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

	text := fmt.Sprintf("column id: %d", idInt)
	response := map[string]model.Booking{text: book}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func PatchBookingByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var patch dto.BookingPatch

	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := bookingUC.PatchBookingByID(r.Context(), patch); err != nil {
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
	response := fmt.Sprintf("column id: %d", patch.ID)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetFilteredBookings(w http.ResponseWriter, r *http.Request) {
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
		bookings, err := roomUC.GetList(r.Context())

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
		err = json.NewEncoder(w).Encode(bookings)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
		return
	}
	responses, err := bookingUC.GetFilteredBookings(r.Context(), filter)
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

func RemoveBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var romovingBookingID int

	err := json.NewDecoder(r.Body).Decode(&romovingBookingID)
	if err != nil {
		http.Error(w, "JSON encoding error"+err.Error(), http.StatusBadRequest)
	}

	if err = bookingUC.RemoveBooking(r.Context(), romovingBookingID); err != nil {
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
	removedRoom := fmt.Sprintf("Removed Coffee id: %d", romovingBookingID)
	err = json.NewEncoder(w).Encode(removedRoom)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
