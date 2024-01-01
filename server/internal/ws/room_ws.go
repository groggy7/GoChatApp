package ws

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	Hub      Hub
	Upgrader websocket.Upgrader
	Logger   *log.Logger
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		Hub: Hub{
			Rooms:   []Room{},
			Clients: []Client{},
		},
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1000,
			WriteBufferSize: 1000,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Logger: log.New(os.Stdout, "Chat Server - ", log.Lshortfile),
	}
}

func (ch *ChatHandler) CreateRoom(ctx *gin.Context) {
	room := Room{
		Id:          len(ch.Hub.Rooms) + 1,
		Clients:     []Client{},
		MessageChan: make(chan Message),
	}

	ch.Hub.Rooms = append(ch.Hub.Rooms, room)
	ctx.JSON(http.StatusCreated, gin.H{"room id": room.Id})
}

func (ch *ChatHandler) JoinRoom(ctx *gin.Context) {

}
