package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/model"
	repo "golangHotelProject/internal/repository"
)

type BookingUsecase struct {
	Repo repo.BookingRepository
}

func NewBookingUsecase(repo repo.BookingRepository) *BookingUsecase {
	return &BookingUsecase{Repo: repo}
}
func (uc *BookingUsecase) CreateBooking(ctx context.Context, b model.Booking) error {
	err := validateBooking(b)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	status, _ := uc.Repo.GettingStatus(ctx, b.GuestID)
	if status {
		return errors.Join(ErrValidation, errors.New("already have active booking with this guest_id"))
	}
	ArrivaledRoomBoolean, _ := uc.Repo.ArrivalStatusOfRoom(ctx, b.RoomID)
	if ArrivaledRoomBoolean {
		return errors.Join(ErrValidation, errors.New("already have active booking in this room_id"))
	}
	return uc.Repo.CreateBooking(ctx, b)
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
