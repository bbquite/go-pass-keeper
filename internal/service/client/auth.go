package client

import (
	"context"
	"errors"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrIncorrectLoginData = errors.New("incorrect login or password")
)

type clientAuthStorageRepo interface {
	SetUserID(userID *uint32) error
	SetToken(token *jwttoken.JWT) error
	Debug() ([]byte, error)
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

func (service *ClientAuthService) RegisterUser(ctx context.Context, userData *models.UserAccountData) error {
	var token jwttoken.JWT

	resp, err := service.grpcClient.PBService.RegisterUser(ctx, &pb.RegisterUserRequest{
		Username: userData.Username,
		Password: userData.Password,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	token.Token = resp.GetToken()
	err = service.store.SetToken(&token)
	if err != nil {
		return err
	}

	return nil
}

func (service *ClientAuthService) AuthUser(ctx context.Context, userData *models.UserAccountData) error {
	var token jwttoken.JWT

	resp, err := service.grpcClient.PBService.AuthUser(ctx, &pb.AuthUserRequest{
		Username: userData.Username,
		Password: userData.Password,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				return ErrIncorrectLoginData
			}
		}
		return err
	}

	token.Token = resp.GetToken()
	err = service.store.SetToken(&token)
	if err != nil {
		return err
	}

	return nil
}
