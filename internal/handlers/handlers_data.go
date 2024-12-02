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

func (h *GRPCHandler) CreatePairData(ctx context.Context, in *pb.CreatePairDataRequest) (*pb.CreatePairDataResponse, error) {
	response := pb.CreatePairDataResponse{}

	if in.Pair == nil {
		return &response, status.Error(codes.InvalidArgument, "")
	} else {
		if in.Pair.Key == "" || in.Pair.Pwd == "" {
			return &response, status.Error(codes.InvalidArgument, "")
		}
	}

	pairData := models.PairData{
		Key:  in.Pair.Key,
		Pwd:  in.Pair.Pwd,
		Meta: in.Pair.Meta,
	}

	resultPairData, err := h.dataService.CreatePairData(ctx, &pairData)
	if err != nil {
		if errors.Is(err, serverServices.ErrUniqueViolation) {
			return &response, status.Error(codes.InvalidArgument, err.Error())
		}
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePairDataResponse{
		Pair: &pb.PairData{
			Id:         resultPairData.ID,
			Key:        resultPairData.Key,
			Pwd:        resultPairData.Pwd,
			Meta:       resultPairData.Meta,
			UploadedAt: resultPairData.UploadedAt.String(),
		},
	}, nil
}

func (h *GRPCHandler) GetPairsDataList(ctx context.Context, in *pb.Empty) (*pb.GetPairsDataResponse, error) {
	response := pb.GetPairsDataResponse{}

	resultPairsDataList, err := h.dataService.GetPairsDataList(ctx)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	for _, pair := range resultPairsDataList {
		response.Pairs = append(response.Pairs, &pb.PairData{
			Id:         pair.ID,
			Key:        pair.Key,
			Pwd:        pair.Pwd,
			Meta:       pair.Meta,
			UploadedAt: pair.UploadedAt.String(),
		})
	}

	return &response, nil
}

func (h *GRPCHandler) UpdatePairData(ctx context.Context, in *pb.UpdatePairDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}

	if in.Pair == nil {
		return &response, status.Error(codes.InvalidArgument, "")
	} else {
		if in.Pair.Id == 0 || in.Pair.Key == "" || in.Pair.Pwd == "" {
			return &response, status.Error(codes.InvalidArgument, "")
		}
	}

	pairData := models.PairData{
		ID:   in.Pair.Id,
		Key:  in.Pair.Key,
		Pwd:  in.Pair.Pwd,
		Meta: in.Pair.Meta,
	}

	err := h.dataService.UpdatePairData(ctx, &pairData)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}

func (h *GRPCHandler) DeletePairData(ctx context.Context, in *pb.DeletePairDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}
	pairID := in.GetId()

	err := h.dataService.DeletePairData(ctx, pairID)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}
