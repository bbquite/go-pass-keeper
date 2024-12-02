package handlers

import (
	"context"
	"errors"
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	serverServices "github.com/bbquite/go-pass-keeper/internal/service/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	response := pb.RegisterUserResponse{}
	userData := models.UserRegisterData{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
	}

	token, err := h.authService.RegisterUser(ctx, &userData)

	if err != nil {
		if errors.Is(err, serverServices.ErrUserAlreadyExists) {
			h.logger.Info(err)
			return &response, status.Error(codes.AlreadyExists, "user already exists")
		}

		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterUserResponse{
		Token: token.Token,
	}, nil
}

func (h *GRPCHandler) AuthUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	response := pb.AuthUserResponse{}
	userData := models.UserLoginData{
		Username: in.Username,
		Password: in.Password,
	}

	token, err := h.authService.AuthUser(ctx, &userData)

	if err != nil {
		if errors.Is(err, serverServices.ErrIncorrectLoginData) {
			h.logger.Info(err)
			return &response, status.Error(codes.Unauthenticated, "incorrect login or password")
		}

		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.AuthUserResponse{
		Token: token.Token,
	}, nil
}
