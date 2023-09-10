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
	UpdateShiftUserId(ctx context.Context, params sqlc.UpdateShiftUserIdParams) (sqlc.Shift, error)
	UpdateShiftWorkDate(ctx context.Context, params sqlc.UpdateShiftWorkDateParams) (sqlc.Shift, error)
	UpdateShiftLength(ctx context.Context, params sqlc.UpdateShiftLengthParams) (sqlc.Shift, error)
	GetShiftByUid(ctx context.Context, uid string) (sqlc.Shift, error)
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

func (s *ShiftsRepo) UpdateShiftUserId(ctx context.Context, params sqlc.UpdateShiftUserIdParams) (sqlc.Shift, error) {
	return s.provider.UpdateShiftUserId(ctx, params)
}

func (s *ShiftsRepo) UpdateShiftWorkDate(ctx context.Context, params sqlc.UpdateShiftWorkDateParams) (sqlc.Shift, error) {
	return s.provider.UpdateShiftWorkDate(ctx, params)
}

func (s *ShiftsRepo) UpdateShiftLength(ctx context.Context, params sqlc.UpdateShiftLengthParams) (sqlc.Shift, error) {
	return s.provider.UpdateShiftLength(ctx, params)
}

func (s *ShiftsRepo) GetShiftByUid(ctx context.Context, uid string) (sqlc.Shift, error) {
	return s.provider.GetShiftByUid(ctx, uid)
}
