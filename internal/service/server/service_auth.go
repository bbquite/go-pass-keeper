package server

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrIncorrectLoginData = errors.New("incorrect login or password")
)

type authStorageRepo interface {
	CreateAccount(ctx context.Context, username string, password string, email string) (uint32, error)
	GetAccountByUsername(ctx context.Context, username string) (models.Account, error)
	GetAccountByLoginData(ctx context.Context, username string, password string) (models.Account, error)
}

type AuthService struct {
	store      authStorageRepo
	logger     *zap.SugaredLogger
	jwtManager *jwttoken.JWTManager
}

func NewAuthService(store authStorageRepo, jwtManager *jwttoken.JWTManager, logger *zap.SugaredLogger) *AuthService {
	return &AuthService{
		store:      store,
		logger:     logger.Named("AUTH"),
		jwtManager: jwtManager,
	}
}

func (service *AuthService) RegisterUser(ctx context.Context, userData *models.UserRegisterData) (jwttoken.JWT, error) {
	var token jwttoken.JWT

	_, err := service.store.GetAccountByUsername(ctx, userData.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			userID, err := service.store.CreateAccount(
				ctx, userData.Username, utils.GenerateSHAString(userData.Password), userData.Email)
			if err != nil {
				return token, err
			}

			tokenString, err := service.jwtManager.CreateAccessToken(userID)

			if err != nil {
				return token, err
			}

			token.Token = tokenString
			return token, nil
		}
		return token, err
	}
	return token, ErrUserAlreadyExists
}

func (service *AuthService) AuthUser(ctx context.Context, userData *models.UserLoginData) (jwttoken.JWT, error) {
	var token jwttoken.JWT

	shaInputPassword := utils.GenerateSHAString(userData.Password)

	account, err := service.store.GetAccountByLoginData(ctx, userData.Username, shaInputPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return token, ErrIncorrectLoginData
		}
		return token, err
	}

	tokenString, err := service.jwtManager.CreateAccessToken(account.ID)
	if err != nil {
		return token, err
	}

	token.Token = tokenString
	return token, nil
}
