package middleware

import (
	"net/http"
	"server/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthenticateSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("SessionID")
		if err != nil {
			ctx.Next()
			return
		}

		if ctx.Writer.Written() {
			return
		}

		sessions := auth.GetSessions()

		for _, session := range sessions {
			if sessionID == session.SessionID {
				ctx.JSON(302, "Successful login")
				ctx.Redirect(http.StatusFound, "/home")
				return
			}
		}
		ctx.Next()
	}
}
