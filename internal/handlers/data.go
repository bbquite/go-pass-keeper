package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	formatTimeLayout = "2006-01-02 15:04:05"
)

func (h *GRPCHandler) CreateData(ctx context.Context, in *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	response := pb.CreateDataResponse{}
	data := models.DataStoreFormat{
		DataType: models.DataTypeEnum(in.Data.GetDataType().String()),
		DataInfo: in.Data.GetDataInfo(),
		Meta:     in.Data.GetMeta(),
	}

	resultData, err := h.dataService.CreateData(ctx, &data)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateDataResponse{
		Data: &pb.DataItem{
			Id:         resultData.ID,
			DataType:   pb.DataTypeEnum(pb.DataTypeEnum_value[string(resultData.DataType)]),
			DataInfo:   resultData.DataInfo,
			Meta:       resultData.Meta,
			UploadedAt: resultData.UploadedAt.Format(formatTimeLayout),
		},
	}, nil
}

func (h *GRPCHandler) GetDataList(ctx context.Context, in *pb.Empty) (*pb.GetDataResponse, error) {
	response := pb.GetDataResponse{}

	resultDataList, err := h.dataService.GetDataList(ctx)
	if err != nil {
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}

	for _, item := range resultDataList {
		response.DataList = append(response.DataList, &pb.DataItem{
			Id:         item.ID,
			DataType:   pb.DataTypeEnum(pb.DataTypeEnum_value[string(item.DataType)]),
			DataInfo:   item.DataInfo,
			Meta:       item.Meta,
			UploadedAt: item.UploadedAt.Format(formatTimeLayout),
		})
	}

	return &response, nil
}

func (h *GRPCHandler) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}

	data := models.DataStoreFormat{
		ID:       in.Data.GetId(),
		DataType: models.DataTypeEnum(in.Data.GetDataType().String()),
		DataInfo: in.Data.GetDataInfo(),
		Meta:     in.Data.GetMeta(),
	}

	err := h.dataService.UpdateData(ctx, &data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &response, status.Error(codes.NotFound, err.Error())
		}
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}

func (h *GRPCHandler) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*pb.Empty, error) {
	response := pb.Empty{}
	dataID := in.GetId()

	log.Print(dataID)

	err := h.dataService.DeleteData(ctx, dataID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &response, status.Error(codes.NotFound, err.Error())
		}
		h.logger.Error(err)
		return &response, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}
