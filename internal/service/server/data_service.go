package server

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"go.uber.org/zap"
)

type dataStorageRepo interface {
	CreatePairsData(ctx context.Context, data models.PairsData) (uint32, error)
	GetPairsData(ctx context.Context) ([]models.PairsData, error)
	UpdatePairsData(ctx context.Context, data models.PairsData) error
	DeletePairsData(ctx context.Context, id uint32) error
}

type DataService struct {
	store  dataStorageRepo
	logger *zap.SugaredLogger
}

func NewDataService(store dataStorageRepo, logger *zap.SugaredLogger) *DataService {
	return &DataService{
		store:  store,
		logger: logger.Named("APP"),
	}
}

func (s *DataService) PingDatabase() error {
	return nil
}
