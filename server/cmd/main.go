package main

import (
	"log"

	"server/internal/http/handler"
	"server/internal/http/server"

	"server/internal/ws"
	"server/pkg/db"
	"server/pkg/repository"
	"server/pkg/service"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Println(err)
	}

	userRepo := repository.NewUserRepository(db.Conn)
	userService := service.NewUserService(&userRepo)
	userHandler := handler.NewUserHandler(&userService)

	chatHandler := ws.NewChatHandler()

	engine := server.StartEngine(userHandler, *chatHandler)

	if err := server.StartServer(engine); err != nil {
		log.Fatalln(err)
	}
}
