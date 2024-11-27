package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"github.com/bbquite/go-pass-keeper/internal/service"
	"github.com/bbquite/go-pass-keeper/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	pb.UnimplementedPassKeeperServiceServer
	cfg         *ServerConfig
	service     *service.AppService
	authService *service.AuthService
	DBStorage   *storage.DBStorage
	logger      *zap.SugaredLogger
}

func NewGRPCServer(cfg *ServerConfig, logger *zap.SugaredLogger) (*gRPCServer, error) {

	dbStorage, err := storage.NewDBStorage(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	appService := service.NewAppService(dbStorage, logger)
	authService := service.NewAuthService(dbStorage, cfg.JWTSecret, logger)

	return &gRPCServer{
		cfg:         cfg,
		authService: authService,
		service:     appService,
		DBStorage:   dbStorage,
		logger:      logger,
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
	pb.RegisterPassKeeperServiceServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("error occured while running gRPC server: %v", err)
		}
	}()

	sig := <-signalCh
	s.logger.Info("Received signal: %v\n", sig)

	grpcServer.GracefulStop()
	s.DBStorage.DB.Close()

	s.logger.Info("Server shutdown gracefully")
}

func (s *gRPCServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	response := pb.RegisterUserResponse{}
	userData := models.UserRegisterData{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
	}

	token, err := s.authService.RegisterUser(ctx, &userData)

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			s.logger.Info(err)
			return &response, status.Error(codes.AlreadyExists, "user already exists")
		}

		s.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterUserResponse{
		Token: token.Token,
	}, nil
}
