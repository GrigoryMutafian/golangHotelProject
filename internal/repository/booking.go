package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/model"
	"golangHotelProject/internal/repository/db"
	"log"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, b model.Booking) error
	GettingStatus(ctx context.Context, guest_id int) (bool, error)
	ArrivalStatusOfRoom(ctx context.Context, RoomID int) (bool, error)
	ReadBookingByID(ctx context.Context, id int) (model.Booking, error)
	PatchBooking(ctx context.Context, b dto.BookingPatch) error
	ListColumn(ctx context.Context) ([]model.Booking, error)
	FilterBookings(ctx context.Context, filter map[string]interface{}) (map[string][]int, error)
	DeleteBooking(ctx context.Context, id int) error
}

type PgBookingRepository struct {
	DB *sql.DB
}

func (r *PgBookingRepository) CreateBooking(ctx context.Context, b model.Booking) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO bookings (room_id, guest_id, start_date, end_date, status)
	VALUES($1, $2, $3, $4, $5)`, b.RoomID, b.GuestID, b.Start_date, b.End_date, b.Status)
	if err != nil {
		log.Printf("ERROR inserting booking: %v", err)
	}
	return err
}

func (r *PgBookingRepository) GettingStatus(ctx context.Context, guest_id int) (bool, error) {
	const q = `SELECT status FROM bookings WHERE guest_id = $1 AND status = true LIMIT 1`
	row := r.DB.QueryRowContext(ctx, q, guest_id)

	var isTrue bool
	err := row.Scan(&isTrue)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return isTrue, nil
}

func (r *PgBookingRepository) ArrivalStatusOfRoom(ctx context.Context, RoomID int) (bool, error) {
	const q = `SELECT status FROM bookings WHERE room_id = $1 AND status = true LIMIT 1`
	row := r.DB.QueryRowContext(ctx, q, RoomID)

	var isTrue bool
	err := row.Scan(&isTrue)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return isTrue, nil
}

func (r *PgBookingRepository) ReadBookingByID(ctx context.Context, id int) (model.Booking, error) {
	const q = `SELECT id, room_id, guest_id, start_date, end_date, status FROM bookings WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, q, id)

	var b model.Booking

	err := row.Scan(&b.ID, &b.RoomID, &b.GuestID, &b.Start_date, &b.End_date, &b.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Booking{}, err
		}
	}
	return b, nil
}

func (r *PgBookingRepository) PatchBooking(ctx context.Context, b dto.BookingPatch) error {
	const q = `UPDATE bookings SET room_id = $1, guest_id = $2, start_date = $3, end_date = $4, status = $5 WHERE id = $6`
	rows, err := r.DB.ExecContext(ctx, q, b.RoomID, b.GuestID, b.Start_date, b.End_date, b.Status, b.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows")
	}
	return nil
}

func (r *PgBookingRepository) ListColumn(ctx context.Context) ([]model.Booking, error) {

	rows, err := r.DB.QueryContext(ctx, `SELECT id, room_id, guest_id, start_date, end_date, status FROM bookings`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bookings []model.Booking

	for rows.Next() {
		var b model.Booking

		err := rows.Scan(&b.ID, &b.RoomID, &b.GuestID, &b.Start_date, &b.End_date, &b.Status)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, b)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *PgBookingRepository) FilterBookings(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	responses := make(map[string][]int)
	for column, value := range filter {
		query := fmt.Sprintf("SELECT id FROM bookings WHERE %s = $1", column)
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

func (r *PgBookingRepository) DeleteBooking(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM bookings WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
