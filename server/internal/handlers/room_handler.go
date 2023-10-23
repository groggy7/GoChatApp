package handlers

import (
	"net/http"
	"server/internal/models"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type RoomHandler interface {
	CreateRoom(ctx *gin.Context)
	JoinRoom(ctx *gin.Context)
}

type roomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(s *services.RoomService) RoomHandler {
	return &roomHandler{
		roomService: *s,
	}
}

func (rh *roomHandler) CreateRoom(ctx *gin.Context) {
	resp, err := rh.roomService.CreateRoom()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (rh *roomHandler) JoinRoom(ctx *gin.Context) {
	var req models.JoinRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
