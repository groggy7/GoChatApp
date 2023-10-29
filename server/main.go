package main

import (
	"log"
	"server/db"
	"server/src/auth"
	"server/src/router"
	"server/src/user"
	"server/src/ws"

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

	userRepo := user.NewUserRepository(db.GetDB())
	userService := user.NewUserService(&userRepo)
	userHandler := user.NewUserHandler(&userService)

	cs := ws.NewChatServer()
	go cs.StartServer()
	go auth.InitSessionServer()

	router.StartRouter(userHandler, config)
}
