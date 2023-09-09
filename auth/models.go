package auth

import (
	"github.com/Bruary/staff-scheduling/core/models"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Token        string               `json:"token,omitempty"`
}

type TokenType string

type CreateTokenRequest struct {
	UserUid string `json:"user_uid,omitempty"`
}

type CreateTokenResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Token        string               `json:"token,omitempty"`
}

type IsTokenValidRequest struct {
	Token string `json:"token,omitempty"`
}

type IsTokenValidResponse struct {
	BaseResponse *models.BaseResponse `json:"base_response,omitempty"`
	Valid        bool                 `json:"valid,omitempty"`
	Claims       *JWTClaims           `json:"claims,omitempty"`
}
