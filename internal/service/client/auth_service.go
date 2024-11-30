package client

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
)

//var (
//	ErrUserAlreadyExists  = errors.New("user already exists")
//	ErrIncorrectLoginData = errors.New("incorrect login or password")
//)

type clientAuthStorageRepo interface {
	SetUserID(userID *uint32) error
	SetToken(token *jwttoken.JWT) error
}

type ClientAuthService struct {
	grpcClient *client.GRPCClient
	store      clientAuthStorageRepo
	logger     *zap.SugaredLogger
}

func NewClientAuthService(grpcClient *client.GRPCClient, store clientAuthStorageRepo, logger *zap.SugaredLogger) *ClientAuthService {

	return &ClientAuthService{
		grpcClient: grpcClient,
		store:      store,
		logger:     logger.Named("CLIENT AUTH"),
	}
}

func (service *ClientAuthService) RegisterUser(ctx context.Context, userData *models.UserRegisterData) error {
	//var token jwttoken.JWT

	resp, err := service.grpcClient.Client.RegisterUser(ctx, &pb.RegisterUserRequest{
		Username: userData.Username,
		Password: userData.Password,
		Email:    userData.Email,
	})

	if err != nil {
		service.logger.Error(err)
		return err
	}

	service.logger.Debug(resp)

	return nil
}

//func (service *ClientAuthService) AuthUser(ctx context.Context, userData *models.UserLoginData) error {
//	var token jwttoken.JWT
//
//	return nil
//}
