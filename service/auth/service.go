package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/api/dto"
	"github.com/xarest/gobs-collection/lib/logger"
	"github.com/xarest/gobs-collection/schema"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log        logger.ILogger
	jwtToken   *JwtToken
	repository *UserRepository
}

func (s *Auth) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			&JwtToken{},
			&UserRepository{},
		},
	}, nil
}

func (s *Auth) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&s.log, &s.jwtToken, &s.repository)
}

func (s *Auth) Register(ctx context.Context, creds dto.Credentials) (*schema.User, error) {
	user, err := s.repository.GetUserForSignIn(ctx, creds.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, echo.NewHTTPError(http.StatusConflict, "user already exists")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(creds.Pass), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "password is too short or too long")
	}
	if err != nil {
		return nil, err
	}
	return s.repository.AddUser(ctx,
		schema.UserCreateParams{
			Email:    creds.Email,
			Password: string(hashPass),
		})
}

func (s *Auth) SignIn(ctx context.Context, creds dto.Credentials) (*dto.RespToken, error) {
	// Get the expected password from our in memory map
	user, err := s.repository.GetUserForSignIn(ctx, creds.Email)
	if nil != err {
		return nil, err
	}
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Pass)); nil != err {
		return nil, err
	}

	userId := user.ID.String()
	tokenStr, expired, err := s.jwtToken.ComposeToken(userId)
	if nil != err {
		return nil, err
	}

	return &dto.RespToken{
		Token:       tokenStr,
		UserID:      userId,
		AccessToken: "",
		ExpiresAt:   expired,
	}, nil
}

func (s *Auth) Logout(ctx context.Context) (any, error) {
	return nil, nil
}

func (a *Auth) RefreshToken(ctx context.Context, tokenStr string) (*dto.RespToken, error) {
	tokenStr, expired, err := a.jwtToken.RefreshToken(tokenStr)
	if nil != err {
		return nil, err
	}
	return &dto.RespToken{
		Token:     tokenStr,
		ExpiresAt: expired,
	}, nil
}

var _ gobs.IServiceInit = (*Auth)(nil)
var _ gobs.IServiceSetup = (*Auth)(nil)
