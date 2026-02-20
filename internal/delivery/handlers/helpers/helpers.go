package helpers

import (
	"encoding/json"
	"golangHotelProject/internal/usecase"
	"log/slog"
	"net/http"
)

func ReqLogger(r *http.Request, handler string) *slog.Logger {
	return slog.Default().With(
		"handler", handler,
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func WriteTextError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func HandleUsecaseError(w http.ResponseWriter, logger *slog.Logger, op string, err error) {
	switch {
	case usecase.IsValidationErr(err):
		logger.Info("validation error", "op", op, "error", err)
		WriteTextError(w, http.StatusBadRequest, err.Error())
	case usecase.IsConflictErr(err):
		logger.Info("conflict error", "op", op, "error", err)
		WriteTextError(w, http.StatusConflict, err.Error())
	default:
		logger.Error("internal error", "op", op, "error", err)
		WriteTextError(w, http.StatusInternalServerError, "internal error: "+err.Error())
	}
}
