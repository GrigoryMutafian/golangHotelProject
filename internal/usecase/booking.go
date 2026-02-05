package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/delivery/handlers/dto"
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
		return errors.Join(ErrValidation, errors.New("already have active booking in this room_ID"))
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

func (uc *BookingUsecase) ReadByIDUsecase(ctx context.Context, id int) (model.Booking, error) {
	if id <= 0 {
		return model.Booking{}, errors.Join(ErrValidation, errors.New("id <= 0"))
	}
	b, _ := uc.Repo.ReadBookingByID(ctx, id)
	if b.RoomID == 0 && b.GuestID == 0 {
		return model.Booking{}, errors.Join(ErrValidation, errors.New("no rows"))
	}
	return b, nil
}

func (uc *BookingUsecase) PatchBookingByID(ctx context.Context, b dto.BookingPatch) error {
	if b.ID == nil || *b.ID <= 0 {
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
		return errors.Join(ErrValidation, err)
	}
	err = uc.Repo.PatchBooking(ctx, b)
	if err != nil {
		return errors.Join(errors.New("DB manipulating error"), err)
	}
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
	response, _ := uc.Repo.ListColumn(ctx)
	if len(response) == 0 {
		return response, errors.Join(ErrConflict, errors.New("database is clear"))
	}
	return response, nil
}

func (uc *BookingUsecase) GetFilteredBookings(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	response, err := uc.Repo.FilterBookings(ctx, filter)
	return response, err
}

func (uc *BookingUsecase) RemoveBooking(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("ID must be more than 0"))
	}
	err := uc.Repo.DeleteBooking(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
