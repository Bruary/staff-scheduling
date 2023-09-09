package models

type BaseResponse struct {
	Success    bool     `json:"success,omitempty"`
	ErrorType  string   `json:"error_type,omitempty"`
	ErrorMsg   string   `json:"error_msg,omitempty"`
	ErrorStack []string `json:"error_stack,omitempty"`
}

type PermissionLevel string

var (
	Admin PermissionLevel = "admin"
	Basic PermissionLevel = "basic"
)
