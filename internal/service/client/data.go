package client

import (
	"context"
	"encoding/json"
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

	AddPairs(data models.PairData) error
	GetPairs() ([]models.PairData, error)

	AddTexts(data models.TextData) error
	GetTexts() ([]models.TextData, error)

	AddBinaries(data models.BinaryData) error
	GetBinary() ([]models.BinaryData, error)

	AddCards(data models.CardData) error
	GetCards() ([]models.CardData, error)

	Debug() ([]byte, error)
	ClearStorage()
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

func (service *ClientDataService) GetDataByID(ctx context.Context, dataID uint32) (models.DataStoreFormat, error) {

	var resultData models.DataStoreFormat
	token := service.store.GetToken()
	if token == "" {
		return resultData, fmt.Errorf("authorization only")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+service.store.GetToken())

	response, err := service.grpcClient.PBService.GetDataByID(ctx, &pb.GetDataByIDRequest{Id: dataID})
	if err != nil {
		return resultData, err
	}

	resultData.ID = response.Data.GetId()
	resultData.Meta = response.Data.GetMeta()
	resultData.DataInfo = response.Data.GetDataInfo()
	resultData.DataType = models.DataTypeEnum(response.Data.DataType)

	return resultData, nil
}

func (service *ClientDataService) GetData(ctx context.Context) error {

	token := service.store.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+service.store.GetToken())

	response, err := service.grpcClient.PBService.GetDataList(ctx, &pb.Empty{})
	if err != nil {
		return err
	}

	service.store.ClearStorage()
	dataItems := response.GetDataList()

	for _, item := range dataItems {
		switch item.DataType {
		case pb.DataTypeEnum(pb.DataTypeEnum_value[string(models.DataTypePAIR)]):

			m := models.PairData{}
			d := []byte(item.DataInfo)
			err := json.Unmarshal(d, &m)
			if err != nil {
				return err
			}
			m.ID = item.Id
			m.Meta = item.Meta
			service.store.AddPairs(m)

		case pb.DataTypeEnum(pb.DataTypeEnum_value[string(models.DataTypeTEXT)]):

			m := models.TextData{}
			d := []byte(item.DataInfo)
			err := json.Unmarshal(d, &m)
			if err != nil {
				return err
			}
			m.ID = item.Id
			m.Meta = item.Meta
			service.store.AddTexts(m)

		case pb.DataTypeEnum(pb.DataTypeEnum_value[string(models.DataTypeBINARY)]):

			m := models.BinaryData{}
			d := []byte(item.DataInfo)
			err := json.Unmarshal(d, &m)
			if err != nil {
				return err
			}
			m.ID = item.Id
			m.Meta = item.Meta
			service.store.AddBinaries(m)

		case pb.DataTypeEnum(pb.DataTypeEnum_value[string(models.DataTypeCARD)]):

			m := models.CardData{}
			d := []byte(item.DataInfo)
			err := json.Unmarshal(d, &m)
			if err != nil {
				return err
			}
			m.ID = item.Id
			m.Meta = item.Meta
			service.store.AddCards(m)
		}
	}

	return nil
}
