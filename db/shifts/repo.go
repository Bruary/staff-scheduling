package shifts

import (
	"context"
	"database/sql"

	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
)

// mockery --dir=db/shifts --name=ShiftsRepoInterface --output=db/shifts/mocks --case=underscore
type ShiftsRepoInterface interface {
	CreateShift(ctx context.Context, params sqlc.CreateShiftParams) (sqlc.Shift, error)
	GetUserShifts(ctx context.Context, email string) ([]sqlc.Shift, error)
	DeleteShift(ctx context.Context, uid string) (sqlc.Shift, error)
}

type ShiftsRepo struct {
	provider *sqlc.Queries
}

var _ ShiftsRepoInterface = &ShiftsRepo{}

func New(shiftsRepo *sql.DB) ShiftsRepoInterface {
	return &ShiftsRepo{
		provider: sqlc.New(shiftsRepo),
	}
}

var _ ShiftsRepoInterface = &ShiftsRepo{}

func (s *ShiftsRepo) CreateShift(ctx context.Context, params sqlc.CreateShiftParams) (sqlc.Shift, error) {
	return s.provider.CreateShift(ctx, params)
}

func (s *ShiftsRepo) GetUserShifts(ctx context.Context, email string) ([]sqlc.Shift, error) {
	return s.provider.GetUserShifts(ctx, email)
}

func (s *ShiftsRepo) DeleteShift(ctx context.Context, uid string) (sqlc.Shift, error) {
	return s.provider.DeleteShift(ctx, uid)
}
