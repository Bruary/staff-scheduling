package models

import (
	"github.com/Bruary/staff-scheduling/models"
)

type User struct {
	Id        int32  `json:"id,omitempty"`
	Created   string `json:"created,omitempty"`
	Type      string `json:"type,omitempty"`
	Uid       string `json:"uid,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Updated   string `json:"updated,omitempty"`
	Deleted   string `json:"deleted,omitempty"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type CreateUserResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	User         User                 `json:"user,omitempty"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email,omitempty"`
}

type GetUserResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	User         User                 `json:"user,omitempty"`
}

type DeleteUserRequest struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type DeleteUserResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	User         User                 `json:"user,omitempty"`
}

type UpdateUserPermissionLevelRequest struct {
	Email           string                 `json:"email"`
	PermissionLevel models.PermissionLevel `json:"permission_level"`
}

type UpdateUserPermissionLevelResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	User         User                 `json:"user,omitempty"`
}
