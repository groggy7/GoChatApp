package ws

import "github.com/gorilla/websocket"

type Hub struct {
	Rooms   []Room
	Clients []Client
}

type Room struct {
	Id          int          `json:"idd"`
	Clients     []Client     `json:"clients"`
	MessageChan chan Message `json:"messageChan"`
}

type Client struct {
	Username   string
	Connection *websocket.Conn
}

type Message struct {
	Content string `json:"content"`
	Client  Client `json:"client"`
	Room    Room   `json:"room"`
}
