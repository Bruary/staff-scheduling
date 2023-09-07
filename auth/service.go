package auth

import (
	"context"
	"time"

	models "github.com/Bruary/staff-scheduling/models"
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
	UserEmail string
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
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Unknown error",
				ErrorMsg:  "Login. Failed to get user information",
			},
		}
	}

	// if base response is not nil that means there is an error somewhere
	if user.BaseResponse != nil {
		return &LoginResponse{
			BaseResponse: user.BaseResponse,
		}
	}

	// check if credentials are correct
	if req.Email != user.User.Email || req.Password != user.User.Password {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType: "Invalid credentials",
				ErrorMsg:  "Login. Invalid user credentials",
			},
		}
	}

	// generate token
	resp := s.CreateToken(ctx, CreateTokenRequest{
		TokenType: JWT,
		Email:     req.Email,
	})
	if resp.BaseResponse != nil {
		return &LoginResponse{
			BaseResponse: &models.BaseResponse{
				ErrorType:  "Unknown error",
				ErrorMsg:   "Login. Error when trying to create a token.",
				ErrorStack: append(resp.BaseResponse.ErrorStack),
			},
		}
	}

	return &LoginResponse{
		Token: resp.Token,
	}
}

func (s *Service) CreateToken(ctx context.Context, req CreateTokenRequest) *CreateTokenResponse {

	if req.TokenType == JWT {

		token, err := generateJWTToken(req.Email)
		if err != nil {
			return &CreateTokenResponse{
				BaseResponse: &models.BaseResponse{
					ErrorType: "Failed to generate JWT",
					ErrorMsg:  err.Error(),
				},
			}
		}

		return &CreateTokenResponse{
			Token: Token{
				Token:     token,
				TokenType: req.TokenType,
			},
		}
	}

	return nil
}

func (s *Service) IsTokenValid(ctx context.Context, req IsTokenValidRequest) *IsTokenValidResponse {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(req.Token.Token, claims, func(t *jwt.Token) (interface{}, error) {
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
		Valid: token.Valid,
	}
}

func generateJWTToken(email string) (string, error) {

	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * 5),
			},
		},
		UserEmail: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
