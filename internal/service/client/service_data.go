package client

import (
	"github.com/bbquite/go-pass-keeper/internal/app/client"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
)

type clientDataStorageRepo interface {
	SetUserID(userID *uint32) error
	SetToken(token *jwttoken.JWT) error
	Debug() ([]byte, error)
}

type ClientDataService struct {
	grpcClient *client.GRPCClient
	store      clientDataStorageRepo
	logger     *zap.SugaredLogger
}

func NewClientDataService(grpcClient *client.GRPCClient, store clientDataStorageRepo, logger *zap.SugaredLogger) *ClientDataService {

	return &ClientDataService{
		grpcClient: grpcClient,
		store:      store,
		logger:     logger.Named("CLIENT DATA"),
	}
}

func (service *ClientDataService) Debug() error {
	debug, err := service.store.Debug()
	if err != nil {
		return err
	}
	service.logger.Debugf("%s", debug)
	return nil
}
