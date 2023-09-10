package models

import (
	"time"

	"github.com/Bruary/staff-scheduling/core/models"
)

type Shift struct {
	Id                 int32     `json:"id,omitempty"`
	Created            string    `json:"created,omitempty"`
	Uid                string    `json:"uid,omitempty"`
	WorkDate           time.Time `json:"work_date"`
	ShiftLenghtInHours float32   `json:"shift_length_in_hours"`
	UserId             int       `json:"user_id"`
	Updated            string    `json:"updated,omitempty"`
	Deleted            string    `json:"deleted,omitempty"`
}
type CreateShiftRequest struct {
	WorkDate           string  `json:"work_date"`
	UserEmail          string  `json:"user_email"`
	ShiftLenghtInHours float32 `json:"shift_length_in_hours"`
}

type CreateShiftResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Schedule     *Shift               `json:"schedule,omitempty"`
}

type DeleteShiftRequest struct {
	ShiftUid string `json:"shift_uid"`
}

type DeleteShiftResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Shift        *Shift               `json:"shift,omitempty"`
}

type UpdateShiftRequest struct {
	ShiftUid           string  `json:"shift_uid,omitempty"`
	WorkDate           string  `json:"work_date,omitempty"`
	ShiftLenghtInHours float32 `json:"shift_length_in_hours,omitempty"`
	AssignedUserEmail  string  `json:"assigned_user_email,omitempty"`
}

type UpdateShiftResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Shift        *Shift               `json:"shift,omitempty"`
}
