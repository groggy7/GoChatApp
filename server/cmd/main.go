package main

import (
	"log"
	"server/db"
	"server/internal/auth"
	"server/internal/handlers"
	"server/internal/repositories"
	"server/internal/router"
	"server/internal/services"
	"server/internal/ws"

	"github.com/gin-contrib/cors"
)

var chatServer ws.ChatServer

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Println(err)
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5500"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}

	userRepo := repositories.NewUserRepository(db.GetDB())
	userService := services.NewUserService(&userRepo)
	userHandler := handlers.NewUserHandler(&userService)

	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(&roomRepo)
	roomHandler := handlers.NewRoomHandler(&roomService)

	router.StartRouter(userHandler, roomHandler, config)
	chatServer = ws.NewChatServer()

	go auth.InitSessionServer()

}
