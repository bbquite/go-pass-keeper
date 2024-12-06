package handlers

import (
	encryptor "github.com/bbquite/go-pass-keeper/internal/encryption"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	serverServices "github.com/bbquite/go-pass-keeper/internal/service/server"
	"github.com/bbquite/go-pass-keeper/internal/storage/postgres"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"go.uber.org/zap"
)

type GRPCHandler struct {
	pb.UnimplementedPassKeeperServiceServer
	dataService *serverServices.DataService
	authService *serverServices.AuthService
	logger      *zap.SugaredLogger
}

func NewGRPCHandler(jwtManager *jwttoken.JWTManager, encryptorManager *encryptor.Encryptor, dbStorage *postgres.DBStorage, logger *zap.SugaredLogger) *GRPCHandler {

	dataService := serverServices.NewDataService(dbStorage, encryptorManager, logger)
	authService := serverServices.NewAuthService(dbStorage, jwtManager, logger)

	return &GRPCHandler{
		dataService: dataService,
		authService: authService,
		logger:      logger,
	}
}
