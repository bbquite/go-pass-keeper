package handlers

import (
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	serverServices "github.com/bbquite/go-pass-keeper/internal/service/server"
	"github.com/bbquite/go-pass-keeper/internal/storage/postgres"
	"go.uber.org/zap"
)

type GRPCHandler struct {
	pb.UnimplementedPassKeeperServiceServer
	dataService *serverServices.DataService
	authService *serverServices.AuthService
	dbStorage   *postgres.DBStorage
	logger      *zap.SugaredLogger
}

func NewGRPCHandler(jwtSecret string, dbStorage *postgres.DBStorage, logger *zap.SugaredLogger) *GRPCHandler {

	dataService := serverServices.NewDataService(dbStorage, logger)
	authService := serverServices.NewAuthService(dbStorage, jwtSecret, logger)

	return &GRPCHandler{
		dataService: dataService,
		authService: authService,
		dbStorage:   dbStorage,
		logger:      logger,
	}
}
