package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"log"
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

// CreateBooking creates a new booking
// @Summary Create a new booking
// @Description Create a booking for a room
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body model.Booking true "Booking object"
// @Success 201 {object} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /CreateBooking [post]
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body: %v", err)
		}
	}()

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

// ReadBookingByID returns booking by ID
// @Summary Get booking by ID
// @Description Retrieve a specific booking by its ID
// @Tags bookings
// @Produce json
// @Param id query int true "Booking ID"
// @Success 200 {object} map[string]model.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ReadBookingByID [get]
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

// PatchBookingByID updates booking
// @Summary Update booking details
// @Description Update an existing booking with partial data
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.BookingPatch true "Booking object with updates"
// @Success 200 {object} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /PatchBookingByID [patch]
func PatchBookingByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body: %v", err)
		}
	}()

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

// GetFilteredBookings returns filtered bookings
// @Summary Get filtered bookings
// @Description Get bookings by filter parameters
// @Tags bookings
// @Accept json
// @Produce json
// @Param filter body map[string]interface{} false "Filter criteria"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /GetFilteredBookings [get]
func GetFilteredBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body: %v", err)
		}
	}()

	var filter map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(filter) == 0 {
		bookings, err := bookingUC.GetList(r.Context())

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
	err = json.NewEncoder(w).Encode(responses)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// RemoveBooking deletes a booking
// @Summary Delete a booking
// @Description Delete a booking by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id body int true "Booking ID"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /RemoveBooking [delete]
func RemoveBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body: %v", err)
		}
	}()

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
