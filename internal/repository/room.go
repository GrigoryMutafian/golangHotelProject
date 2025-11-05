package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	md "golangHotelProject/internal/model"
	"golangHotelProject/internal/repository/db"
	"strconv"
	"strings"
)

type RoomRepository interface {
	Create(ctx context.Context, room md.Room) error
	List(ctx context.Context) ([]md.Room, error)
	Filter(ctx context.Context, filter map[string]interface{}) (map[string][]int, error)
	IsNumberExists(ctx context.Context, number int) (bool, error)
	Patch(ctx context.Context, id int, p dto.RoomPatch) error
	Delete(ctx context.Context, id int) error
	IsOccupied(ctx context.Context, roomID int) (bool, error)
}

type PgRoomRepository struct {
	DB *sql.DB
}

func (r *PgRoomRepository) Create(ctx context.Context, room md.Room) error { //addRoom
	_, err := r.DB.ExecContext(ctx, `INSERT INTO rooms (number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning)
	VALUES($1, $2, $3, $4, $5, $6, $7)`, room.Number, room.RoomCount, room.IsOccupied, room.Floor, room.SleepingPlaces, room.RoomType, room.NeedCleaning)

	return err
}

func (r *PgRoomRepository) IsNumberExists(ctx context.Context, number int) (bool, error) {
	const q = `SELECT 1 FROM rooms WHERE number = $1 LIMIT 1`
	row := r.DB.QueryRowContext(ctx, q, number)

	var inFindingRow int
	if err := row.Scan(&inFindingRow); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PgRoomRepository) List(ctx context.Context) ([]md.Room, error) {

	rows, err := r.DB.QueryContext(ctx, `SELECT id, number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning FROM rooms`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []md.Room

	for rows.Next() {
		var r md.Room

		err := rows.Scan(&r.ID, &r.Number, &r.RoomCount, &r.IsOccupied, &r.Floor, &r.SleepingPlaces, &r.RoomType, &r.NeedCleaning)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *PgRoomRepository) Filter(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	responses := make(map[string][]int)
	for column, value := range filter {
		query := fmt.Sprintf("SELECT id FROM rooms WHERE %s = $1", column)
		rows, err := db.DB.QueryContext(ctx, query, value)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var currentID int
			rows.Scan(&currentID)

			strValue := fmt.Sprintf("%v", value)
			responses[strValue] = append(responses[strValue], currentID)
		}
	}
	return responses, nil
}
func (r *PgRoomRepository) Patch(ctx context.Context, id int, p dto.RoomPatch) error {
	sets := make([]string, 0, 7)
	args := make([]any, 0, 7)

	next := func() string { return "$" + strconv.Itoa(len(args)+1) }

	if p.RoomCount != nil {
		sets = append(sets, "room_count = "+next())
		args = append(args, *p.RoomCount)
	}
	if p.IsOccupied != nil {
		sets = append(sets, "is_occupied = "+next())
		args = append(args, *p.IsOccupied)
	}
	if p.Floor != nil {
		sets = append(sets, "floor = "+next())
		args = append(args, *p.Floor)
	}
	if p.SleepingPlaces != nil {
		sets = append(sets, "sleeping_places = "+next())
		args = append(args, *p.SleepingPlaces)
	}
	if p.RoomType != nil {
		sets = append(sets, "room_type = "+next())
		args = append(args, *p.RoomType)
	}
	if p.NeedCleaning != nil {
		sets = append(sets, "need_cleaning = "+next())
		args = append(args, *p.NeedCleaning)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	q := "UPDATE rooms SET " + strings.Join(sets, ", ") + " WHERE id = $" + strconv.Itoa(len(args))

	res, err := r.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PgRoomRepository) Delete(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgRoomRepository) IsOccupied(ctx context.Context, roomID int) (bool, error) {
	const q = `SELECT is_occupied FROM rooms WHERE id = $1`
	var occupied bool
	if err := r.DB.QueryRowContext(ctx, q, roomID).Scan(&occupied); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sql.ErrNoRows
		}
		return false, err
	}
	return occupied, nil
}
