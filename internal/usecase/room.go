package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/delivery/handlers/dto"
	md "golangHotelProject/internal/model"
	repo "golangHotelProject/internal/repository"
)

var (
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict error")
)

func IsValidationErr(err error) bool { return errors.Is(err, ErrValidation) }
func IsConflictErr(err error) bool   { return errors.Is(err, ErrConflict) }

type RoomUsecase struct {
	Repo repo.RoomRepository
}

func NewRoomUsecase(repo repo.RoomRepository) *RoomUsecase {
	return &RoomUsecase{Repo: repo}
}

func (uc *RoomUsecase) AddRoom(ctx context.Context, room md.Room) error {
	err := validateRoom(room)

	if err != nil {
		return errors.Join(ErrValidation, err)
	}

	exists, err := uc.Repo.IsNumberExists(ctx, room.Number)
	if err != nil {
		return err
	}
	if exists {
		return errors.Join(ErrConflict, errors.New("room number already exists"))
	}

	if err := uc.Repo.CreateRoom(ctx, room); err != nil {
		return err
	}
	return nil
}

func validateRoom(r md.Room) error {
	if r.Number <= 0 {
		return errors.Join(ErrValidation, errors.New("number must be more then 0"))
	}
	if r.RoomCount < 1 {
		return errors.Join(ErrValidation, errors.New("room count must be more then 0"))
	}
	if r.SleepingPlaces < 1 {
		return errors.Join(ErrValidation, errors.New("sleepng places must be more then 0"))
	}

	if r.Floor < 1 {
		return errors.Join(ErrValidation, errors.New("there no underground floors, number must more then 0"))
	}
	switch r.RoomType {
	case "Standard", "Deluxe", "Suite":
	default:
		return errors.Join(ErrValidation, errors.New("room_type must be one of: Standard, Deluxe, Suite"))
	}
	if r.IsOccupied && r.NeedCleaning {
		return errors.Join(ErrValidation, errors.New("cannot set need_cleaning while room is occupied"))
	}
	return nil
}

func (uc *RoomUsecase) PatchRoom(ctx context.Context, id int, p dto.RoomPatch) error {

	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("invalid id"))
	}

	if isEmptyPatch(p) {
		return nil
	}

	if p.Floor != nil && *p.Floor <= 0 {
		return errors.Join(ErrValidation, errors.New("floor must be more then 0"))
	}

	if p.SleepingPlaces != nil && *p.SleepingPlaces <= 0 {
		return errors.Join(ErrValidation, errors.New("sleepng places must be more then 0"))
	}
	if p.RoomType != nil {
		switch *p.RoomType {
		case "Standard", "Deluxe", "Suite":
		default:
			return errors.Join(ErrValidation, errors.New("room_type must be one of: Standard, Deluxe, Suite"))
		}
	}
	if p.RoomCount != nil && *p.RoomCount <= 0 {
		return errors.Join(ErrValidation, errors.New("room count must be more then 0"))

	}
	if p.IsOccupied != nil && p.NeedCleaning != nil {
		if *p.IsOccupied && *p.NeedCleaning {
			return errors.Join(ErrValidation, errors.New("cannot set need_cleaning while room is occupied"))
		}
	}

	if err := uc.Repo.PatchRoom(ctx, id, p); err != nil {
		return err
	}
	return nil
}

func isEmptyPatch(p dto.RoomPatch) bool {
	return p.RoomCount == nil &&
		p.IsOccupied == nil &&
		p.Floor == nil &&
		p.SleepingPlaces == nil &&
		p.RoomType == nil &&
		p.NeedCleaning == nil
}

func (uc *RoomUsecase) RemoveRoom(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("ID must be more than 0"))
	}
	occupied, err := uc.Repo.IsOccupied(ctx, id)
	if err != nil {
		return err
	}
	if occupied {
		return errors.Join(ErrValidation, errors.New("room is occupied"))
	}

	err = uc.Repo.DeleteRoom(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *RoomUsecase) GetList(ctx context.Context) ([]md.Room, error) {
	response, err := uc.Repo.ListRoom(ctx)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return response, errors.Join(ErrConflict, errors.New("database is clear"))
	}
	return response, nil
}

func (uc *RoomUsecase) GetFilteredRooms(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	response, err := uc.Repo.FilterRoom(ctx, filter)
	return response, err
}
