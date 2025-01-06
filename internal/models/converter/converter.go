package converter

import (
	"github.com/bbquite/go-pass-keeper/internal/models"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
)

func DataStoreFormatToProtoFormat(data *models.DataStoreFormat) *pb.DataItem {
	dataConverted := &pb.DataItem{
		Id:       data.ID,
		DataType: pb.DataTypeEnum(pb.DataTypeEnum_value[string(data.DataType)]),
		DataInfo: data.DataInfo,
		Meta:     data.Meta,
	}
	return dataConverted
}
