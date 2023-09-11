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

// mockery --dir=shifts --name=ServiceInterface --output=shifts/mocks --case=underscore
type ServiceInterface interface {
	CreateShift(ctx context.Context, req shiftsModels.CreateShiftRequest) *shiftsModels.CreateShiftResponse
	DeleteShift(ctx context.Context, req shiftsModels.DeleteShiftRequest) *shiftsModels.DeleteShiftResponse
	UpdateShift(ctx context.Context, req shiftsModels.UpdateShiftRequest) *shiftsModels.UpdateShiftResponse
	GetShifts(ctx context.Context, req shiftsModels.GetShiftsRequest) *shiftsModels.GetShiftsResponse
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
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "User email is missing"),
			},
		}
	}

	if req.ShiftLenghtInHours == 0 {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Shift length is missing"),
			},
		}
	}

	if req.ShiftLenghtInHours < 0 || req.ShiftLenghtInHours > 12 {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.InvalidParamError.ErrorType,
				ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
				ErrorStack: append(coreModels.InvalidParamError.ErrorStack, "Shift length can not be less than 0.5 hours or more than 12 hours"),
			},
		}
	}

	if req.WorkDate == "" {
		return &shiftsModels.CreateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Work date is missing"),
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
				ErrorStack: append(coreModels.InvalidParamError.ErrorStack, fmt.Sprintf("Work date is invalid, correct format:%s ", coreModels.YYYYMMDD_format)),
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
				ErrorStack: append(coreModels.UnknownError.ErrorStack, "Wailed to get user by email, err="+err.Error()),
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
				ErrorStack: append(coreModels.UnknownError.ErrorStack, "Failed to get user shifts, err="+err.Error()),
			},
		}
	}

	// check if new schedule overlaps with existing schedule
	for _, shift := range shifts {
		doesUserHaveShiftAlready := newShiftDate.Equal(shift.WorkDate)

		if doesUserHaveShiftAlready {
			return &shiftsModels.CreateShiftResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.ShiftAlreadyExistError.ErrorType,
					ErrorMsg:   coreModels.ShiftAlreadyExistError.ErrorMsg,
					ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Shift already exist on %s for user_id=%v", shift.WorkDate, shift.UserID)),
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
				ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to create a new shift, user_id=%v", shift.UserID)),
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

func (s *Service) DeleteShift(ctx context.Context, req shiftsModels.DeleteShiftRequest) *shiftsModels.DeleteShiftResponse {

	// validate
	if req.ShiftUid == "" {
		return &shiftsModels.DeleteShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Shift uid is missing"),
			},
		}
	}

	shift, err := s.shiftsRepo.DeleteShift(ctx, req.ShiftUid)
	if err != nil {
		return &shiftsModels.DeleteShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.UnknownError.ErrorType,
				ErrorMsg:   coreModels.UnknownError.ErrorMsg,
				ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to delete shift or its already deleted, shift_uid=%s, err=%s", req.ShiftUid, err.Error())),
			},
		}
	}

	return &shiftsModels.DeleteShiftResponse{
		Shift: &shiftsModels.Shift{
			Id:                 shift.ID,
			Created:            shift.Created.String(),
			Uid:                shift.Uid,
			WorkDate:           shift.WorkDate,
			ShiftLenghtInHours: float32(shift.ShiftLengthHours),
			UserId:             int(shift.UserID),
			Updated:            shift.Updated.Time.String(),
			Deleted:            shift.Deleted.Time.String(),
		},
	}
}

func (s *Service) UpdateShift(ctx context.Context, req shiftsModels.UpdateShiftRequest) *shiftsModels.UpdateShiftResponse {

	var (
		updatedShift sqlc.Shift
	)

	// validate
	if req.ShiftUid == "" {
		return &shiftsModels.UpdateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Missing shift_uid in the request"),
			},
		}
	}

	if req.AssignedUserEmail == "" && req.ShiftLenghtInHours == 0 && req.WorkDate == "" {
		return &shiftsModels.UpdateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Atleast one of the following params are required: work_date, shift_length_in_hours, assigned_user_email"),
			},
		}
	}

	// check if shift exists
	shift, err := s.shiftsRepo.GetShiftByUid(ctx, req.ShiftUid)
	if err != nil {
		return &shiftsModels.UpdateShiftResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.ShiftNotFoundError.ErrorType,
				ErrorMsg:   coreModels.ShiftNotFoundError.ErrorMsg,
				ErrorStack: append(coreModels.ShiftNotFoundError.ErrorStack, fmt.Sprintf("Failed to get shift, shift_uid=%s, err=%s", req.ShiftUid, err.Error())),
			},
		}
	}

	// if user email is passed, then:
	// 1) check if new user exists
	// 2) update shift with new user
	if req.AssignedUserEmail != "" {
		user, err := s.usersRepo.GetUserByEmail(ctx, req.AssignedUserEmail)
		if err != nil {
			return &shiftsModels.UpdateShiftResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.UnknownError.ErrorType,
					ErrorMsg:   coreModels.UnknownError.ErrorMsg,
					ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to get user, user_email=%s, err=%s", req.AssignedUserEmail, err.Error())),
				},
			}
		}

		// update
		if shift.UserID != user.ID {
			updatedShift, err = s.shiftsRepo.UpdateShiftUserId(ctx, sqlc.UpdateShiftUserIdParams{
				Email: req.AssignedUserEmail,
				Uid:   req.ShiftUid,
			})
			if err != nil {
				return &shiftsModels.UpdateShiftResponse{
					BaseResponse: &coreModels.BaseResponse{
						ErrorType:  coreModels.UnknownError.ErrorType,
						ErrorMsg:   coreModels.UnknownError.ErrorMsg,
						ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to update shift user id, shift_uid=%s, user_email=%s, err=%s", req.ShiftUid, req.AssignedUserEmail, err.Error())),
					},
				}
			}

			fmt.Printf("Service.UpdateShift: shift user_id have been updated, shift_uid=%s, old user_id=%v, new user_id=%v", req.ShiftUid, shift.UserID, user.ID)
		}
	}

	// if work date is passed, then:
	// 1) parse new date
	// 2) update shift with new work date
	if req.WorkDate != "" {
		newWorkDate, err := time.Parse(coreModels.YYYYMMDD_format, req.WorkDate)
		if err != nil {
			return &shiftsModels.UpdateShiftResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.UnknownError.ErrorType,
					ErrorMsg:   coreModels.UnknownError.ErrorMsg,
					ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to parse new work date, correct date format:%s shift_uid=%s, work_date=%s, err=%s", coreModels.YYYYMMDD_format, req.ShiftUid, req.WorkDate, err.Error())),
				},
			}
		}

		if !shift.WorkDate.Equal(newWorkDate) {
			updatedShift, err = s.shiftsRepo.UpdateShiftWorkDate(ctx, sqlc.UpdateShiftWorkDateParams{
				WorkDate: newWorkDate,
				Uid:      req.ShiftUid,
			})
			if err != nil {
				return &shiftsModels.UpdateShiftResponse{
					BaseResponse: &coreModels.BaseResponse{
						ErrorType:  coreModels.UnknownError.ErrorType,
						ErrorMsg:   coreModels.UnknownError.ErrorMsg,
						ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed update shift work date, shift_uid=%s, work_date=%s, err=%s", req.ShiftUid, req.WorkDate, err.Error())),
					},
				}
			}

			fmt.Printf("Service.UpdateShift: shift work_date have been updated, shift_uid=%s, old work_date=%s, new work_date=%s", req.ShiftUid, shift.WorkDate.String(), newWorkDate.String())
		}
	}

	// if shift length is passed, then:
	// 1) validate new shift length
	// 2) update shift with new length
	if req.ShiftLenghtInHours != 0 {

		if req.ShiftLenghtInHours > 12 || req.ShiftLenghtInHours < 0.5 {
			return &shiftsModels.UpdateShiftResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.InvalidParamError.ErrorType,
					ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
					ErrorStack: append(coreModels.InvalidParamError.ErrorStack, "Shift_length_in_hours can not be less than 0.5 hours or more than 12 hours"),
				},
			}
		}

		if req.ShiftLenghtInHours != float32(shift.ShiftLengthHours) {
			updatedShift, err = s.shiftsRepo.UpdateShiftLength(ctx, sqlc.UpdateShiftLengthParams{
				ShiftLengthHours: float64(req.ShiftLenghtInHours),
				Uid:              req.ShiftUid,
			})
			if err != nil {
				return &shiftsModels.UpdateShiftResponse{
					BaseResponse: &coreModels.BaseResponse{
						ErrorType:  coreModels.UnknownError.ErrorType,
						ErrorMsg:   coreModels.UnknownError.ErrorMsg,
						ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed update shift length, shift_uid=%s, shift_length_in_hours=%v, err=%s", req.ShiftUid, req.ShiftLenghtInHours, err.Error())),
					},
				}
			}

			fmt.Printf("Service.UpdateShift: shift shift_length_in_hours have been updated, shift_uid=%s, old shift_length_in_hours=%v, new shift_length_in_hours=%v", req.ShiftUid, shift.ShiftLengthHours, req.ShiftLenghtInHours)
		}
	}

	return &shiftsModels.UpdateShiftResponse{
		Shift: &shiftsModels.Shift{
			Id:                 updatedShift.ID,
			Created:            updatedShift.Created.String(),
			Uid:                updatedShift.Uid,
			WorkDate:           updatedShift.WorkDate,
			ShiftLenghtInHours: float32(updatedShift.ShiftLengthHours),
			UserId:             int(updatedShift.UserID),
			Updated:            updatedShift.Updated.Time.String(),
		},
	}
}

func (s *Service) GetShifts(ctx context.Context, req shiftsModels.GetShiftsRequest) *shiftsModels.GetShiftsResponse {
	var (
		resp           []shiftsModels.Shift
		fromDate       time.Time
		toDate         time.Time
		oneYearAgoDate = time.Now().UTC().Add(-time.Hour * 24 * 365)
		err            error
	)

	// validate
	if len(req.UsersEmails) == 0 {
		return &shiftsModels.GetShiftsResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.MissingParamError.ErrorType,
				ErrorMsg:   coreModels.MissingParamError.ErrorMsg,
				ErrorStack: append(coreModels.MissingParamError.ErrorStack, "Missing user_emails in the request"),
			},
		}
	}

	// Validate dates
	if req.FromDate != "" {
		fromDate, err = time.Parse(coreModels.YYYYMMDD_format, req.FromDate)
		if err != nil {
			return &shiftsModels.GetShiftsResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.InvalidParamError.ErrorType,
					ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
					ErrorStack: append(coreModels.InvalidParamError.ErrorStack, fmt.Sprintf("Failed to parse from_date, format should be like: %s", coreModels.YYYYMMDD_format)),
				},
			}
		}
	}

	// if from_date is not passed OR its more than 1 year, set default to 1 year ago
	if req.FromDate == "" || (req.FromDate != "" && fromDate.UTC().Sub(oneYearAgoDate) < 0) {
		fromDate = time.Now().UTC().Add(-time.Hour * 24 * 365)
	}

	if req.ToDate != "" {
		toDate, err = time.Parse(coreModels.YYYYMMDD_format, req.ToDate)
		if err != nil {
			return &shiftsModels.GetShiftsResponse{
				BaseResponse: &coreModels.BaseResponse{
					ErrorType:  coreModels.InvalidParamError.ErrorType,
					ErrorMsg:   coreModels.InvalidParamError.ErrorMsg,
					ErrorStack: append(coreModels.InvalidParamError.ErrorStack, fmt.Sprintf("Failed to parse to_date, format should be like: %s", coreModels.YYYYMMDD_format)),
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
		return &shiftsModels.GetShiftsResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.InvalidDateRangeError.ErrorType,
				ErrorMsg:   coreModels.InvalidDateRangeError.ErrorMsg,
				ErrorStack: append(coreModels.InvalidDateRangeError.ErrorStack, "to_date can not be before from_date"),
			},
		}
	}

	// fetch all users shifts in the list provided and order them by work_date in desc format
	userShifts, err := s.shiftsRepo.GetUsersShiftsByDateRange(ctx, sqlc.GetUsersShiftsByDateRangeParams{
		Column1:    req.UsersEmails,
		WorkDate:   fromDate,
		WorkDate_2: toDate,
	})
	if err != nil {
		return &shiftsModels.GetShiftsResponse{
			BaseResponse: &coreModels.BaseResponse{
				ErrorType:  coreModels.UnknownError.ErrorType,
				ErrorMsg:   coreModels.UnknownError.ErrorMsg,
				ErrorStack: append(coreModels.UnknownError.ErrorStack, fmt.Sprintf("Failed to fetch users shifts err=%s", err.Error())),
			},
		}
	}

	// if no records, just return empty array
	if len(userShifts) == 0 {
		return &shiftsModels.GetShiftsResponse{
			Shifts: []shiftsModels.Shift{},
		}
	}

	// fill shifts extracted into the response
	for _, shift := range userShifts {
		resp = append(resp, shiftsModels.Shift{
			Id:                 shift.ID,
			Created:            shift.Created.String(),
			Uid:                shift.Uid,
			WorkDate:           shift.WorkDate,
			ShiftLenghtInHours: float32(shift.ShiftLengthHours),
			UserId:             int(shift.UserID),
			Updated:            shift.Updated.Time.String(),
		})
	}

	return &shiftsModels.GetShiftsResponse{
		Shifts: resp,
	}
}
