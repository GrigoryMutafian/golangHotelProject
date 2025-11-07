package dto

type RoomPatch struct {
	RoomCount      *int    `json:"room_count,omitempty"`
	IsOccupied     *bool   `json:"is_occupied,omitempty"`
	Floor          *int    `json:"floor,omitempty"`
	SleepingPlaces *int    `json:"sleeping_places,omitempty"`
	RoomType       *string `json:"room_type,omitempty"`
	NeedCleaning   *bool   `json:"need_cleaning,omitempty"`
}

type BookingDTO struct {
	ID         int    `json:"id,omitempty"`
	RoomID     int    `json:"room_id"`
	GuestID    int    `json:"guest_id"`
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
	Status     bool   `json:"status"`
}
