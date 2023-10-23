package services

import (
	"server/internal/models"
	"server/internal/repositories"
)

type RoomService interface {
	CreateRoom() (models.CreateRoomResponse, error)
	JoinRoom(req *models.JoinRoomRequest) error
}

type roomService struct {
	roomRepository repositories.RoomRepository
}

func NewRoomService(r *repositories.RoomRepository) RoomService {
	return &roomService{
		roomRepository: *r,
	}
}

func (rs *roomService) CreateRoom() (models.CreateRoomResponse, error) {
	room, err := rs.roomRepository.CreateRoom()
	resp := models.CreateRoomResponse{
		RoomId: room.RoomId,
	}

	return resp, err
}

func (rs *roomService) JoinRoom(req *models.JoinRoomRequest) error {
	return rs.roomRepository.JoinRoom(&models.Client{Username: req.Username, Connection: req.Connection}, req.RoomId)
}
