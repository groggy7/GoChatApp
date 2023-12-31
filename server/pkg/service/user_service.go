package service

import (
	"context"
	"server/pkg/model"
	"server/pkg/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.CreateUserResponse, error)
	GetUserByEmail(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(r *repository.UserRepository) UserService {
	return &userService{
		userRepository: *r,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := model.CreateUserResponse{
		Id:       user.Id,
		Username: req.Username,
		Email:    req.Email,
	}

	return &resp, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	resp := &model.GetUserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	return resp, nil
}
