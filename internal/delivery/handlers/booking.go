package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"log/slog"
	"net/http"
	"strconv"
)

func reqLogger(r *http.Request, handler string) *slog.Logger {
	return slog.Default().With(
		"handler", handler,
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func writeTextError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

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
	log := slog.Default().With(
		"handler", "booking.create",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	if r.Method != http.MethodPost {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error("error closing request body", "err", err)
		}
	}()

	var NewBooking model.Booking

	err := json.NewDecoder(r.Body).Decode(&NewBooking)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("creating booking", "booking_id", NewBooking.ID)

	if err := bookingUC.CreateBooking(r.Context(), NewBooking); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "booking_id", NewBooking.ID, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("booking conflict", "booking_id", NewBooking.ID, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("create booking failed", "booking_id", NewBooking.ID, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("booking created", "booking_id", NewBooking.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := "Booking created"
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", NewBooking.ID)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
	log.Info("response sent", "status", http.StatusCreated, "booking_id", NewBooking.ID)
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
	log := slog.Default().With(
		"handler", "booking.readByID",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	if r.Method != http.MethodGet {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		log.Warn("missing id")
		http.Error(w, "id input is clear", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		log.Warn("invalid id", "id", idStr)
		http.Error(w, "pars error", http.StatusBadRequest)
		return
	}

	log.Info("reading booking", "booking_id", idInt)

	book, err := bookingUC.ReadByIDUsecase(r.Context(), idInt)
	if err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "booking_id", idInt, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("booking conflict", "booking_id", idInt, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("read booking failed", "booking_id", idInt, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("booking retrieved", "booking_id", idInt)

	text := fmt.Sprintf("column id: %d", idInt)
	response := map[string]model.Booking{text: book}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", idInt)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Info("response sent", "status", http.StatusOK, "booking_id", idInt)
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
	log := slog.Default().With(
		"handler", "booking.patch",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	if r.Method != http.MethodPatch {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error("error closing request body", "err", err)
		}
	}()

	var patch dto.BookingPatch

	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("patching booking", "booking_id", patch.ID)

	if err := bookingUC.PatchBookingByID(r.Context(), patch); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "booking_id", patch.ID, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("booking conflict", "booking_id", patch.ID, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("patch booking failed", "booking_id", patch.ID, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("booking patched", "booking_id", patch.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("column id: %d", patch.ID)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", patch.ID)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response sent", "status", http.StatusOK, "booking_id", patch.ID)
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
	log := slog.Default().With(
		"handler", "booking.getFiltered",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	if r.Method != http.MethodGet {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error("error closing request body", "err", err)
		}
	}()

	var filter map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("getting filtered bookings", "filter", filter)

	if len(filter) == 0 {
		log.Info("getting all bookings")
		bookings, err := bookingUC.GetList(r.Context())

		if err != nil {
			switch {
			case usecase.IsValidationErr(err):
				log.Info("validation error", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			case usecase.IsConflictErr(err):
				log.Info("booking conflict", "error", err)
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				log.Error("get bookings failed", "err", err)
				http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		log.Info("bookings retrieved", "count", len(bookings))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(bookings)
		if err != nil {
			log.Error("JSON encode error", "error", err)
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
		log.Info("response sent", "status", http.StatusOK, "count", len(bookings))
		return
	}

	responses, err := bookingUC.GetFilteredBookings(r.Context(), filter)
	if err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "filter", filter, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("booking conflict", "filter", filter, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("get filtered bookings failed", "filter", filter, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("filtered bookings retrieved", "filter", filter, "count", len(responses))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responses)
	if err != nil {
		log.Error("JSON encode error", "error", err, "filter", filter)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response sent", "status", http.StatusOK, "filter", filter, "count", len(responses))
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
	log := slog.Default().With(
		"handler", "booking.remove",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	if r.Method != http.MethodDelete {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error("error closing request body", "err", err)
		}
	}()

	var removingBookingID int

	err := json.NewDecoder(r.Body).Decode(&removingBookingID)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "JSON encoding error"+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("removing booking", "booking_id", removingBookingID)

	if err = bookingUC.RemoveBooking(r.Context(), removingBookingID); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "booking_id", removingBookingID, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("booking conflict", "booking_id", removingBookingID, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("remove booking failed", "booking_id", removingBookingID, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("booking removed", "booking_id", removingBookingID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	removedBooking := fmt.Sprintf("Removed Coffee id: %d", removingBookingID)
	err = json.NewEncoder(w).Encode(removedBooking)
	if err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", removingBookingID)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response sent", "status", http.StatusOK, "booking_id", removingBookingID)
}
