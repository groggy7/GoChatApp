package main

import (
	"log"
	"server/internal/auth"
	"server/internal/auth/middleware"
	"server/internal/db"
	"server/internal/handlers"
	"server/internal/repositories"
	"server/internal/services"

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
	r.Use(middleware.AuthenticateSession())

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)
	r.GET("/home", userHandler.Homepage)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalln(err)
	}
}
