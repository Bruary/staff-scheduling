package auth

import (
	"context"
	"fmt"
	"time"

	models "github.com/Bruary/staff-scheduling/core/models"
	"github.com/Bruary/staff-scheduling/users"
	userModels "github.com/Bruary/staff-scheduling/users/models"
	jwt "github.com/golang-jwt/jwt/v4"
)

type ServiceInterface interface {
	Login(context.Context, LoginRequest) *LoginResponse
	CreateToken(context.Context, CreateTokenRequest) *CreateTokenResponse
	IsTokenValid(context.Context, IsTokenValidRequest) *IsTokenValidResponse
}

type Service struct {
	UsersService users.ServiceInterface
}

var _ ServiceInterface = &Service{}

func New(usersService users.ServiceInterface) *Service {
	return &Service{
		UsersService: usersService,
	}
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserUid string
}

var key = []byte("helloWorld")

func (s *Service) Login(ctx context.Context, req LoginRequest) *LoginResponse {

	// validate req body
	if req.Email == "" {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing Parameters",
				ErrorMsg:  "Username parameters are missing",
			},
		}
	}

	if req.Password == "" {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Missing Parameters",
				ErrorMsg:  "Password parameter is missing",
			},
		}
	}

	// check if user exists
	user := s.UsersService.GetUserByEmail(ctx, userModels.GetUserByEmailRequest{
		Email: req.Email,
	})
	if user.BaseResponse != nil {

		fmt.Printf("Service.Login: failed to get user, user_email=%s err=%s", req.Email, user.BaseResponse.ErrorMsg)

		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.InvalidCredentialsError.ErrorType,
				ErrorMsg:   models.InvalidCredentialsError.ErrorMsg,
				ErrorStack: append(models.InvalidCredentialsError.ErrorStack, "Invalid email or passowrd used, please try again"),
			},
		}
	}

	// check if credentials are correct
	if req.Email != user.User.Email {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.InvalidCredentialsError.ErrorType,
				ErrorMsg:   models.InvalidCredentialsError.ErrorMsg,
				ErrorStack: append(models.InvalidCredentialsError.ErrorStack, "Invalid email or passowrd used, please try again"),
			},
		}
	}

	if !users.CheckPasswordHash(req.Password, user.User.Password) {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  models.InvalidCredentialsError.ErrorType,
				ErrorMsg:   models.InvalidCredentialsError.ErrorMsg,
				ErrorStack: append(models.InvalidCredentialsError.ErrorStack, "Invalid email or passowrd used, please try again"),
			},
		}
	}

	// generate JWT token
	resp := s.CreateToken(ctx, CreateTokenRequest{
		UserUid: user.User.Uid,
	})
	if resp.BaseResponse != nil {

		fmt.Printf("Service.Login: failed to create JWT token, err=%s", resp.BaseResponse.ErrorMsg)

		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: models.UnknownError.ErrorType,
				ErrorMsg:  models.UnknownError.ErrorMsg,
			},
		}
	}

	return &LoginResponse{
		Token: resp.Token,
	}
}

func (s *Service) CreateToken(ctx context.Context, req CreateTokenRequest) *CreateTokenResponse {

	token, err := generateJWTToken(req.UserUid)
	if err != nil {
		return &CreateTokenResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Failed to generate JWT",
				ErrorMsg:  err.Error(),
			},
		}
	}

	return &CreateTokenResponse{
		Token: token,
	}
}

func (s *Service) IsTokenValid(ctx context.Context, req IsTokenValidRequest) *IsTokenValidResponse {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(req.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return &IsTokenValidResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Failed to parse JWT token",
				ErrorMsg:  err.Error(),
			},
		}
	}

	return &IsTokenValidResponse{
		Valid:  token.Valid,
		Claims: claims,
	}
}

func generateJWTToken(uid string) (string, error) {

	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * 30),
			},
		},
		UserUid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
