package db

import (
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
