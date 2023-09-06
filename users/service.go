package users

import (
	"context"
	"database/sql"

	db "github.com/Bruary/staff-scheduling/db"
	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	"github.com/Bruary/staff-scheduling/models"
	userModels "github.com/Bruary/staff-scheduling/users/models"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	CreateUser(context.Context, userModels.CreateUserRequest) (*userModels.CreateUserResponse, error)
	GetUserByEmail(context.Context, userModels.GetUserByEmailRequest) (*userModels.GetUserResponse, error)
	GetUserByUID(ctx context.Context, userUID string) (*userModels.GetUserResponse, error)
	DeleteUser(context.Context, userModels.DeleteUserRequest) (*userModels.DeleteUserResponse, error)
}

type Service struct {
	Db *db.DbQueries
}

var _ ServiceInterface = &Service{}

func New(db *db.DbQueries) *Service {
	return &Service{
		Db: db,
	}
}

var _ ServiceInterface = &Service{}

func (s *Service) CreateUser(ctx context.Context, req userModels.CreateUserRequest) (*userModels.CreateUserResponse, error) {

	if req.Email == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params.",
				ErrorMsg:  "Users email is missing.",
			},
		}, nil
	}

	if req.FirstName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params.",
				ErrorMsg:  "Users first name is missing.",
			},
		}, nil
	}

	if req.LastName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params.",
				ErrorMsg:  "Users last name is missing.",
			},
		}, nil
	}

	if req.Password == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing params.",
				ErrorMsg:  "Password is required.",
			},
		}, nil
	}

	resp, _ := s.GetUserByEmail(ctx, userModels.GetUserByEmailRequest{
		Email: req.Email,
	})

	// check if get user worked, which means user already exists
	if resp.BaseResponse == nil {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "User already exists.",
				ErrorMsg:  "User is already registered",
			},
		}, nil
	}

	params := sqlc.CreateUserParams{
		Uid:       uuid.NewString(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	user, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Unknown error",
				ErrorMsg:  err.Error(),
			},
		}, nil
	}

	return &userModels.CreateUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type.String,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, req userModels.GetUserByEmailRequest) (*userModels.GetUserResponse, error) {

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
		}, nil
	}

	user, err = s.Db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "User does not exist.",
					ErrorMsg:  err.Error(),
				},
			}, nil

		} else {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "Unknown error.",
					ErrorMsg:  err.Error(),
				},
			}, nil
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
	}, nil

}

func (s *Service) GetUserByUID(ctx context.Context, userUID string) (*userModels.GetUserResponse, error) {

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
		}, nil
	}

	user, err = s.Db.GetUserByUid(ctx, userUID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "User does not exist.",
					ErrorMsg:  err.Error(),
				},
			}, nil

		} else {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "Unknown error.",
					ErrorMsg:  err.Error(),
				},
			}, nil
		}
	}

	return &userModels.GetUserResponse{
		User: userModels.User{
			Id:        user.ID,
			Created:   user.Created.String(),
			Type:      user.Type.String,
			Uid:       user.Uid,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Updated:   user.Updated.String(),
			Deleted:   user.Deleted.Time.String(),
		},
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req userModels.DeleteUserRequest) (*userModels.DeleteUserResponse, error) {

	return nil, nil
}
