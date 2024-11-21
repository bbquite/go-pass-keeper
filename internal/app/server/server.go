package app

import (
	"encoding/json"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/service"
	"github.com/bbquite/go-pass-keeper/internal/storage"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type ServerConfig struct {
	Host        string
	DatabaseURI string
}

type gRPCServer struct {
	//pb.UnimplementedMetricServiceServer
	cfg     *ServerConfig
	service *service.AppService
	logger  *zap.SugaredLogger
}

func initServerLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	defer logger.Sync()
	return sugar, nil
}

func NewGRPCServer(cfg *ServerConfig) (*gRPCServer, error) {

	logger, err := initServerLogger()
	if err != nil {
		return nil, err
	}

	dbStorage, err := storage.NewDBStorage(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	defer dbStorage.DB.Close()

	appService := service.NewAppService(dbStorage, logger)

	return &gRPCServer{
		cfg:     cfg,
		service: appService,
		logger:  logger,
	}, nil
}

func (s *gRPCServer) RunGRPCServer() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//listen, err := net.Listen("tcp", s.cfg.Host)
	//if err != nil {
	//	log.Fatalf("error occured while running gRPC server: %v", err)
	//}
	//
	//grpcServer := grpc.NewServer()
	//reflection.Register(grpcServer)
	//pb.RegisterMetricServiceServer(grpcServer, s)
	//
	//go func() {
	//	if err := grpcServer.Serve(listen); err != nil {
	//		log.Fatalf("error occured while running gRPC server: %v", err)
	//	}
	//}()
	//

	jsonConfig, _ := json.Marshal(s.cfg)
	s.logger.Infof("Server run with config: %s", jsonConfig)

	sig := <-signalCh
	s.logger.Info("Received signal: %v\n", sig)

	//grpcServer.GracefulStop()

	s.logger.Info("Server shutdown gracefully")
}
