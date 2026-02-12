package dto

import "time"

type RoomDTO struct {
	ID             int     `json:"id" example:"1"`
	RoomType       string  `json:"roomType" example:"standard"`
	Price          float64 `json:"price" example:"150.00"`
	IsAvailable    bool    `json:"isAvailable" example:"true"`
	RoomCount      int     `json:"roomCount" example:"1"`
	IsOccupied     bool    `json:"isOccupied" example:"false"`
	Floor          int     `json:"floor" example:"1"`
	SleepingPlaces int     `json:"sleepingPlaces" example:"2"`
	NeedCleaning   bool    `json:"needCleaning" example:"false"`
}

type RoomPatch struct {
	RoomCount      *int    `json:"roomCount,omitempty" example:"2"`
	IsOccupied     *bool   `json:"isOccupied,omitempty" example:"true"`
	Floor          *int    `json:"floor,omitempty" example:"2"`
	SleepingPlaces *int    `json:"sleepingPlaces,omitempty" example:"3"`
	RoomType       *string `json:"roomType,omitempty" example:"luxury"`
	NeedCleaning   *bool   `json:"needCleaning,omitempty" example:"true"`
}

type BookingDTO struct {
	ID         int       `json:"id" example:"1"`
	RoomID     int       `json:"roomId" example:"1"`
	GuestID    int       `json:"guestId" example:"1"`
	Start_date time.Time `json:"startDate" example:"2024-01-15T00:00:00Z"`
	End_date   time.Time `json:"endDate" example:"2024-01-20T00:00:00Z"`
	Status     string    `json:"status" example:"confirmed"`
}

type BookingPatch struct {
	ID         *int       `json:"id" example:"1"`
	RoomID     *int       `json:"roomId,omitempty" example:"2"`
	GuestID    *int       `json:"guestId,omitempty" example:"2"`
	Start_date *time.Time `json:"startDate,omitempty" example:"2024-01-16T00:00:00Z"`
	End_date   *time.Time `json:"endDate,omitempty" example:"2024-01-21T00:00:00Z"`
	Status     *string    `json:"status,omitempty" example:"cancelled"`
}

type CreatingRoomResponse struct {
	Message string `json:"message"`
	RoomID  int    `json:"roomId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RoomPatchResponse struct {
	Status string `json:"status"`
}

type RemoveRoomRequest struct {
	RoomID int `json:"roomId"`
}
