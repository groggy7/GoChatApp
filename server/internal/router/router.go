package router

import (
	"log"
	"server/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartRouter(userHandler handlers.UserHandler, roomHandler handlers.RoomHandler, corsConfig cors.Config) {
	r := gin.Default()

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	r.Use(cors.New(corsConfig))
	r.GET("/createroom", roomHandler.CreateRoom)
	r.POST("/joinroom", roomHandler.JoinRoom)

	log.Println("Started http server at port 8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalln(err)
	}
}
