package commands

import (
	"context"
	"encoding/json"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	"strconv"
)

func (cm *CommandManager) exportCommand(dataType models.DataTypeEnum, params CommandParams) error {

	err := cm.dataService.GetData(context.Background())
	if err != nil {
		return err
	}

	switch dataType {
	case models.DataTypePAIR:
		return cm.exportPairData()

	case models.DataTypeTEXT:
		return cm.exportTextData()

	case models.DataTypeCARD:
		return cm.exportCardData()

	case models.DataTypeBINARY:
		paramsValidated := cm.validateParams(params)
		if _, ok := paramsValidated["id"]; ok {
			dataID, err := strconv.ParseUint(paramsValidated["id"].value, 10, 32)
			if err != nil {
				return err
			}
			return cm.exportBinaryData(uint32(dataID))
		}
	}

	return nil
}

func (cm *CommandManager) exportPairData() error {
	jsOut, err := json.MarshalIndent(cm.localStorage.PairsList, "", "	")
	if err != nil {
		return err
	}
	err = utils.SaveFile(cm.pairExportFilePath, jsOut)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CommandManager) exportTextData() error {
	jsOut, err := json.MarshalIndent(cm.localStorage.TextsList, "", "	")
	if err != nil {
		return err
	}
	err = utils.SaveFile(cm.textExportFilePath, jsOut)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CommandManager) exportCardData() error {
	jsOut, err := json.MarshalIndent(cm.localStorage.CardsList, "", "	")
	if err != nil {
		return err
	}
	err = utils.SaveFile(cm.cardExportFilePath, jsOut)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CommandManager) exportBinaryData(dataID uint32) error {
	dataItem, err := cm.dataService.GetDataByID(context.Background(), dataID)
	if err != nil {
		return err
	}

	m := models.BinaryData{}
	d := []byte(dataItem.DataInfo)
	err = json.Unmarshal(d, &m)
	if err != nil {
		return err
	}

	err = utils.SaveFile("./"+m.FileName, m.Binary)
	if err != nil {
		return err
	}

	return nil
}
