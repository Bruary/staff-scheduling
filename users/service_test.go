package users_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	usersRepoMock "github.com/Bruary/staff-scheduling/db/users/mocks"
	coreModels "github.com/Bruary/staff-scheduling/models"
	"github.com/Bruary/staff-scheduling/users"
	"github.com/Bruary/staff-scheduling/users/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("environment", "testing")
}

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		insertUserParams  sqlc.CreateUserParams
		createUserResp    sqlc.User
		resp              models.CreateUserResponse
		getUserResp       sqlc.User
		getUserErr        error
		userAlreadyExists bool
		errorType         string
	}{
		"Should fail if user already exists": {
			insertUserParams: sqlc.CreateUserParams{
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			resp: models.CreateUserResponse{
				User: models.User{
					Id:        1,
					Type:      string(coreModels.Basic),
					Uid:       users.UidForTests,
					FirstName: "Bakir",
					LastName:  "Kais",
					Email:     "bakirtest@test.com",
					Password:  "helloWorld",
				},
			},
			createUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.Basic),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			getUserResp: sqlc.User{
				ID:        1,
				Type:      string(coreModels.Basic),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest@test.com",
				Password:  "helloWorld",
			},
			userAlreadyExists: true,
			getUserErr:        nil,
			errorType:         "User already exists",
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
					Type:      string(coreModels.Basic),
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
				Type:      string(coreModels.Basic),
				Uid:       users.UidForTests,
				FirstName: "Bakir",
				LastName:  "Kais",
				Email:     "bakirtest_2@test.com",
				Password:  "helloWorld",
			},
			getUserResp:       sqlc.User{},
			getUserErr:        sql.ErrNoRows,
			userAlreadyExists: false,
		},
	}

	userRepoMock := usersRepoMock.UsersRepoInterface{}
	usersService := users.New(&userRepoMock)

	ctx := context.Background()

	for _, tc := range tests {

		userRepoMock.On("GetUserByEmail", ctx, tc.insertUserParams.Email).Return(tc.getUserResp, tc.getUserErr)
		defer userRepoMock.AssertExpectations(t)

		if !tc.userAlreadyExists {
			userRepoMock.On("CreateUser", ctx, tc.insertUserParams).Return(tc.createUserResp, nil)
		}

		resp := usersService.CreateUser(ctx, models.CreateUserRequest{
			FirstName: tc.insertUserParams.FirstName,
			LastName:  tc.insertUserParams.LastName,
			Email:     tc.insertUserParams.Email,
			Password:  tc.insertUserParams.Password,
		})
		if tc.userAlreadyExists {
			assert.Equal(t, resp.BaseResponse.ErrorType, tc.errorType)
		}

		if !tc.userAlreadyExists {
			assert.Equal(t, tc.resp.User.Email, resp.User.Email)
			assert.Equal(t, tc.resp.User.Type, resp.User.Type)
			assert.Equal(t, tc.resp.User.Uid, resp.User.Uid)
		}
	}
}
