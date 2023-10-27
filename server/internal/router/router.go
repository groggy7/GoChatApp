package router

import (
	"log"
	"server/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartRouter(userHandler handlers.UserHandler, corsConfig cors.Config) {
	r := gin.Default()

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	r.Use(cors.New(corsConfig))

	log.Println("Started http server at port 8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalln(err)
	}
}
