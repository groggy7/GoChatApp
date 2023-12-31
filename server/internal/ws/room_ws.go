package ws

import (
	"log"
	"net/http"
	"os"
	"strconv"

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
	conn, err := ch.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := ctx.Query("username")
	roomStr := ctx.Query("room")
	roomId, err := strconv.Atoi(roomStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := Client{
		Username:   username,
		Connection: conn,
	}

	matchedRoom := &Room{}
	for i := range ch.Hub.Rooms {
		if ch.Hub.Rooms[i].Id == roomId {
			matchedRoom = &ch.Hub.Rooms[i]
			break
		}
	}
	if matchedRoom.Id == 0 {
		ch.Logger.Println("No matching room found")
		return
	}

	matchedRoom.Clients = append(matchedRoom.Clients, client)
	ctx.JSON(http.StatusOK, gin.H{"client": client.Username, "room": roomId})

	matchedRoom.MessageChan <- Message{
		Content: "User " + client.Username + " has joined the room",
		Client:  client,
		Room:    *matchedRoom,
	}

	for {
		message := Message{}
		if err := conn.ReadJSON(&message); err != nil {
			ch.Logger.Println("Error reading json.", err)
			break
		}

		ch.Logger.Println("Message received: ", message)
		matchedRoom.MessageChan <- message
	}
}

func (ch *ChatHandler) WaitForMessages() {
	ch.Logger.Println("WaitForMessages function started")
	for _, room := range ch.Hub.Rooms {
		go func(room Room) {
			ch.Logger.Printf("Waiting for messages in room %s\n", room.Name)

			for message := range room.MessageChan {
				ch.Logger.Println("Message received: ", message)

				for _, client := range room.Clients {
					if err := client.Connection.WriteJSON(message); err != nil {
						ch.Logger.Println("Error writing json.", err)
					}
				}
			}
		}(room)
	}
}
