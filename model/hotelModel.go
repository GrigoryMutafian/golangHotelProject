package model

type Room struct {
	ID             int    `json:"id"`
	Number         int    `json:"number"`
	RoomCount      int    `json:"room_count"`
	IsOccupied     bool   `json:"is_occupied"`
	Floor          int    `json:"floor"`
	SleepingPlaces int    `json:"sleeping_places"`
	RoomType       string `json:"room_type"`
	NeedCleaning   bool   `json:"need_cleaning"`
}
