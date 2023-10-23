package models

import "github.com/gorilla/websocket"

type JoinRoomRequest struct {
	Username   string          `json:"username"`
	RoomId     int             `json:"roomid"`
	Connection *websocket.Conn `json:"connection"`
}

type CreateRoomResponse struct {
	RoomId int
}
