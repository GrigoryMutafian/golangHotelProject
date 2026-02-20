package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/delivery/handlers/helpers"
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
	log := helpers.ReqLogger(r, "booking.create")

	if r.Method != http.MethodPost {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		helpers.WriteTextError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
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
		helpers.WriteTextError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	log.Info("creating booking", "booking_id", NewBooking.ID)

	if err := bookingUC.CreateBooking(r.Context(), NewBooking); err != nil {
		helpers.HandleUsecaseError(w, log, "create booking", err)
		return
	}

	log.Info("booking created", "booking_id", NewBooking.ID)

	response := "Booking created"
	if err := helpers.WriteJSON(w, http.StatusCreated, response); err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", NewBooking.ID)
		helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error: "+err.Error())
		return
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
	log := helpers.ReqLogger(r, "booking.readByID")

	if r.Method != http.MethodGet {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		helpers.WriteTextError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		log.Warn("missing id")
		helpers.WriteTextError(w, http.StatusBadRequest, "id input is clear")
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		log.Warn("invalid id", "id", idStr)
		helpers.WriteTextError(w, http.StatusBadRequest, "pars error")
		return
	}

	log.Info("reading booking", "booking_id", idInt)

	book, err := bookingUC.ReadByIDUsecase(r.Context(), idInt)
	if err != nil {
		helpers.HandleUsecaseError(w, log, "reading booking", err)
		return
	}

	log.Info("booking retrieved", "booking_id", idInt)

	text := fmt.Sprintf("column id: %d", idInt)
	response := map[string]model.Booking{text: book}
	if err := helpers.WriteJSON(w, http.StatusOK, response); err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", idInt)
		helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error: "+err.Error())
		return
	}
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
	log := helpers.ReqLogger(r, "booking.patch")

	if r.Method != http.MethodPatch {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		helpers.WriteTextError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		helpers.WriteTextError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	log.Info("patching booking", "booking_id", patch.ID)

	if err := bookingUC.PatchBookingByID(r.Context(), patch); err != nil {
		helpers.HandleUsecaseError(w, log, "patch booking", err)
		return
	}

	log.Info("booking patched", "booking_id", patch.ID)

	response := fmt.Sprintf("column id: %d", patch.ID)
	if err := helpers.WriteJSON(w, http.StatusOK, response); err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", patch.ID)
		helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error: "+err.Error())
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
	log := helpers.ReqLogger(r, "booking.getFiltered")

	if r.Method != http.MethodGet {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		helpers.WriteTextError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		helpers.WriteTextError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	log.Info("getting filtered bookings", "filter", filter)

	if len(filter) == 0 {
		log.Info("getting all bookings")
		bookings, err := bookingUC.GetList(r.Context())

		if err != nil {
			helpers.HandleUsecaseError(w, log, "get all bookings", err)
			return
		}

		log.Info("bookings retrieved", "count", len(bookings))
		if err := helpers.WriteJSON(w, http.StatusOK, bookings); err != nil {
			log.Error("JSON encode error", "error", err)
			helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error")
			return
		}
		log.Info("response sent", "status", http.StatusOK, "count", len(bookings))
		return
	}

	responses, err := bookingUC.GetFilteredBookings(r.Context(), filter)
	if err != nil {
		helpers.HandleUsecaseError(w, log, "get filtered bookings", err)
		return
	}

	log.Info("filtered bookings retrieved", "filter", filter, "count", len(responses))
	if err := helpers.WriteJSON(w, http.StatusOK, responses); err != nil {
		log.Error("JSON encode error", "error", err, "filter", filter)
		helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error: "+err.Error())
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
	log := helpers.ReqLogger(r, "booking.remove")

	if r.Method != http.MethodDelete {
		log.Warn(
			"method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
		)
		helpers.WriteTextError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		helpers.WriteTextError(w, http.StatusBadRequest, "JSON encoding error: "+err.Error())
		return
	}

	log.Info("removing booking", "booking_id", removingBookingID)

	if err = bookingUC.RemoveBooking(r.Context(), removingBookingID); err != nil {
		helpers.HandleUsecaseError(w, log, "remove booking", err)
		return
	}

	log.Info("booking removed", "booking_id", removingBookingID)

	removedBooking := fmt.Sprintf("Removed Booking id: %d", removingBookingID)
	if err := helpers.WriteJSON(w, http.StatusOK, removedBooking); err != nil {
		log.Error("JSON encode error", "error", err, "booking_id", removingBookingID)
		helpers.WriteTextError(w, http.StatusInternalServerError, "JSON encoding error: "+err.Error())
		return
	}
	log.Info("response sent", "status", http.StatusOK, "booking_id", removingBookingID)
}
