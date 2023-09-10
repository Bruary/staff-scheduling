package models

type BaseResponse struct {
	Success    bool     `json:"success,omitempty"`
	ErrorType  string   `json:"error_type,omitempty"`
	ErrorMsg   string   `json:"error_msg,omitempty"`
	ErrorStack []string `json:"error_stack,omitempty"`
}

type PermissionLevel string

var (
	AdminPermissionLevel PermissionLevel = "admin"
	BasicPermissionLevel PermissionLevel = "basic"
)

// Date formats
var (
	YYYYMMDD_format = "2006-01-02"
)
