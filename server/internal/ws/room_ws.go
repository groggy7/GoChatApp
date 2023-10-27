package ws

import (
	"log"
	"net/http"
	"os"
	"server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatServer struct {
	Hub      models.Hub
	Upgrader websocket.Upgrader
	Logger   *log.Logger
}

func NewChatServer() ChatServer {
	return ChatServer{
		Hub: models.Hub{
			Rooms:   []models.Room{},
			Clients: []models.Client{},
		},
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Logger: log.New(os.Stdout, "Chat Server - ", log.Lshortfile),
	}
}

func (cs *ChatServer) StartServer() {
	go cs.handleIncomingMessages()

	r := gin.Default()

	r.GET("/ws/:user", cs.handleWebSocket)
	r.GET("/ws/createroom", cs.CreateRoom)
	r.POST("/ws/joinroom", cs.JoinRoom)
	r.POST("/ws/sendmessage", cs.SendMessage)

	log.Println("Started websocket server at port 4444")
	if err := r.Run("0.0.0.0:4444"); err != nil {
		log.Fatalln(err)
	}
}

func (cs *ChatServer) handleIncomingMessages() {
	for _, room := range cs.Hub.Rooms {
		go func(room models.Room) {
			message := <-room.MessageChan

			for _, client := range room.Clients {
				err := client.Connection.WriteJSON(message)
				if err != nil {
					log.Println(err)
				}
			}
		}(room)
	}
}

func (cs *ChatServer) handleWebSocket(ctx *gin.Context) {
	username := ctx.Param("user")

	conn, err := cs.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		cs.Logger.Println(err)
		return
	}

	newClient := models.Client{
		Username:   username,
		Connection: conn,
	}

	cs.Hub.Clients = append(cs.Hub.Clients, newClient)
}

func (cs *ChatServer) CreateRoom(ctx *gin.Context) {
	newRoom := models.Room{
		Id:          len(cs.Hub.Rooms) + 1,
		Clients:     nil,
		MessageChan: make(chan models.Message),
	}

	cs.Hub.Rooms = append(cs.Hub.Rooms, newRoom)
}

func (cs *ChatServer) JoinRoom(ctx *gin.Context) {
	var req models.JoinRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var matchedClient *models.Client
	for _, client := range cs.Hub.Clients {
		if client.Username == req.Username {
			matchedClient = &client
		}
	}
	if matchedClient == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "client not found"})
		return
	}

	var matchedRoom *models.Room
	for _, room := range cs.Hub.Rooms {
		if room.Id == req.RoomId {
			matchedRoom = &room
		}
	}
	if matchedRoom == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room not found"})
		return
	}

	matchedRoom.Clients = append(matchedRoom.Clients, *matchedClient)
}

func (cs *ChatServer) SendMessage(ctx *gin.Context) {
	var req models.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var matchedClient *models.Client
	for _, client := range cs.Hub.Clients {
		if client.Username == req.Username {
			matchedClient = &client
		}
	}
	if matchedClient == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "client not found"})
		return
	}

	var matchedRoom *models.Room
	for _, room := range cs.Hub.Rooms {
		if room.Id == req.RoomId {
			matchedRoom = &room
		}
	}
	if matchedRoom == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room not found"})
		return
	}

	message := models.Message{
		Content: req.Content,
		Client:  *matchedClient,
		Room:    *matchedRoom,
	}

	matchedRoom.MessageChan <- message
}
