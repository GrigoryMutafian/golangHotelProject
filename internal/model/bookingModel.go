package model

import "time"

type Booking struct {
	ID         int       `json:"id"`
	Guest_ID   int       `json:"guest_id"`
	Start_date time.Time `json:"start_date"`
	End_date   time.Time `json:"end_date"`
	Status     string    `json:"status"`
}
