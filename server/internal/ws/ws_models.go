package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms []Room
}

type Room struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Clients     []Client     `json:"clients"`
	MessageChan chan Message `json:"messageChan"`
}

type Client struct {
	Username   string
	Connection *websocket.Conn
}

type Message struct {
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"-"`
}
