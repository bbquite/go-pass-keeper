package service

import (
	"go.uber.org/zap"
)

type StorageRepo interface {
	//Ping() error
}

type AppService struct {
	store  StorageRepo
	logger *zap.SugaredLogger
}

func NewAppService(store StorageRepo, logger *zap.SugaredLogger) *AppService {
	return &AppService{
		store:  store,
		logger: logger.Named("SERVICE"),
	}
}

//func (s *AppService) PingDatabase() error {
//	err := s.store.Ping()
//	if err != nil {
//		return err
//	}
//	return nil
//}
