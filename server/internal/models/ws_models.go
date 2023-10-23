package models

import "github.com/gorilla/websocket"

type Room struct {
	RoomId      int          `json:"roomId"`
	Clients     []Client     `json:"clients"`
	MessageChan chan Message `json:"messageChan"`
}

type Message struct {
	Content string `json:"content"`
	Client  Client `json:"client"`
	Room    Room   `json:"room"`
}

type Client struct {
	Username   string          `json:"username"`
	Connection *websocket.Conn `json:"connection"`
}

type Hub struct {
	Rooms []Room
}
