package handlers

import (
	"encoding/json"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	md "golangHotelProject/internal/model"
	"golangHotelProject/internal/usecase"
	"log/slog"
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

// @Summary create room
// @Tags room
// @Description create room
// @ID createRoom
// @Accept json
// @Produce json
// @Param input body md.Room true "new room data"
// @Success 201 {object} dto.CreatingRoomResponse "Created"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON or validation error"
// @Failure 409 {object} dto.ErrorResponse "Conflict (room already exists)"
// @Failure 413 {object} dto.ErrorResponse "Request entity too large"
// @Failure 415 {object} dto.ErrorResponse "Unsupported media type"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /Create [post]
func Create(w http.ResponseWriter, r *http.Request) {
	log := slog.Default().With(
		"handler", "room.create",
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

	var NewRoom md.Room

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&NewRoom)
	if err != nil {
		log.Warn("invalid json",
			"error", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Info("decoded room",
			"room number", NewRoom.Number)
	}

	if err := roomUC.AddRoom(r.Context(), NewRoom); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error",
				"room number", NewRoom.Number,
				"error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("room conflict",
				"room number", NewRoom.Number,
				"err", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("add room failed",
				"room_number", NewRoom.Number,
				"err", err,
			)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("room added",
		"room_number", NewRoom.Number)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]int{"Number of added Room is": NewRoom.Number}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("JSON encode error",
			"error", err,
			"room number", NewRoom.Number)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
	log.Info("response sent",
		"status", http.StatusCreated,
		"room_number", NewRoom.Number)
}

// @Summary patch room
// @Tags room
// @Description patch an existing room
// @ID patchRoom
// @Accept json
// @Produce json
// @Param id query int true "room id"
// @Param input body dto.RoomPatch true "patch data"
// @Success 200 {object} dto.RoomPatchResponse "rooms updated"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON or validation error"
// @Failure 409 {object} dto.ErrorResponse "Conflict"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /Patch [patch]
func Patch(w http.ResponseWriter, r *http.Request) {
	log := slog.Default().With(
		"handler", "room.patch",
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

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		log.Warn("missing id")
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn("invalid id", "id", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var patch dto.RoomPatch

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&patch)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("patching room", "room_id", id)

	if err = roomUC.PatchRoom(r.Context(), id, patch); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "room_id", id, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("room conflict", "room_id", id, "err", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("patch room failed", "room_id", id, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("room patched", "room_id", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "rooms updated"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("JSON encode error", "error", err, "room_id", id)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response sent", "status", http.StatusOK, "room_id", id)
}

// @Summary remove room
// @Tags room
// @Description remove an existing room by id
// @ID removeRoom
// @Accept json
// @Produce json
// @Param input body dto.RemoveRoomRequest true "room id to remove"
// @Success 200 {string} string "Removed Room id: {id}"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON or validation error"
// @Failure 409 {object} dto.ErrorResponse "Conflict"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /RemoveRoom [delete]
func RemoveRoom(w http.ResponseWriter, r *http.Request) {
	log := slog.Default().With(
		"handler", "room.remove",
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

	var romovingRoomID int

	err := json.NewDecoder(r.Body).Decode(&romovingRoomID)
	if err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, "JSON encoding error"+err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("removing room", "room_id", romovingRoomID)

	if err = roomUC.RemoveRoom(r.Context(), romovingRoomID); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "room_id", romovingRoomID, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("room conflict", "room_id", romovingRoomID, "err", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("remove room failed", "room_id", romovingRoomID, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("room removed", "room_id", romovingRoomID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	removedRoom := fmt.Sprintf("Removed Room id: %d", romovingRoomID)
	err = json.NewEncoder(w).Encode(removedRoom)
	if err != nil {
		log.Error("JSON encode error", "error", err, "room_id", romovingRoomID)
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response sent", "status", http.StatusOK, "room_id", romovingRoomID)
}

// @Summary get filtered rooms
// @Tags room
// @Description get list of rooms with optional filters
// @ID getFilteredRooms
// @Accept json
// @Produce json
// @Param input body map[string]interface{} false "filter criteria"
// @Success 200 {array} md.Room "list of rooms"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON or validation error"
// @Failure 409 {object} dto.ErrorResponse "Conflict"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /GetFilteredRooms [post]
func GetFilteredRooms(w http.ResponseWriter, r *http.Request) {
	log := slog.Default().With(
		"handler", "room.getFilteredRooms",
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

	log.Info("getting filtered rooms", "filter", filter)

	if len(filter) == 0 {
		log.Info("getting all rooms")
		rooms, err := roomUC.GetList(r.Context())

		if err != nil {
			switch {
			case usecase.IsValidationErr(err):
				log.Info("validation error", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			case usecase.IsConflictErr(err):
				log.Info("room conflict", "error", err)
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				log.Error("get rooms failed", "err", err)
				http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		log.Info("rooms retrieved", "count", len(rooms))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(rooms)
		if err != nil {
			log.Error("JSON encode error", "error", err)
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
		log.Info("response sent", "status", http.StatusOK, "count", len(rooms))
		return
	}

	responses, err := roomUC.GetFilteredRooms(r.Context(), filter)
	if err != nil {
		switch {
		case usecase.IsValidationErr(err):
			log.Info("validation error", "filter", filter, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			log.Info("room conflict", "filter", filter, "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Error("get filtered rooms failed", "filter", filter, "err", err)
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Info("filtered rooms retrieved", "filter", filter, "count", len(responses))
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
