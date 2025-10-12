package hotelModel

type Room struct {
	ID             int    `json:"id"`
	Number         int    `json:"number"`
	RoomCount      int    `json:"room_count"`
	IsOccupied     bool   `json:"is_occupied"`
	Floor          int    `json:"floor"`
	SleepingPlaces int    `json:"sleeping_places"`
	RoomQuality    string `json:"room_quality"`
	NeedCleaning   bool   `json:"need_cleaning"`
}
