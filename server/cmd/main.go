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

	cs := ws.NewChatServer()
	go cs.StartServer()
	go auth.InitSessionServer()

	router.StartRouter(userHandler, config)
}
