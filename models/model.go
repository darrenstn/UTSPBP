package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Room struct {
	ID       int    `json:"room_id" gorm:"primaryKey"`
	RoomName string `json:"room_name"`
}

type Rooms struct {
	Rooms []Room `json:"rooms"`
}

type RoomsResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Rooms  `json:"data"`
}

type Participant struct {
	ID        int    `json:"id"`
	AccountID int    `json:"id_account"`
	Username  string `json:"username"`
}

type RoomParticipants struct {
	ID           int           `json:"id"`
	RoomName     string        `json:"room_name"`
	Participants []Participant `json:"participants"`
}

type DetailRoomParticipants struct {
	RoomsParticipants RoomParticipants `json:"room"`
}

type RoomParticipantsResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    DetailRoomParticipants `json:"data"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
