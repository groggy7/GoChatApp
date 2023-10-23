package ws

import (
	"log"
	"os"
	"server/internal/models"

	"github.com/gorilla/websocket"
)

type ChatServer struct {
	clients     []models.Client
	rooms       []models.Room
	messageChan chan models.Message
	upgrader    websocket.Upgrader
	logger      *log.Logger
}

func NewChatServer() ChatServer {
	return ChatServer{
		clients:     []models.Client{},
		rooms:       []models.Room{},
		messageChan: make(chan models.Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		logger: log.New(os.Stdout, "Chat Server - ", log.Lshortfile),
	}
}
