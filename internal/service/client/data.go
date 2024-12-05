package client

import (
	"context"
	"fmt"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type clientDataStorageRepo interface {
	SetUserID(userID *uint32) error
	SetToken(token *jwttoken.JWT) error
	GetToken() string
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

func (service *ClientDataService) CreateData(ctx context.Context, data *models.DataStoreFormat) error {

	token := service.store.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+service.store.GetToken())

	_, err := service.grpcClient.PBService.CreateData(ctx, &pb.CreateDataRequest{
		Data: &pb.DataItem{
			DataType: pb.DataTypeEnum(pb.DataTypeEnum_value[string(data.DataType)]),
			DataInfo: data.DataInfo,
			Meta:     data.Meta,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *ClientDataService) UpdateData(ctx context.Context, data *models.DataStoreFormat) error {

	token := service.store.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+service.store.GetToken())

	_, err := service.grpcClient.PBService.UpdateData(ctx, &pb.UpdateDataRequest{
		Data: &pb.DataItem{
			Id:       data.ID,
			DataType: pb.DataTypeEnum(pb.DataTypeEnum_value[string(data.DataType)]),
			DataInfo: data.DataInfo,
			Meta:     data.Meta,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *ClientDataService) DeleteData(ctx context.Context, dataID uint32) error {

	token := service.store.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+service.store.GetToken())

	_, err := service.grpcClient.PBService.DeleteData(ctx, &pb.DeleteDataRequest{
		Id: dataID,
	})

	if err != nil {
		return err
	}

	return nil
}
