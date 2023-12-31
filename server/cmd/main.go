package main

import (
	"log"

	"server/internal/http/auth"
	"server/internal/http/handler"
	"server/internal/http/router"

	"server/internal/ws"
	"server/pkg/db"
	"server/pkg/repository"
	"server/pkg/service"

	"github.com/gin-contrib/cors"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Println(err)
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5500"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(&userRepo)
	userHandler := handler.NewUserHandler(&userService)

	chatHandler := ws.NewChatHandler()
	go chatHandler.WaitForMessages()

	go auth.InitSessionServer()

	router.StartRouter(userHandler, *chatHandler, config)
}
