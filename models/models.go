package models

type BaseResponse struct {
	ErrorType  string   `json:"error_type,omitempty"`
	ErrorMsg   string   `json:"error_msg,omitempty"`
	ErrorStack []string `json:"error_stack,omitempty"`
}

type PermissionLevel string

var (
	Admin PermissionLevel = "admin"
	Basic PermissionLevel = "basic"
)
