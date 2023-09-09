package auth

import "github.com/Bruary/staff-scheduling/models"

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Token        string               `json:"token,omitempty"`
}

type TokenType string

type CreateTokenRequest struct {
	Email string `json:"email,omitempty"`
}

type CreateTokenResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Token        string               `json:"token,omitempty"`
}

type IsTokenValidRequest struct {
	Token string `json:"token,omitempty"`
}

type IsTokenValidResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Valid        bool                 `json:"Valid,omitempty"`
}
