package users

import (
	"context"
	"database/sql"
	"os"

	"github.com/Bruary/staff-scheduling/core/models"
	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	userRepo "github.com/Bruary/staff-scheduling/db/users"
	userModels "github.com/Bruary/staff-scheduling/users/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	UidForTests      string = "1426ec8a-b7bb-47a8-bbba-ab5553b90c70"
	PasswordForTests string = "helloWorld"
)

type ServiceInterface interface {
	CreateUser(context.Context, userModels.CreateUserRequest) *userModels.CreateUserResponse
	GetUserByEmail(context.Context, userModels.GetUserByEmailRequest) *userModels.GetUserResponse
	GetUserByUID(ctx context.Context, userUID string) *userModels.GetUserResponse
	DeleteUser(context.Context, userModels.DeleteUserRequest) *userModels.DeleteUserResponse
	UpdateUserPermissionLevel(context.Context, userModels.UpdateUserPermissionLevelRequest) *userModels.UpdateUserPermissionLevelResponse
}

type Service struct {
	UserRepo userRepo.UsersRepoInterface
}

var _ ServiceInterface = &Service{}

func New(userRepo userRepo.UsersRepoInterface) *Service {
	return &Service{
		UserRepo: userRepo,
	}
}

var _ ServiceInterface = &Service{}

func (s *Service) CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse {

	if req.Email == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params",
				ErrorMsg:  "Users email is missing",
			},
		}
	}

	if req.FirstName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params",
				ErrorMsg:  "Users first name is missing",
			},
		}
	}

	if req.LastName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params",
				ErrorMsg:  "Users last name is missing",
			},
		}
	}

	if req.Password == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params",
				ErrorMsg:  "Password is required",
			},
		}
	}

	resp, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Unknown Error",
				ErrorMsg:  err.Error(),
			},
		}
	} else if (err == nil && resp != sqlc.User{}) {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "User already exists",
				ErrorMsg:  "User is already registered",
			},
		}
	}

	// hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Password hashing failed",
				ErrorMsg:  err.Error(),
			},
		}
	}

	params := sqlc.CreateUserParams{
		Uid:       uuid.NewString(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	// to unify the uid/password used for testing
	if os.Getenv("environment") == "testing" {
		params.Uid = UidForTests
		params.Password = PasswordForTests
	}

	user, err := s.UserRepo.CreateUser(ctx, params)
	if err != nil {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Unknown error",
				ErrorMsg:  err.Error(),
			},
		}
	}

	return &userModels.CreateUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}
}

func (s *Service) GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) *userModels.GetUserResponse {

	var (
		user sqlc.User
		err  error
	)

	if req.Email == "" {
		return &userModels.GetUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing Parameters",
				ErrorMsg:  "Missing parameters, email and phone number could not be found.",
			},
		}
	}

	user, err = s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "User does not exist.",
					ErrorMsg:  err.Error(),
				},
			}

		} else {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "Unknown error.",
					ErrorMsg:  err.Error(),
				},
			}
		}
	}

	return &userModels.GetUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}

}

func (s *Service) GetUserByUID(ctx context.Context, userUID string) *userModels.GetUserResponse {

	var (
		user sqlc.User
		err  error
	)

	if userUID == "" {
		return &userModels.GetUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing Parameters",
				ErrorMsg:  "Missing parameters, user uid could not be found.",
			},
		}
	}

	user, err = s.UserRepo.GetUserByUid(ctx, userUID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "User does not exist.",
					ErrorMsg:  err.Error(),
				},
			}

		} else {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "Unknown error.",
					ErrorMsg:  err.Error(),
				},
			}
		}
	}

	return &userModels.GetUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}
}

func (s *Service) DeleteUser(ctx context.Context, req userModels.DeleteUserRequest) *userModels.DeleteUserResponse {
	// validation
	if req.Email == "" {
		return &userModels.DeleteUserResponse{
			BaseResponse: &models.MissingParamError,
		}
	}

	// check if user exists
	_, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err == sql.ErrNoRows {
		return &userModels.DeleteUserResponse{
			BaseResponse: &models.UserDoesNotExistError,
		}
	}

	// delete user
	user, err := s.UserRepo.DeleteUser(ctx, req.Email)
	if err != nil {
		return &userModels.DeleteUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
			},
		}
	}

	return &userModels.DeleteUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}
}

func (s *Service) UpdateUserPermissionLevel(ctx context.Context, req userModels.UpdateUserPermissionLevelRequest) *userModels.UpdateUserPermissionLevelResponse {
	// validate payload
	if req.Email == "" {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.MissingParamError,
		}
	}

	if req.PermissionLevel == "" {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.MissingParamError,
		}
	}

	// check if user exists
	_, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err == sql.ErrNoRows {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.UserDoesNotExistError,
		}
	} else if err != nil {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
			},
		}
	}

	params := sqlc.UpdateUserPermissionLevelParams{
		Type:  string(req.PermissionLevel),
		Email: req.Email,
	}

	// update user permission
	user, err := s.UserRepo.UpdateUserPermissionLevel(ctx, params)
	if err != nil {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
			},
		}
	}

	return &userModels.UpdateUserPermissionLevelResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
		},
	}
}

func HashPassword(plainPassword string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
