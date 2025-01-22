package commands

import (
	"encoding/json"

	"github.com/bbquite/go-pass-keeper/internal/cli/validator"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/utils"
)

var (
	emptyParams = CommandParams{}

	authParams = CommandParams{
		"username": {validateFunc: validator.StringValidation},
		"password": {validateFunc: validator.StringValidation},
	}

	pairParams = CommandParams{
		"key":  {validateFunc: validator.StringValidation},
		"pwd":  {validateFunc: validator.StringValidation},
		"meta": {validateFunc: validator.StringValidation},
	}

	textParams = CommandParams{
		"text": {validateFunc: validator.StringValidationUnlimit},
		"meta": {validateFunc: validator.StringValidation},
	}

	binaryParams = CommandParams{
		"filepath": {validateFunc: validator.FilePathValidation},
		"meta":     {validateFunc: validator.StringValidation},
	}

	cardParams = CommandParams{
		"number": {
			validateFunc: validator.CardNumberValidation,
			usage:        "ex. 4242 4242 4242 4242",
		},
		"cvv": {
			validateFunc: validator.CardCvvValidation,
			usage:        "ex. 777",
		},
		"owner": {
			validateFunc: validator.StringValidation,
			usage:        "ex. IVAN IVANOV",
		},
		"exp": {
			validateFunc: validator.DateValidation,
			usage:        "ex. 01.06",
		},
		"meta": {validateFunc: validator.StringValidation},
	}
)

func wrapIDParam(params CommandParams) CommandParams {
	params["id"] = CommandParamsItem{validateFunc: validator.IntValidation}
	return params
}

func marshalDataByType(dataType models.DataTypeEnum, paramsValidated CommandParams) ([]byte, error) {
	var dataMarshal []byte
	var err error

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
		filePath := paramsValidated["filepath"].value
		fileContent, fileName, fileSize, ferr := utils.GetFileInfo(filePath)
		if ferr != nil {
			return nil, ferr
		}

		dataMarshal, err = json.Marshal(&models.BinaryData{
			FileName: fileName,
			FileSize: fileSize,
			Binary:   fileContent,
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
		return nil, err
	}
	return dataMarshal, nil
}
