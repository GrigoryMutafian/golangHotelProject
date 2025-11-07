package repository

import (
	"context"
	"database/sql"
	md "golangHotelProject/internal/model"
	"log"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, b md.Booking) error
	GettingStatus(ctx context.Context, guest_id int) (bool, error)
	ArrivalStatusOfRoom(ctx context.Context, RoomID int) (bool, error)
}

type PgBookingRepository struct {
	DB *sql.DB
}

func (r *PgBookingRepository) CreateBooking(ctx context.Context, b md.Booking) error {
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
