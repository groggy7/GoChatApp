package user

import (
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	GetUserByEmail(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(r *UserRepository) UserService {
	return &userService{
		userRepository: *r,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := CreateUserResponse{
		Id:       user.Id,
		Username: req.Username,
		Email:    req.Email,
	}

	return &resp, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	resp := &GetUserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	return resp, nil
}
