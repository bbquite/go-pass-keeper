package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bbquite/go-pass-keeper/internal/config"
	"github.com/bbquite/go-pass-keeper/internal/handlers"
	"github.com/bbquite/go-pass-keeper/internal/storage/postgres"

	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	cfg       *config.ServerConfig
	dbStorage *postgres.DBStorage
	handler   *handlers.GRPCHandler
	logger    *zap.SugaredLogger
}

func NewGRPCServer(cfg *config.ServerConfig, logger *zap.SugaredLogger) (*gRPCServer, error) {

	dbStorage, err := postgres.NewDBStorage(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	handler := handlers.NewGRPCHandler(cfg.JWTSecret, dbStorage, logger)

	return &gRPCServer{
		cfg:       cfg,
		handler:   handler,
		dbStorage: dbStorage,
		logger:    logger.Named("SERVER"),
	}, nil
}

func (s *gRPCServer) RunGRPCServer() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	listen, err := net.Listen("tcp", s.cfg.Host)
	if err != nil {
		log.Fatalf("error occured while running gRPC server: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPassKeeperServiceServer(grpcServer, s.handler)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("error occured while running gRPC server: %v", err)
		}
	}()

	sig := <-signalCh
	s.logger.Info("Received signal: %v\n", sig)

	grpcServer.GracefulStop()
	s.dbStorage.DB.Close()

	s.logger.Info("Server shutdown gracefully")
}
