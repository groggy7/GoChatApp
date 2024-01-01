package handler

import (
	"net/http"
	"server/pkg/dto"
	"server/pkg/service"
	"server/pkg/util"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(svc *service.UserService) UserHandler {
	return &userHandler{
		userService: *svc,
	}
}

func (h *userHandler) Signup(ctx *gin.Context) {
	var req dto.CreateUserRequest

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

	service.CreateSession(req.Username)
	ctx.JSON(http.StatusCreated, response)
}

func (h *userHandler) Login(ctx *gin.Context) {
	var req dto.GetUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = util.CheckHashAndPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong username or password"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Login Successful"})
}
