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

func NewChatHandler() ChatHandler {
	return ChatHandler{
		Hub: Hub{
			Rooms: []Room{},
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
		ch.Logger.Println("WebSocket upgrade failed:", err)
		return
	}

	username := ctx.Param("username")
	roomStr := ctx.Param("roomid")
	roomid, err := strconv.Atoi(roomStr)
	if err != nil {
		ch.Logger.Println("Invalid room ID:", err)
		return
	}

	cl := Client{
		Username:   username,
		Connection: conn,
	}

	var rm *Room
	for i, room := range ch.Hub.Rooms {
		if room.Id == int(roomid) {
			rm = &ch.Hub.Rooms[i]
			break
		}
	}

	if rm == nil {
		cl.Connection.WriteMessage(websocket.CloseMessage, []byte("No such room found"))
		cl.Connection.Close()
		return
	}

	rm.Clients = append(rm.Clients, cl)

	if len(rm.Clients) > 1 {
		message := Message{
			Content: username + " joined to the room",
		}
		rm.MessageChan <- message
	}

	go ch.WaitForMessage(cl.Connection, rm)
	go ch.ProcessMessages()
}

func (ch *ChatHandler) WaitForMessage(conn *websocket.Conn, rm *Room) {
	for {
		msg := Message{}
		if err := conn.ReadJSON(&msg); err != nil {
			ch.Logger.Println("Error reading from WebSocket:", err)
			break
		}
		rm.MessageChan <- msg
	}
}

func (ch *ChatHandler) ProcessMessages() {
	for {
		for _, room := range ch.Hub.Rooms {
			select {
			case message := <-room.MessageChan:
				log.Println(message)
			default:
				// No message available, continue to the next room
			}
		}
	}
}
