package services

import (
	"chatapp/server/internal/models"
	"chatapp/server/internal/repositories"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.CreateUserResponse, error)
	GetUserByEmail(ctx context.Context, req *models.GetUserRequest) (*models.GetUserResponse, error)
}

type service struct {
	userRepository repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) UserService {
	return &service{
		userRepository: *r,
	}
}

func (s *service) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := models.CreateUserResponse{
		Id:       user.Id,
		Username: req.Username,
		Email:    req.Email,
	}

	return &resp, nil
}

func (s *service) GetUserByEmail(ctx context.Context, req *models.GetUserRequest) (*models.GetUserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	resp := &models.GetUserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	return resp, nil
}
