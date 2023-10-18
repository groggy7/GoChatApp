package handlers

import (
	"chatapp/server/internal/auth"
	"chatapp/server/internal/models"
	"chatapp/server/internal/services"
	"chatapp/server/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type handler struct {
	userService services.UserService
}

func NewUserHandler(svc *services.UserService) UserHandler {
	return &handler{
		userService: *svc,
	}
}

func (h *handler) Signup(ctx *gin.Context) {
	var req models.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := util.HashThePassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Password = hashedPassword

	response, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auth.CreateSession(req.Username)
	ctx.JSON(http.StatusCreated, response)
}

func (h *handler) Login(ctx *gin.Context) {
	var req models.GetUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = util.CheckHashAndPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong username or password"})
		return
	}

	session, err := auth.GetSession(user.Username)
	if err != nil {
		log.Println(err)
	}

	ctx.SetCookie("SessionID", session.SessionID, session.Expiry.Day(), "/", "localhost", true, true)

	ctx.JSON(200, gin.H{"message": "Login Successful"})
}
