package users

import (
	"context"

	userModels "github.com/Bruary/staff-scheduling/users/models"
)

type ControllerInterface interface {
	CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse
	GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) *userModels.GetUserResponse
	DeleteUser(context.Context, userModels.DeleteUserRequest) *userModels.DeleteUserResponse
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

// @Title Create User
// @Summary Create a new user
// @ID create-a-new-user
// @Produce json
// @Param req body userModels.CreateUserRequest true "create user request"
// @Success 200 {object} userModels.CreateUserResponse
// @Router /user [post]
func (s *ControllerService) CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse {
	return s.usersService.CreateUser(ctx, req)
}

func (s *ControllerService) GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) *userModels.GetUserResponse {
	return s.usersService.GetUserByEmail(ctx, req)
}

func (s *ControllerService) DeleteUser(ctx context.Context, req userModels.DeleteUserRequest) *userModels.DeleteUserResponse {
	return s.usersService.DeleteUser(ctx, req)
}
