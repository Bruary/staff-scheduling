package models

type Errors BaseResponse

var (
	// Validations
	MissingParamError = BaseResponse{Success: false, ErrorType: "MISSING_PARAMETER", ErrorMsg: "Request parameter(s) missing", ErrorStack: []string{}}
	InvalidParamError = BaseResponse{Success: false, ErrorType: "INVALID_PARAMETER", ErrorMsg: "Request parameter(s) is invalid", ErrorStack: []string{}}

	// User
	UserDoesNotExistError = BaseResponse{Success: false, ErrorType: "DOES_NOT_EXIST", ErrorMsg: "User does not exist", ErrorStack: []string{}}

	// Schedules
	ScheduleAlreadyExistError = BaseResponse{Success: false, ErrorType: "SCHEDULE_ALREADY_EXIST", ErrorMsg: "User already have a schedule booked for this day", ErrorStack: []string{}}

	// Core
	UnknownError        = BaseResponse{Success: false, ErrorType: "UNKNOWN_ERROR", ErrorMsg: "Unknown error occured", ErrorStack: []string{}}
	JWTMissingError     = BaseResponse{Success: false, ErrorType: "TOKEN_MISSING", ErrorMsg: "JWT token is missing", ErrorStack: []string{}}
	UnauthorizedError   = BaseResponse{Success: false, ErrorType: "UNAUTHORIZED", ErrorMsg: "Unauthorized access", ErrorStack: []string{}}
	UserPermissionError = BaseResponse{Success: false, ErrorType: "PERMISSION_ERROR", ErrorMsg: "User is not authorized", ErrorStack: []string{}}
)
