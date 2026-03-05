package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/logger"
	md "golangHotelProject/internal/model"
	repo "golangHotelProject/internal/repository"
	"log/slog"
)

var (
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict error")
)

func IsValidationErr(err error) bool { return errors.Is(err, ErrValidation) }
func IsConflictErr(err error) bool   { return errors.Is(err, ErrConflict) }

type RoomUsecase struct {
	Repo   repo.RoomRepository
	Logger *slog.Logger
}

func NewRoomUsecase(repo repo.RoomRepository, log logger.Logger) *RoomUsecase {
	return &RoomUsecase{
		Repo:   repo,
		Logger: log.With("component", "RoomUsecase"),
	}
}

func (uc *RoomUsecase) AddRoom(ctx context.Context, room md.Room) error {
	const op = "AddRoom"

	uc.Logger.Debug("adding new room",
		"op", op,
		"room_number", room.Number,
		"room_type", room.RoomType,
	)

	err := validateRoom(room)
	if err != nil {
		uc.Logger.Warn("room validation failed",
			"op", op,
			"room_number", room.Number,
			"error", err.Error(),
		)
		return errors.Join(ErrValidation, err)
	}

	exists, err := uc.Repo.IsNumberExists(ctx, room.Number)
	if err != nil {
		uc.Logger.Error("failed to check room existence",
			"op", op,
			"room_number", room.Number,
			"error", err.Error(),
		)
		return err
	}
	if exists {
		uc.Logger.Warn("room already exists",
			"op", op,
			"room_number", room.Number,
		)
		return errors.Join(ErrConflict, errors.New("room number already exists"))
	}

	if err := uc.Repo.CreateRoom(ctx, room); err != nil {
		uc.Logger.Error("failed to create room",
			"op", op,
			"room_number", room.Number,
			"error", err.Error(),
		)
		return err
	}

	uc.Logger.Info("room created successfully",
		"op", op,
		"room_number", room.Number,
	)
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
	const op = "PatchRoom"

	uc.Logger.Debug("patching room",
		"op", op,
		"room_id", id,
	)

	if id <= 0 {
		uc.Logger.Warn("invalid room id",
			"op", op,
			"room_id", id,
		)
		return errors.Join(ErrValidation, errors.New("invalid id"))
	}

	if isEmptyPatch(p) {
		uc.Logger.Debug("empty patch, skipping update",
			"op", op,
			"room_id", id,
		)
		return nil
	}

	if p.Floor != nil && *p.Floor <= 0 {
		uc.Logger.Warn("invalid floor value",
			"op", op,
			"room_id", id,
			"floor", *p.Floor,
		)
		return errors.Join(ErrValidation, errors.New("floor must be more then 0"))
	}

	if p.SleepingPlaces != nil && *p.SleepingPlaces <= 0 {
		uc.Logger.Warn("invalid sleeping places value",
			"op", op,
			"room_id", id,
			"sleeping_places", *p.SleepingPlaces,
		)
		return errors.Join(ErrValidation, errors.New("sleepng places must be more then 0"))
	}
	if p.RoomType != nil {
		switch *p.RoomType {
		case "Standard", "Deluxe", "Suite":
		default:
			uc.Logger.Warn("invalid room type",
				"op", op,
				"room_id", id,
				"room_type", *p.RoomType,
			)
			return errors.Join(ErrValidation, errors.New("room_type must be one of: Standard, Deluxe, Suite"))
		}
	}
	if p.RoomCount != nil && *p.RoomCount <= 0 {
		uc.Logger.Warn("invalid room count value",
			"op", op,
			"room_id", id,
			"room_count", *p.RoomCount,
		)
		return errors.Join(ErrValidation, errors.New("room count must be more then 0"))

	}
	if p.IsOccupied != nil && p.NeedCleaning != nil {
		if *p.IsOccupied && *p.NeedCleaning {
			uc.Logger.Warn("cannot set need_cleaning while room is occupied",
				"op", op,
				"room_id", id,
			)
			return errors.Join(ErrValidation, errors.New("cannot set need_cleaning while room is occupied"))
		}
	}

	if err := uc.Repo.PatchRoom(ctx, id, p); err != nil {
		uc.Logger.Error("failed to patch room",
			"op", op,
			"room_id", id,
			"error", err.Error(),
		)
		return err
	}

	uc.Logger.Info("room patched successfully",
		"op", op,
		"room_id", id,
	)
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
	const op = "RemoveRoom"

	uc.Logger.Debug("removing room",
		"op", op,
		"room_id", id,
	)

	if id <= 0 {
		uc.Logger.Warn("invalid room id",
			"op", op,
			"room_id", id,
		)
		return errors.Join(ErrValidation, errors.New("ID must be more than 0"))
	}
	occupied, err := uc.Repo.IsOccupied(ctx, id)
	if err != nil {
		uc.Logger.Error("failed to check room occupation status",
			"op", op,
			"room_id", id,
			"error", err.Error(),
		)
		return err
	}
	if occupied {
		uc.Logger.Warn("cannot remove occupied room",
			"op", op,
			"room_id", id,
		)
		return errors.Join(ErrValidation, errors.New("room is occupied"))
	}

	err = uc.Repo.DeleteRoom(ctx, id)
	if err != nil {
		uc.Logger.Error("failed to delete room",
			"op", op,
			"room_id", id,
			"error", err.Error(),
		)
		return err
	}

	uc.Logger.Info("room removed successfully",
		"op", op,
		"room_id", id,
	)
	return nil
}

func (uc *RoomUsecase) GetList(ctx context.Context) ([]md.Room, error) {
	const op = "GetList"

	uc.Logger.Debug("fetching room list", "op", op)

	response, err := uc.Repo.ListRoom(ctx)
	if err != nil {
		uc.Logger.Error("failed to fetch room list",
			"op", op,
			"error", err.Error(),
		)
		return nil, err
	}
	if len(response) == 0 {
		uc.Logger.Info("room list is empty", "op", op)
		return response, errors.Join(ErrConflict, errors.New("database is clear"))
	}

	uc.Logger.Debug("room list fetched successfully",
		"op", op,
		"count", len(response),
	)
	return response, nil
}

func (uc *RoomUsecase) GetFilteredRooms(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	const op = "GetFilteredRooms"

	uc.Logger.Debug("filtering rooms",
		"op", op,
		"filter", filter,
	)

	response, err := uc.Repo.FilterRoom(ctx, filter)
	if err != nil {
		uc.Logger.Error("failed to filter rooms",
			"op", op,
			"filter", filter,
			"error", err.Error(),
		)
		return nil, err
	}

	uc.Logger.Debug("rooms filtered successfully",
		"op", op,
		"filter", filter,
		"results_count", len(response),
	)
	return response, err
}
