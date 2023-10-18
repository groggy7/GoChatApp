package main

import (
	"chatapp/server/db"
	"chatapp/server/internal/auth"
	"chatapp/server/internal/handlers"
	"chatapp/server/internal/repositories"
	"chatapp/server/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Println(err)
	}

	userRepository := repositories.NewUserRepository(db.GetDB())
	userService := services.NewUserService(&userRepository)
	userHandler := handlers.NewUserHandler(&userService)

	r := gin.Default()
	go auth.InitSessionServer()

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalln(err)
	}
}
