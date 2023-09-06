package models

type BaseResponse struct {
	ErrorType string `json:"errorType,omitempty"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
}
