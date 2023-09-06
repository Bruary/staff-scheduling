package db

import (
	"context"
	"database/sql"

	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
)

type DbQueries struct {
	provider *sqlc.Queries
}

func New(db *sql.DB) *DbQueries {
	return &DbQueries{
		provider: sqlc.New(db),
	}
}

func (s *DbQueries) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	return s.provider.CreateUser(ctx, params)
}

func (s *DbQueries) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return s.provider.GetUserByEmail(ctx, email)
}

func (s *DbQueries) GetUserByUid(ctx context.Context, uid string) (sqlc.User, error) {
	return s.provider.GetUserByUid(ctx, uid)
}
