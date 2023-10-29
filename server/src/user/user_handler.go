package user

import (
	"log"
	"net/http"
	"server/src/auth"
	"server/src/util"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(svc *UserService) UserHandler {
	return &userHandler{
		userService: *svc,
	}
}

func (h *userHandler) Signup(ctx *gin.Context) {
	var req CreateUserRequest

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

func (h *userHandler) Login(ctx *gin.Context) {
	var req GetUserRequest

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
