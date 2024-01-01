package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"server/pkg/model"
	"server/pkg/service"

	"github.com/gin-gonic/gin"
)

var client = service.GetClient()

func GetSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		log.Println("in session middleware")
		log.Println(username)
		sessionJSON, err := client.Get(context.Background(), username).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		var session model.Session
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("session", &session)
		ctx.Next()
	}
}
