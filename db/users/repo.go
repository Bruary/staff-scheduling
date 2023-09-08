package users

import (
	"context"

	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
)

// mockery --dir=db/users --name=UsersRepoInterface --output=db/users/mocks --case=underscore
type UsersRepoInterface interface {
	CreateUser(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	GetUserByUid(ctx context.Context, uid string) (sqlc.User, error)
}

type UsersRepo struct {
	provider *sqlc.Queries
}

var _ UsersRepoInterface = &UsersRepo{}

func (s *UsersRepo) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	return s.provider.CreateUser(ctx, params)
}

func (s *UsersRepo) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return s.provider.GetUserByEmail(ctx, email)
}

func (s *UsersRepo) GetUserByUid(ctx context.Context, uid string) (sqlc.User, error) {
	return s.provider.GetUserByUid(ctx, uid)
}
