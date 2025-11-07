package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"net/http"
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
