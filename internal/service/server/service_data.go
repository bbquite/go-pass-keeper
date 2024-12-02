package server

import (
	"context"
	"errors"

	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var (
	ErrUniqueViolation = errors.New("record already exists")
)

type dataStorageRepo interface {
	CreatePairData(ctx context.Context, data *models.PairData) (models.PairData, error)
	GetPairsDataList(ctx context.Context) ([]models.PairData, error)
	UpdatePairData(ctx context.Context, data *models.PairData) error
	DeletePairData(ctx context.Context, pairID uint32) error
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

func (service *DataService) CreatePairData(ctx context.Context, pairData *models.PairData) (models.PairData, error) {
	var resultPairData models.PairData

	resultPairData, err := service.store.CreatePairData(ctx, pairData)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if pgErr.Code == pgerrcode.UniqueViolation {
			return resultPairData, ErrUniqueViolation
		}

		return resultPairData, err
	}

	return resultPairData, nil
}

func (service *DataService) GetPairsDataList(ctx context.Context) ([]models.PairData, error) {
	var resultPairsDataList []models.PairData

	resultPairsDataList, err := service.store.GetPairsDataList(ctx)
	if err != nil {
		return resultPairsDataList, err
	}

	return resultPairsDataList, nil
}

func (service *DataService) UpdatePairData(ctx context.Context, pairData *models.PairData) error {
	err := service.store.UpdatePairData(ctx, pairData)
	if err != nil {
		return err
	}

	return nil
}

func (service *DataService) DeletePairData(ctx context.Context, pairID uint32) error {
	err := service.store.DeletePairData(ctx, pairID)
	if err != nil {
		return err
	}

	return nil
}
