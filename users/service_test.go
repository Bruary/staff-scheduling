package users_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	coreModels "github.com/Bruary/staff-scheduling/core/models"
	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	usersRepoMock "github.com/Bruary/staff-scheduling/db/users/mocks"
	shiftsServiceMock "github.com/Bruary/staff-scheduling/shifts/mocks"
	"github.com/Bruary/staff-scheduling/users"
	"github.com/Bruary/staff-scheduling/users/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("environment", "testing")
}

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		insertUserParams     sqlc.CreateUserParams
		createUserResp       sqlc.User
		resp                 models.CreateUserResponse
		getUserResp          sqlc.User
		getUserErr           error
		shouldCallCreateUser bool
		shouldCallGetUser    bool
		missingParam         bool
		shouldError          bool
		errorType            string
	}{
		"Should fail if user email is missing": {
			insertUserParams: sqlc.CreateUserParams{
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "",
				Password:  "helloWorld",
			},
			createUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.BasicPermissionLevel),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			shouldCallCreateUser: false,
			shouldError:          true,
			errorType:            coreModels.MissingParamError.ErrorType,
		},
		"Should fail if user already exists": {
			insertUserParams: sqlc.CreateUserParams{
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			createUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.BasicPermissionLevel),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			getUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.BasicPermissionLevel),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			shouldCallCreateUser: false,
			shouldError:          true,
			errorType:            coreModels.UserAlreadyExistError.ErrorType,
			shouldCallGetUser:    true,
		},
		"Should successed if all params are passed and valid": {
			insertUserParams: sqlc.CreateUserParams{
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest_2@test.com",
				Password:  "helloWorld",
			},
			resp: models.CreateUserResponse{
				User: models.User{
					Id:        1,
					Created:   time.Now().String(),
					Type:      string(coreModels.BasicPermissionLevel),
					Uid:       users.UidForTests,
					FirstName: "Bakir",
					LastName:  "Kais",
					Email:     "bakirtest_2@test.com",
					Password:  "helloWorld",
					Updated:   time.Now().String(),
				},
			},
			createUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.BasicPermissionLevel),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest_2@test.com",
				Password:  "helloWorld",
			},
			getUserErr:           sql.ErrNoRows,
			shouldError:          false,
			shouldCallCreateUser: true,
			shouldCallGetUser:    true,
		},
	}

	userRepoMock := usersRepoMock.UsersRepoInterface{}
	shiftsServiceMock := shiftsServiceMock.ServiceInterface{}
	usersService := users.New(&userRepoMock, &shiftsServiceMock)

	ctx := context.Background()

	for _, tc := range tests {

		if tc.shouldCallGetUser {
			userRepoMock.On("GetUserByEmail", ctx, tc.insertUserParams.Email).Return(tc.getUserResp, tc.getUserErr)
		}
		defer userRepoMock.AssertExpectations(t)

		if tc.shouldCallCreateUser {
			userRepoMock.On("CreateUser", ctx, tc.insertUserParams).Return(tc.createUserResp, nil)
		}

		resp := usersService.CreateUser(ctx, models.CreateUserRequest{
			FirstName: tc.insertUserParams.FirstName,
			LastName:  tc.insertUserParams.LastName,
			Email:     tc.insertUserParams.Email,
			Password:  tc.insertUserParams.Password,
		})

		if tc.shouldError {
			assert.Equal(t, resp.BaseResponse.ErrorType, tc.errorType)
		}

		if !tc.shouldError {
			assert.Equal(t, tc.resp.User.Email, resp.User.Email)
			assert.Equal(t, tc.resp.User.Type, resp.User.Type)
			assert.Equal(t, tc.resp.User.Uid, resp.User.Uid)
		}
	}
}
