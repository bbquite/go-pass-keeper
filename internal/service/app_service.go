package service

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"go.uber.org/zap"
)

type appStorageRepo interface {
	CreatePairsData(ctx context.Context, data models.PairsData) (uint32, error)
	GetPairsData(ctx context.Context) ([]models.PairsData, error)
	UpdatePairsData(ctx context.Context, data models.PairsData) error
	DeletePairsData(ctx context.Context, id uint32) error
}

type AppService struct {
	store  appStorageRepo
	logger *zap.SugaredLogger
}

func NewAppService(store appStorageRepo, logger *zap.SugaredLogger) *AppService {
	return &AppService{
		store:  store,
		logger: logger.Named("APP"),
	}
}

func (s *AppService) PingDatabase() error {
	return nil
}
