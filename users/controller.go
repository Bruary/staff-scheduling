package users

import (
	"context"

	userModels "github.com/Bruary/staff-scheduling/users/models"
)

type ControllerInterface interface {
	CreateUser(ctx context.Context, req userModels.CreateUserRequest) (*userModels.CreateUserResponse, error)
	GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) (*userModels.GetUserResponse, error)
	DeleteUser(context.Context, userModels.DeleteUserRequest) (*userModels.DeleteUserResponse, error)
}

type ControllerService struct {
	usersService ServiceInterface
}

var _ ControllerInterface = &ControllerService{}

func NewControllerService(usersService ServiceInterface) *ControllerService {
	return &ControllerService{
		usersService: usersService,
	}
}

func (s *ControllerService) CreateUser(ctx context.Context, req userModels.CreateUserRequest) (*userModels.CreateUserResponse, error) {

	resp, _ := s.usersService.CreateUser(ctx, req)

	response := &userModels.CreateUserResponse{
		BaseResponse: resp.BaseResponse,
		User:         resp.User,
	}

	return response, nil
}

func (s *ControllerService) GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) (*userModels.GetUserResponse, error) {

	resp, _ := s.usersService.GetUserByEmail(ctx, req)

	response := &userModels.GetUserResponse{
		BaseResponse: resp.BaseResponse,
		User:         resp.User,
	}

	return response, nil
}

func (s *ControllerService) DeleteUser(ctx context.Context, req userModels.DeleteUserRequest) (*userModels.DeleteUserResponse, error) {

	resp, _ := s.usersService.DeleteUser(ctx, req)

	response := &userModels.DeleteUserResponse{
		BaseResponse: resp.BaseResponse,
		User:         resp.User,
	}

	return response, nil
}
