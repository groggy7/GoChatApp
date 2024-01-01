package main

import (
	"log"

	"server/internal/http/handler"
	"server/internal/http/server"

	"server/internal/ws"
	"server/pkg/db"
	"server/pkg/repository"
	"server/pkg/service"

	"github.com/gin-contrib/cors"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Println(err)
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5500"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}

	userRepo := repository.NewUserRepository(db.Conn)
	userService := service.NewUserService(&userRepo)
	userHandler := handler.NewUserHandler(&userService)

	chatHandler := ws.NewChatHandler()

	server.StartEngine(userHandler, *chatHandler, config)
}
