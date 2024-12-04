package server

import (
	"context"
	"errors"

	"github.com/bbquite/go-pass-keeper/internal/models"
	"go.uber.org/zap"
)

var (
	ErrUniqueViolation = errors.New("record already exists")
)

type dataStorageRepo interface {
	CreateData(ctx context.Context, data *models.DataStoreFormat) (models.DataStoreFormat, error)
	GetDataList(ctx context.Context) ([]models.DataStoreFormat, error)
	UpdateData(ctx context.Context, data *models.DataStoreFormat) error
	DeleteData(ctx context.Context, pairID uint32) error
}

type DataService struct {
	store  dataStorageRepo
	logger *zap.SugaredLogger
}

func NewDataService(store dataStorageRepo, logger *zap.SugaredLogger) *DataService {
	return &DataService{
		store:  store,
		logger: logger.Named("DATA"),
	}
}

func (service *DataService) CreateData(ctx context.Context, data *models.DataStoreFormat) (models.DataStoreFormat, error) {
	var resultData models.DataStoreFormat

	resultData, err := service.store.CreateData(ctx, data)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (service *DataService) GetDataList(ctx context.Context) ([]models.DataStoreFormat, error) {
	var resultDataList []models.DataStoreFormat

	resultDataList, err := service.store.GetDataList(ctx)
	if err != nil {
		return resultDataList, err
	}

	return resultDataList, nil
}

func (service *DataService) UpdateData(ctx context.Context, pairData *models.DataStoreFormat) error {
	err := service.store.UpdateData(ctx, pairData)
	if err != nil {
		return err
	}

	return nil
}

func (service *DataService) DeleteData(ctx context.Context, pairID uint32) error {
	err := service.store.DeleteData(ctx, pairID)
	if err != nil {
		return err
	}

	return nil
}
