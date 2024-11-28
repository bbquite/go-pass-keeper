package handlers

import (
	"github.com/bbquite/go-pass-keeper/internal/app/server"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"github.com/bbquite/go-pass-keeper/internal/service"
	"github.com/bbquite/go-pass-keeper/internal/storage"
	"go.uber.org/zap"
)

type GRPCHandler struct {
	pb.UnimplementedPassKeeperServiceServer
	appService  *service.AppService
	authService *service.AuthService
	dbStorage   *storage.DBStorage
	logger      *zap.SugaredLogger
}

func NewGRPCHandler(cfg *server.ServerConfig, dbStorage *storage.DBStorage, logger *zap.SugaredLogger) *GRPCHandler {

	appService := service.NewAppService(dbStorage, logger)
	authService := service.NewAuthService(dbStorage, cfg.JWTSecret, logger)

	return &GRPCHandler{
		appService:  appService,
		authService: authService,
		dbStorage:   dbStorage,
		logger:      logger,
	}
}
