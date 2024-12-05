package command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/bbquite/go-pass-keeper/internal/models"
)

func (cm *CommandManager) createCommand(dataType models.DataTypeEnum, params CommandParams) error {

	var dataMarshal []byte
	var err error
	paramsValidated := cm.validateParams(params)

	switch dataType {
	case models.DataTypePAIR:
		dataMarshal, err = json.Marshal(&models.PairData{
			Key: paramsValidated["key"].value,
			Pwd: paramsValidated["pwd"].value,
		})
	case models.DataTypeTEXT:
		dataMarshal, err = json.Marshal(&models.TextData{
			Text: paramsValidated["text"].value,
		})
	case models.DataTypeBINARY:
		dataMarshal, err = json.Marshal(&models.BinaryData{
			FileName: paramsValidated["filename"].value,
		})
	case models.DataTypeCARD:
		dataMarshal, err = json.Marshal(&models.CardData{
			CardNum:   paramsValidated["number"].value,
			CardCvv:   paramsValidated["cvv"].value,
			CardExp:   paramsValidated["exp"].value,
			CardOwner: paramsValidated["owner"].value,
		})
	}

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

	var dataMarshal []byte
	var err error
	paramsValidated := cm.validateParams(params)

	switch dataType {
	case models.DataTypePAIR:
		dataMarshal, err = json.Marshal(&models.PairData{
			Key: paramsValidated["key"].value,
			Pwd: paramsValidated["pwd"].value,
		})
	case models.DataTypeTEXT:
		dataMarshal, err = json.Marshal(&models.TextData{
			Text: paramsValidated["text"].value,
		})
	case models.DataTypeBINARY:
		dataMarshal, err = json.Marshal(&models.BinaryData{
			FileName: paramsValidated["filename"].value,
		})
	case models.DataTypeCARD:
		dataMarshal, err = json.Marshal(&models.CardData{
			CardNum:   paramsValidated["number"].value,
			CardCvv:   paramsValidated["cvv"].value,
			CardExp:   paramsValidated["exp"].value,
			CardOwner: paramsValidated["owner"].value,
		})
	}

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

func (cm *CommandManager) deleteCommand(params CommandParams) error {

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
