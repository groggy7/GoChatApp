package repositories

import (
	"fmt"
	"server/internal/models"
)

type RoomRepository interface {
	CreateRoom() (models.Room, error)
	JoinRoom(client *models.Client, roomId int) error
}

type roomRepository struct {
	rooms map[int]models.Room
}

func NewRoomRepository() RoomRepository {
	return &roomRepository{
		rooms: make(map[int]models.Room),
	}
}

func (r *roomRepository) CreateRoom() (models.Room, error) {
	newRoom := models.Room{
		RoomId:      len(r.rooms) + 1,
		Clients:     nil,
		MessageChan: make(chan models.Message),
	}

	r.rooms[newRoom.RoomId] = newRoom
	return newRoom, nil
}

func (r *roomRepository) JoinRoom(client *models.Client, roomId int) error {
	room, ok := r.rooms[roomId]
	if !ok {
		return fmt.Errorf("failed to join room: no room found with the ID %d", roomId)
	}

	room.Clients = append(room.Clients, *client)
	return nil
}
