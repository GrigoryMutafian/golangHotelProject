package dto

type RoomPatch struct {
	RoomCount      *int    `json:"room_count,omitempty"`
	IsOccupied     *bool   `json:"is_occupied,omitempty"`
	Floor          *int    `json:"floor,omitempty"`
	SleepingPlaces *int    `json:"sleeping_places,omitempty"`
	RoomType       *string `json:"room_type,omitempty"`
	NeedCleaning   *bool   `json:"need_cleaning,omitempty"`
}
