package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/logger"
	"golangHotelProject/internal/model"
	repo "golangHotelProject/internal/repository"
	"log/slog"
)

type BookingUsecase struct {
	Repo   repo.BookingRepository
	Logger *slog.Logger
}

func NewBookingUsecase(repo repo.BookingRepository, log logger.Logger) *BookingUsecase {
	return &BookingUsecase{
		Repo:   repo,
		Logger: log.With("component", "BookingUsecase"),
	}
}

func (uc *BookingUsecase) CreateBooking(ctx context.Context, b model.Booking) error {
	const op = "CreateBooking"

	uc.Logger.Debug("creating booking",
		"op", op,
		"room_id", b.RoomID,
		"guest_id", b.GuestID,
	)

	err := validateBooking(b)
	if err != nil {
		uc.Logger.Warn("booking validation failed",
			"op", op,
			"room_id", b.RoomID,
			"guest_id", b.GuestID,
			"error", err.Error(),
		)
		return errors.Join(ErrValidation, err)
	}

	status, _ := uc.Repo.GettingStatus(ctx, b.GuestID)
	if status {
		uc.Logger.Warn("guest already has active booking",
			"op", op,
			"guest_id", b.GuestID,
		)
		return errors.Join(ErrValidation, errors.New("already have active booking with this guest_id"))
	}

	ArrivaledRoomBoolean, _ := uc.Repo.ArrivalStatusOfRoom(ctx, b.RoomID)
	if ArrivaledRoomBoolean {
		uc.Logger.Warn("room already booked",
			"op", op,
			"room_id", b.RoomID,
		)
		return errors.Join(ErrValidation, errors.New("already have active booking in this room_ID"))
	}

	if err := uc.Repo.CreateBooking(ctx, b); err != nil {
		uc.Logger.Error("failed to create booking",
			"op", op,
			"room_id", b.RoomID,
			"guest_id", b.GuestID,
			"error", err.Error(),
		)
		return err
	}

	uc.Logger.Info("booking created successfully",
		"op", op,
		"room_id", b.RoomID,
		"guest_id", b.GuestID,
	)
	return nil
}

func validateBooking(b model.Booking) error {
	if b.RoomID <= 0 {
		return errors.Join(ErrValidation, errors.New("room id <= 0"))
	}
	if b.GuestID <= 0 {
		return errors.Join(ErrValidation, errors.New("guest id <= 0"))
	}
	if b.Start_date.IsZero() || b.End_date.IsZero() {
		return errors.New("some date is clear")
	}
	if !b.Start_date.Before(b.End_date) {
		return errors.New("start_date должен быть раньше end_date")
	}
	return nil
}

func (uc *BookingUsecase) ReadByIDUsecase(ctx context.Context, id int) (model.Booking, error) {
	const op = "ReadByIDUsecase"

	uc.Logger.Debug("reading booking by id",
		"op", op,
		"booking_id", id,
	)

	if id <= 0 {
		uc.Logger.Warn("invalid booking id",
			"op", op,
			"booking_id", id,
		)
		return model.Booking{}, errors.Join(ErrValidation, errors.New("id <= 0"))
	}

	b, _ := uc.Repo.ReadBookingByID(ctx, id)
	if b.RoomID == 0 && b.GuestID == 0 {
		uc.Logger.Warn("booking not found",
			"op", op,
			"booking_id", id,
		)
		return model.Booking{}, errors.Join(ErrValidation, errors.New("no rows"))
	}

	uc.Logger.Debug("booking retrieved successfully",
		"op", op,
		"booking_id", id,
	)
	return b, nil
}

func (uc *BookingUsecase) PatchBookingByID(ctx context.Context, b dto.BookingPatch) error {
	const op = "PatchBookingByID"

	uc.Logger.Debug("patching booking",
		"op", op,
		"booking_id", *b.ID,
	)

	if b.ID == nil || *b.ID <= 0 {
		uc.Logger.Warn("invalid booking id",
			"op", op,
		)
		return errors.Join(ErrValidation, errors.New("id <= 0"))
	}

	old, _ := uc.Repo.ReadBookingByID(ctx, *b.ID)
	if b.RoomID == nil {
		b.RoomID = &old.RoomID
	}
	if b.GuestID == nil {
		b.GuestID = &old.GuestID
	}
	if b.Start_date == nil {
		b.Start_date = &old.Start_date
	}
	if b.End_date == nil {
		b.End_date = &old.End_date
	}
	if b.Status == nil {
		b.Status = &old.Status
	}

	err := validateBookingPatch(b)
	if err != nil {
		uc.Logger.Warn("booking patch validation failed",
			"op", op,
			"booking_id", *b.ID,
			"error", err.Error(),
		)
		return errors.Join(ErrValidation, err)
	}

	err = uc.Repo.PatchBooking(ctx, b)
	if err != nil {
		uc.Logger.Error("failed to patch booking",
			"op", op,
			"booking_id", *b.ID,
			"error", err.Error(),
		)
		return errors.Join(errors.New("DB manipulating error"), err)
	}

	uc.Logger.Info("booking patched successfully",
		"op", op,
		"booking_id", *b.ID,
	)
	return nil
}

func validateBookingPatch(b dto.BookingPatch) error {
	if b.RoomID == nil || *b.RoomID <= 0 {
		return errors.Join(ErrValidation, errors.New("room id <= 0"))
	}
	if b.GuestID == nil || *b.GuestID <= 0 {
		return errors.Join(ErrValidation, errors.New("guest id <= 0"))
	}
	if b.Start_date == nil || b.Start_date.IsZero() || b.End_date == nil || b.End_date.IsZero() {
		return errors.New("some date is clear")
	}
	if !b.Start_date.Before(*b.End_date) {
		return errors.New("start_date должен быть раньше end_date")
	}
	return nil
}

func (uc *BookingUsecase) GetList(ctx context.Context) ([]model.Booking, error) {
	const op = "GetList"

	uc.Logger.Debug("fetching booking list", "op", op)

	response, _ := uc.Repo.ListColumn(ctx)
	if len(response) == 0 {
		uc.Logger.Info("booking list is empty", "op", op)
		return response, errors.Join(ErrConflict, errors.New("database is clear"))
	}

	uc.Logger.Debug("booking list fetched successfully",
		"op", op,
		"count", len(response),
	)
	return response, nil
}

func (uc *BookingUsecase) GetFilteredBookings(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	const op = "GetFilteredBookings"

	uc.Logger.Debug("filtering bookings",
		"op", op,
		"filter", filter,
	)

	response, err := uc.Repo.FilterBookings(ctx, filter)
	if err != nil {
		uc.Logger.Error("failed to filter bookings",
			"op", op,
			"filter", filter,
			"error", err.Error(),
		)
		return nil, err
	}

	uc.Logger.Debug("bookings filtered successfully",
		"op", op,
		"filter", filter,
		"results_count", len(response),
	)
	return response, err
}

func (uc *BookingUsecase) RemoveBooking(ctx context.Context, id int) error {
	const op = "RemoveBooking"

	uc.Logger.Debug("removing booking",
		"op", op,
		"booking_id", id,
	)

	if id <= 0 {
		uc.Logger.Warn("invalid booking id",
			"op", op,
			"booking_id", id,
		)
		return errors.Join(ErrValidation, errors.New("ID must be more than 0"))
	}

	err := uc.Repo.DeleteBooking(ctx, id)
	if err != nil {
		uc.Logger.Error("failed to delete booking",
			"op", op,
			"booking_id", id,
			"error", err.Error(),
		)
		return err
	}

	uc.Logger.Info("booking removed successfully",
		"op", op,
		"booking_id", id,
	)
	return nil
}
