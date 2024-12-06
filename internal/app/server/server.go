package server

import (
	"fmt"
	encryptor "github.com/bbquite/go-pass-keeper/internal/encryption"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbquite/go-pass-keeper/internal/interceptors"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"

	"github.com/bbquite/go-pass-keeper/internal/config"
	"github.com/bbquite/go-pass-keeper/internal/handlers"
	"github.com/bbquite/go-pass-keeper/internal/storage/postgres"

	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	cfg           *config.ServerConfig
	dbStorage     *postgres.DBStorage
	handler       *handlers.GRPCHandler
	interceptors  []grpc.UnaryServerInterceptor
	jwtManager    *jwttoken.JWTManager
	noAuthMethods []string
	logger        *zap.SugaredLogger
}

func NewGRPCServer(cfg *config.ServerConfig, logger *zap.SugaredLogger) (*gRPCServer, error) {

	dbStorage, err := postgres.NewDBStorage(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	jwtManager := jwttoken.NewJWTTokenManager(time.Hour*3, cfg.JWTSecret)
	encryptorManager := encryptor.NewEncryptor([]byte(cfg.CryptoKey))
	handler := handlers.NewGRPCHandler(jwtManager, encryptorManager, dbStorage, logger)

	noAuthMethods := []string{
		"/internal.proto.PassKeeperService/RegisterUser",
		"/internal.proto.PassKeeperService/AuthUser",
	}

	serverInit := &gRPCServer{
		cfg:        cfg,
		handler:    handler,
		dbStorage:  dbStorage,
		jwtManager: jwtManager,

		noAuthMethods: noAuthMethods,
		logger:        logger.Named("SERVER"),
	}

	err = serverInit.loadServerInterceptors()
	if err != nil {
		log.Fatal(err)
	}

	return serverInit, nil
}

func (s *gRPCServer) loadServerInterceptors() error {
	var grpcServerInterceptors []grpc.UnaryServerInterceptor
	grpcServerInterceptors = append(grpcServerInterceptors, interceptors.NewAuthInterceptor(s.jwtManager, s.noAuthMethods).Unary())
	s.interceptors = grpcServerInterceptors
	return nil
}

func (s *gRPCServer) RunGRPCServer() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	listen, err := net.Listen("tcp", s.cfg.Host)
	if err != nil {
		log.Fatalf("error occured while running gRPC server: %v", err)
	}

	grpcCredos, err := credentials.NewServerTLSFromFile(s.cfg.GetServerCrtPath(), s.cfg.GetServerKeyPath())
	if err != nil {
		log.Fatalf("failed to load TLS certificates: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(grpcCredos),
		grpc.ChainUnaryInterceptor(s.interceptors...))
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
