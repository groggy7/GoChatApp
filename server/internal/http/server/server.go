package server

import (
	"log"

	"server/internal/http/handler"
	"server/internal/http/middleware"
	"server/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartEngine(userHandler handler.UserHandler, chatHandler ws.ChatHandler, corsConfig cors.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(corsConfig))
	r.Use(middleware.GetSessionMiddleware())

	r.POST("/signup", userHandler.Signup)
	r.POST("/login/:username", userHandler.Login)

	r.GET("/ws/create", chatHandler.CreateRoom)
	r.GET("/ws/join", chatHandler.JoinRoom)

	return r
}

func StartServer(engine *gin.Engine) error {
	log.Println("Started http server at port 8080")
	if err := engine.Run("0.0.0.0:8080"); err != nil {
		return err
	}

	return nil
}
