package users

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Bruary/staff-scheduling/core/models"
	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	userRepo "github.com/Bruary/staff-scheduling/db/users"
	"github.com/Bruary/staff-scheduling/shifts"
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
	GetAllUsersWithShifts(ctx context.Context, req userModels.GetAllUsersRequest) *userModels.GetAllUsersResponse
}

type Service struct {
	UserRepo      userRepo.UsersRepoInterface
	ShiftsService shifts.ServiceInterface
}

var _ ServiceInterface = &Service{}

func New(userRepo userRepo.UsersRepoInterface, ShiftsService shifts.ServiceInterface) *Service {
	return &Service{
		UserRepo:      userRepo,
		ShiftsService: ShiftsService,
	}
}

var _ ServiceInterface = &Service{}

func (s *Service) CreateUser(ctx context.Context, req userModels.CreateUserRequest) *userModels.CreateUserResponse {

	if req.Email == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User email is missing"),
			},
		}
	}

	if req.FirstName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User first_name is missing"),
			},
		}
	}

	if req.LastName == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User last_name is missing"),
			},
		}
	}

	if req.Password == "" {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User password is missing"),
			},
		}
	}

	resp, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
			},
		}
	} else if (err == nil && resp != sqlc.User{}) {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UserAlreadyExistError.ErrorType,
				ErrorMsg:   models.UserAlreadyExistError.ErrorMsg,
				ErrorStack: append(models.UserAlreadyExistError.ErrorStack, "User is already registered, use login instead"),
			},
		}
	}

	// hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return &userModels.CreateUserResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
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
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, err.Error()),
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
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User email is missing"),
			},
		}
	}

	user, err = s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: models.UserDoesNotExistError.ErrorType,
					ErrorMsg:  models.UserDoesNotExistError.ErrorMsg,
				},
			}

		} else {
			fmt.Printf("Service.GetUserByEmail: failed to get user, user_email=%s err=%s", req.Email, err.Error())

			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: models.UnknownError.ErrorType,
					ErrorMsg:  models.UnknownError.ErrorMsg,
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
				ErrorType:  models.MissingParamError.ErrorType,
				ErrorMsg:   models.MissingParamError.ErrorMsg,
				ErrorStack: append(models.MissingParamError.ErrorStack, "User uid is missing"),
			},
		}
	}

	user, err = s.UserRepo.GetUserByUid(ctx, userUID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: models.UserDoesNotExistError.ErrorType,
					ErrorMsg:  models.UserDoesNotExistError.ErrorMsg,
				},
			}

		} else {
			fmt.Printf("Service.GetUserByUID: failed to get user, user_uid=%s err=%s", userUID, err.Error())

			return &userModels.GetUserResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: models.UnknownError.ErrorType,
					ErrorMsg:  models.UnknownError.ErrorMsg,
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

	if req.PermissionLevel != models.AdminPermissionLevel && req.PermissionLevel != models.BasicPermissionLevel {
		return &userModels.UpdateUserPermissionLevelResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.InvalidParamError.ErrorType,
				ErrorMsg:   models.InvalidParamError.ErrorMsg,
				ErrorStack: append(models.InvalidParamError.ErrorStack, "permission level does not exist"),
			},
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

func (s *Service) GetAllUsersWithShifts(ctx context.Context, req userModels.GetAllUsersRequest) *userModels.GetAllUsersResponse {
	var (
		fromDate       time.Time
		toDate         time.Time
		oneYearAgoDate = time.Now().UTC().Add(-time.Hour * 24 * 365)
		err            error
		resp           []userModels.UserWithShifts
	)

	// validation
	if req.FromDate != "" {
		fromDate, err = time.Parse(models.YYYYMMDD_format, req.FromDate)
		if err != nil {
			return &userModels.GetAllUsersResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType:  models.InvalidParamError.ErrorType,
					ErrorMsg:   models.InvalidParamError.ErrorMsg,
					ErrorStack: append(models.InvalidParamError.ErrorStack, fmt.Sprintf("Failed to parse from_date, format should be like: %s", models.YYYYMMDD_format)),
				},
			}
		}
	}

	// if from_date is not passed OR its more than 1 year, set default to 1 year ago
	if req.FromDate == "" || (req.FromDate != "" && fromDate.UTC().Sub(oneYearAgoDate) < 0) {
		fromDate = time.Now().UTC().Add(-time.Hour * 24 * 365)
	}

	if req.ToDate != "" {
		toDate, err = time.Parse(models.YYYYMMDD_format, req.ToDate)
		if err != nil {
			return &userModels.GetAllUsersResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType:  models.InvalidParamError.ErrorType,
					ErrorMsg:   models.InvalidParamError.ErrorMsg,
					ErrorStack: append(models.InvalidParamError.ErrorStack, fmt.Sprintf("Failed to parse to_date, format should be like: %s", models.YYYYMMDD_format)),
				},
			}
		}
	} else {
		// if to_date is not provided, then default is now
		toDate = time.Now().UTC()
	}

	// check if toDate is smaller than from date
	difference := toDate.UTC().Sub(fromDate)
	if difference < 0 {
		return &userModels.GetAllUsersResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.InvalidDateRangeError.ErrorType,
				ErrorMsg:   models.InvalidDateRangeError.ErrorMsg,
				ErrorStack: append(models.InvalidDateRangeError.ErrorStack, "to_date can not be before from_date"),
			},
		}
	}

	// get all users with shifts
	usersWithShifts, err := s.UserRepo.GetAllUsersWithShifts(ctx, sqlc.GetAllUsersWithShiftsParams{
		WorkDate:   fromDate,
		WorkDate_2: toDate,
	})
	if err != nil {
		return &userModels.GetAllUsersResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.UnknownError.ErrorType,
				ErrorMsg:   models.UnknownError.ErrorMsg,
				ErrorStack: append(models.UnknownError.ErrorStack, fmt.Sprintf("Failed to get all users with shifts, err=%s", err.Error())),
			},
		}
	}

	// if no records, just return empty array
	if len(usersWithShifts) == 0 {
		return &userModels.GetAllUsersResponse{
			UserWithShifts: []userModels.UserWithShifts{},
		}
	}

	// fill response
	for _, userWithShift := range usersWithShifts {
		resp = append(resp, userModels.UserWithShifts{
			User: userModels.User{
				Id:        userWithShift.ID,
				Created:   userWithShift.Created.String(),
				Type:      userWithShift.Type,
				Uid:       userWithShift.Uid,
				FirstName: userWithShift.FirstName,
				LastName:  userWithShift.LastName,
				Email:     userWithShift.Email,
			},
			ShiftsLengthInHours: float32(userWithShift.TotalHours),
		})
	}

	return &userModels.GetAllUsersResponse{
		UserWithShifts: resp,
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
