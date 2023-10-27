package models

type JoinRoomRequest struct {
	Username string `json:"username"`
	RoomId   int    `json:"roomid"`
}

type CreateRoomResponse struct {
	RoomId int
}

type SendMessageRequest struct {
	Content  string `json:"content"`
	Username string `json:"client"`
	RoomId   int    `json:"roomid"`
}
