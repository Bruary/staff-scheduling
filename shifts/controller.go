package shifts

import (
	"context"

	shiftsModels "github.com/Bruary/staff-scheduling/shifts/models"
)

type ControllerInterface interface {
	CreateShift(ctx context.Context, req shiftsModels.CreateShiftRequest) *shiftsModels.CreateShiftResponse
	DeleteShift(ctx context.Context, req shiftsModels.DeleteShiftRequest) *shiftsModels.DeleteShiftResponse
	UpdateShift(ctx context.Context, req shiftsModels.UpdateShiftRequest) *shiftsModels.UpdateShiftResponse
	GetShifts(ctx context.Context, req shiftsModels.GetShiftsRequest) *shiftsModels.GetShiftsResponse
}

type ControllerService struct {
	shiftsService ServiceInterface
}

var _ ControllerInterface = &ControllerService{}

func NewControllerService(shiftsService ServiceInterface) ControllerInterface {
	return &ControllerService{
		shiftsService: shiftsService,
	}
}

// @Title Create shift
// @Summary Create new shift
// @ID create-new-shift
// @Produce json
// @Param req body shiftsModels.CreateShiftRequest true "create shift request"
// @Success 200 {object} shiftsModels.CreateShiftResponse
// @Router /api/v1/shift [post]
func (s *ControllerService) CreateShift(ctx context.Context, req shiftsModels.CreateShiftRequest) *shiftsModels.CreateShiftResponse {
	return s.shiftsService.CreateShift(ctx, req)
}

// @Title Delete shift
// @Summary Delete shift
// @ID delete-shift
// @Produce json
// @Param req body shiftsModels.DeleteShiftRequest true "delete shift request"
// @Success 200 {object} shiftsModels.DeleteShiftResponse
// @Router /api/v1/shift [delete]
func (s *ControllerService) DeleteShift(ctx context.Context, req shiftsModels.DeleteShiftRequest) *shiftsModels.DeleteShiftResponse {
	return s.shiftsService.DeleteShift(ctx, req)
}

// @Title Update shift
// @Summary Update shift
// @ID update-shift
// @Produce json
// @Param req body shiftsModels.UpdateShiftRequest true "update shift request"
// @Success 200 {object} shiftsModels.UpdateShiftResponse
// @Router /api/v1/shift [patch]
func (s *ControllerService) UpdateShift(ctx context.Context, req shiftsModels.UpdateShiftRequest) *shiftsModels.UpdateShiftResponse {
	return s.shiftsService.UpdateShift(ctx, req)
}

// @Title Get shifts
// @Summary Get shifts
// @ID get-shifts
// @Produce json
// @Param req body shiftsModels.GetShiftsRequest true "gets shifts request"
// @Success 200 {object} shiftsModels.GetShiftsResponse
// @Router /api/v1/shifts [get]
func (s *ControllerService) GetShifts(ctx context.Context, req shiftsModels.GetShiftsRequest) *shiftsModels.GetShiftsResponse {
	return s.shiftsService.GetShifts(ctx, req)
}
