package handlers

import (
	"net/http"
	"server/models"
	"server/services"
	"server/util"

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
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	hashedPassword, err := util.HashThePassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	req.Password = hashedPassword

	response, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *handler) Login(ctx *gin.Context) {
	var req models.GetUserRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	err = util.CheckHashAndPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": "wrong username or password"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"message": "Login Successful"})
}
