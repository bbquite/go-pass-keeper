package client

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/models"
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

func (service *ClientDataService) CreatePairData(ctx context.Context, data *models.PairData) error {
	return nil
}

func (service *ClientDataService) UpdatePairData(ctx context.Context, data *models.PairData) error {
	return nil
}

func (service *ClientDataService) DeletePairData(ctx context.Context, data *models.PairData) error {
	return nil
}

func (service *ClientDataService) CreateTextData(ctx context.Context, data *models.TextData) error {
	return nil
}

func (service *ClientDataService) UpdateTextData(ctx context.Context, data *models.TextData) error {
	return nil
}

func (service *ClientDataService) DeleteTextData(ctx context.Context, data *models.TextData) error {
	return nil
}

func (service *ClientDataService) CreateBinaryData(ctx context.Context, data *models.BinaryData) error {
	return nil
}

func (service *ClientDataService) UpdateBinaryData(ctx context.Context, data *models.BinaryData) error {
	return nil
}

func (service *ClientDataService) DeleteBinaryData(ctx context.Context, data *models.BinaryData) error {
	return nil
}

func (service *ClientDataService) CreateCardData(ctx context.Context, data *models.CardData) error {
	return nil
}

func (service *ClientDataService) UpdateCardData(ctx context.Context, data *models.CardData) error {
	return nil
}

func (service *ClientDataService) DeleteCardData(ctx context.Context, data *models.CardData) error {
	return nil
}
