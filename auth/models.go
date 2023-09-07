package auth

import "github.com/Bruary/staff-scheduling/models"

type LoginRequest struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type LoginResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Token        Token                `json:"Token,omitempty"`
}

type TokenType string

var JWT TokenType = "jwt"

type Token struct {
	Token     string    `json:"Token,omitempty"`
	TokenType TokenType `json:"token_type,omitempty"`
}

type CreateTokenRequest struct {
	TokenType TokenType `json:"token_type,omitempty"`
	Email     string    `json:"email,omitempty"`
}

type CreateTokenResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Token        Token                `json:"token,omitempty"`
}

type IsTokenValidRequest struct {
	Token *Token `json:"token,omitempty"`
}

type IsTokenValidResponse struct {
	BaseResponse *models.BaseResponse `json:"baseResponse,omitempty"`
	Valid        bool                 `json:"Valid,omitempty"`
}
