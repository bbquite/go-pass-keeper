package handlers

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *GRPCHandler) CreatePairData(ctx context.Context, in *pb.CreatePairsDataRequest) (*pb.CreatePairsDataResponse, error) {
	response := pb.CreatePairsDataResponse{}
	pairData := models.PairData{
		Key:  in.Pairs.Key,
		Pwd:  in.Pairs.Pwd,
		Meta: in.Pairs.Meta,
	}

	resultPairData, err := h.dataService.CreatePairData(ctx, &pairData)
	h.logger.Info(resultPairData)

	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePairsDataResponse{
		Pairs: &pb.PairsData{
			Key:  in.Pairs.Key,
			Pwd:  in.Pairs.Pwd,
			Meta: in.Pairs.Meta,
		},
	}, nil
}

func (h *GRPCHandler) GetPairsDataList(ctx context.Context) (*pb.GetPairsDataResponse, error) {
	response := pb.GetPairsDataResponse{}

	resultPairsDataList, err := h.dataService.GetPairsDataList(ctx)
	h.logger.Info(resultPairsDataList)

	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	for _, pair := range resultPairsDataList {
		response.Pairs = append(response.Pairs, &pb.PairsData{
			Key:  pair.Key,
			Pwd:  pair.Pwd,
			Meta: pair.Meta,
		})
	}

	return &response, nil
}

func (h *GRPCHandler) UpdatePairData(ctx context.Context, in *pb.UpdatePairsDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}
	pairData := models.PairData{
		Key:  in.Pairs.Key,
		Pwd:  in.Pairs.Pwd,
		Meta: in.Pairs.Meta,
	}

	err := h.dataService.UpdatePairData(ctx, &pairData)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}

func (h *GRPCHandler) DeletePairData(ctx context.Context, in *pb.DeletePairsDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}
	pairID := in.GetId()

	err := h.dataService.DeletePairData(ctx, pairID)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}
