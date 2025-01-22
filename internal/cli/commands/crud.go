package commands

import (
	"context"
	"strconv"

	"github.com/bbquite/go-pass-keeper/internal/models"
)

func (cm *CommandManager) createCommand(dataType models.DataTypeEnum, params CommandParams) error {
	paramsValidated := cm.validateParams(params)

	dataMarshal, err := marshalDataByType(dataType, paramsValidated)
	if err != nil {
		return err
	}

	err = cm.dataService.CreateData(context.Background(), &models.DataStoreFormat{
		DataType: dataType,
		DataInfo: string(dataMarshal),
		Meta:     paramsValidated["meta"].value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (cm *CommandManager) updateCommand(dataType models.DataTypeEnum, params CommandParams) error {
	paramsValidated := cm.validateParams(params)

	dataMarshal, err := marshalDataByType(dataType, paramsValidated)
	if err != nil {
		return err
	}

	dataID, err := strconv.ParseUint(paramsValidated["id"].value, 10, 32)
	if err != nil {
		return err
	}

	err = cm.dataService.UpdateData(context.Background(), &models.DataStoreFormat{
		ID:       uint32(dataID),
		DataType: dataType,
		DataInfo: string(dataMarshal),
		Meta:     paramsValidated["meta"].value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (cm *CommandManager) deleteCommand(dataType models.DataTypeEnum, params CommandParams) error {
	paramsValidated := cm.validateParams(params)

	dataID, err := strconv.ParseUint(paramsValidated["id"].value, 10, 32)
	if err != nil {
		return err
	}

	err = cm.dataService.DeleteData(context.Background(), uint32(dataID))
	if err != nil {
		return err
	}
	return nil
}
