package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"server/pkg/model"
	"server/pkg/service"
	"strings"

	"github.com/gin-gonic/gin"
)

var client = service.GetClient()

func GetSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.Path, "/login") {
			username := ctx.Param("username")
			if username == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no username"})
				ctx.Abort()
				return
			}

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
		}
		ctx.Next()
	}
}
