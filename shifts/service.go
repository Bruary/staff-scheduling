package shifts

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	coreModels "github.com/Bruary/staff-scheduling/core/models"
	shiftsRepo "github.com/Bruary/staff-scheduling/db/shifts"
	sqlc "github.com/Bruary/staff-scheduling/db/sqlc"
	usersRepo "github.com/Bruary/staff-scheduling/db/users"
	shiftsModels "github.com/Bruary/staff-scheduling/shifts/models"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	CreateShift(ctx context.Context, req shiftsModels.CreateShiftRequest) *shiftsModels.CreateShiftResponse
}

type Service struct {
	usersRepo  usersRepo.UsersRepoInterface
	shiftsRepo shiftsRepo.ShiftsRepoInterface
}

var _ ServiceInterface = &Service{}

func New(usersRepo usersRepo.UsersRepoInterface, shiftsRepo shiftsRepo.ShiftsRepoInterface) Service {
	return Service{
		usersRepo:  usersRepo,
		shiftsRepo: shiftsRepo,
	}
}

func (s *Service) CreateShift(ctx context.Context, req shiftsModels.CreateShiftRequest) *shiftsModels.CreateShiftResponse {

	// validate
	if req.UserEmail == "" {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Service.CreateSchedule: user email is missing"),
			},
		}
	}

	if req.ShiftLenghtInHours == 0 {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Service.CreateSchedule: shift length is missing"),
			},
		}
	}

	if req.ShiftLenghtInHours < 0 || req.ShiftLenghtInHours > 12 {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.InvalidParamError.ErrorType,
				ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
				ErrorStack: append(coreModels.InvalidParamError.ErrorStack, "Service.CreateSchedule: shift length can not be less than 0.5 hours or more than 12 hours"),
			},
		}
	}

	if req.WorkDate == "" {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Service.CreateSchedule: work date is missing"),
			},
		}
	}

	// check if date is valid
	newShiftDate, err := time.Parse(coreModels.YYYYMMDD_format, req.WorkDate)
	if err != nil {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.InvalidParamError.ErrorType,
				ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
				ErrorStack: append(coreModels.InvalidParamError.ErrorStack, fmt.Sprintf("Service.CreateSchedule: work date is invalid, correct format:%s ", coreModels.YYYYMMDD_format)),
			},
		}
	}

	// check if user exists
	user, err := s.usersRepo.GetUserByEmail(ctx, req.UserEmail)
	if err != nil && err == sql.ErrNoRows {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.UserDoesNotExistError,
		}
	} else if err != nil {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.UnknownError.ErrorType,
				ErrorMsg:   coreModels.UnknownError.ErrorMsg,
				ErrorStack: append(coreModels.UnknownError.ErrorStack, "Service.CreateSchedule: failed to get user by email, err="+err.Error()),
			},
		}
	}

	// get all existing shifts for this user
	shifts, err := s.shiftsRepo.GetUserShifts(ctx, req.UserEmail)
	if err != nil {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.UnknownError.ErrorType,
				ErrorMsg:   coreModels.UnknownError.ErrorMsg,
				ErrorStack: append(coreModels.UnknownError.ErrorStack, "Service.CreateSchedule: failed to get user shifts, err="+err.Error()),
			},
		}
	}

	// check if new schedule overlaps with existing schedule
	for _, shift := range shifts {
		doesUserHaveShiftAlready := newShiftDate.Equal(shift.WorkDate)

		if doesUserHaveShiftAlready {
			return &shiftsModels.CreateShiftResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.ScheduleAlreadyExistError.ErrorType,
					ErrorMsg:   coreModels.ScheduleAlreadyExistError.ErrorMsg,
					ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Service.CreateSchedule: shift already exist on %s for user_id=%v", shift.WorkDate, shift.UserID)),
				},
			}
		}
	}

	// create shift
	params := sqlc.CreateShiftParams{
		Uid:              uuid.NewString(),
		WorkDate:         newShiftDate,
		ShiftLengthHours: math.Round(float64(req.ShiftLenghtInHours)*100) / 100,
		UserID:           user.ID,
	}

	shift, err := s.shiftsRepo.CreateShift(ctx, params)
	if err != nil {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.UnknownError.ErrorType,
				ErrorMsg:   coreModels.UnknownError.ErrorMsg,
				ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Service.CreateSchedule: failed to create a new shift, user_id=%v", shift.UserID)),
			},
		}
	}

	return &shiftsModels.CreateShiftResponse{
		Schedule: &shiftsModels.Shift{
			Id:                 shift.ID,
			Created:            shift.Created.String(),
			Uid:                shift.Uid,
			WorkDate:           shift.WorkDate,
			ShiftLenghtInHours: float32(shift.ShiftLengthHours),
			UserId:             int(shift.UserID),
		},
	}
}
