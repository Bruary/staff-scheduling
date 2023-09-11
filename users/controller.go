package users

import (
	"context"

	userModels "github.com/Bruary/staff-scheduling/users/models"
)

type ControllerInterface interface {
	CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse
	GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) *userModels.GetUserResponse
	DeleteUser(context.Context, userModels.DeleteUserRequest) *userModels.DeleteUserResponse
	UpdateUserPermissionLevel(ctx context.Context, req userModels.UpdateUserPermissionLevelRequest) *userModels.UpdateUserPermissionLevelResponse
	GetAllUsersWithShifts(ctx context.Context, req userModels.GetAllUsersRequest) *userModels.GetAllUsersResponse
}

type ControllerService struct {
	usersService ServiceInterface
}

var _ ControllerInterface = &ControllerService{}

func NewControllerService(usersService ServiceInterface) ControllerInterface {
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
// @Router /api/v1/signup [post]
func (s *ControllerService) CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse {
	return s.usersService.CreateUser(ctx, req)
}

func (s *ControllerService) GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) *userModels.GetUserResponse {
	return s.usersService.GetUserByEmail(ctx, req)
}

// @Title Delete User
// @Summary Delete user
// @ID delete-user
// @Produce json
// @Param req body userModels.DeleteUserRequest true "delete user request"
// @Success 200 {object} userModels.DeleteUserResponse
// @Router /api/v1/user [delete]
func (s *ControllerService) DeleteUser(ctx context.Context, req userModels.DeleteUserRequest) *userModels.DeleteUserResponse {
	return s.usersService.DeleteUser(ctx, req)
}

// @Title Update User Permission Level
// @Summary Update user permission level
// @ID update-user-permission-level
// @Produce json
// @Param req body userModels.UpdateUserPermissionLevelRequest true "update user permission request"
// @Success 200 {object} userModels.UpdateUserPermissionLevelResponse
// @Router /api/v1/user/permission [put]
func (s *ControllerService) UpdateUserPermissionLevel(ctx context.Context, req userModels.UpdateUserPermissionLevelRequest) *userModels.UpdateUserPermissionLevelResponse {
	return s.usersService.UpdateUserPermissionLevel(ctx, req)
}

// @Title Get All Users With Shifts
// @Summary Get all users with accumulated work hours
// @ID get-all-users-with-shifts
// @Produce json
// @Param req body userModels.GetAllUsersRequest true "get all users with shifts request"
// @Success 200 {object} userModels.GetAllUsersResponse
// @Router /api/v1/users/shifts [get]
func (s *ControllerService) GetAllUsersWithShifts(ctx context.Context, req userModels.GetAllUsersRequest) *userModels.GetAllUsersResponse {
	return s.usersService.GetAllUsersWithShifts(ctx, req)
}
