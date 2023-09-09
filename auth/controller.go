package auth

import (
	"context"
)

type ControllerInterface interface {
	Login(ctx context.Context, req LoginRequest) *LoginResponse
	CreateToken(ctx context.Context, req CreateTokenRequest) *CreateTokenResponse
	IsTokenValid(ctx context.Context, req IsTokenValidRequest) *IsTokenValidResponse
}

type ControllerService struct {
	authService ServiceInterface
}

var _ ControllerInterface = &ControllerService{}

func NewControllerService(authService ServiceInterface) *ControllerService {
	return &ControllerService{
		authService: authService,
	}
}

// @Title Login
// @Summary Login
// @ID login
// @Produce json
// @Param req body LoginRequest true "login request"
// @Success 200 {object} LoginResponse
// @Router /api/v1/login [post]
func (s *ControllerService) Login(ctx context.Context, req LoginRequest) *LoginResponse {
	return s.authService.Login(ctx, req)
}

func (s *ControllerService) CreateToken(ctx context.Context, req CreateTokenRequest) *CreateTokenResponse {
	return s.authService.CreateToken(ctx, req)
}

func (s *ControllerService) IsTokenValid(ctx context.Context, req IsTokenValidRequest) *IsTokenValidResponse {
	return s.authService.IsTokenValid(ctx, req)
}
