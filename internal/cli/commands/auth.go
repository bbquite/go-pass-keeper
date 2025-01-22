package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bbquite/go-pass-keeper/internal/utils"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"

	"github.com/bbquite/go-pass-keeper/internal/models"
)

func (cm *CommandManager) accountAction(params CommandParams, action func(ctx context.Context, userData *models.UserAccountData) error) error {
	paramsValidated := cm.validateParams(params)
	loginData := &models.UserAccountData{
		Username: paramsValidated["username"].value,
		Password: paramsValidated["password"].value,
	}

	err := action(context.Background(), loginData)
	if err != nil {
		return err
	}

	return cm.saveTokenToFile()
}

func (cm *CommandManager) checkTokenWrapper(dataType models.DataTypeEnum, params CommandParams, action CommandActionWithTypeParams) error {
	token := cm.localStorage.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only, run \"AUTH\"")
	}

	return action(dataType, params)
}

func (cm *CommandManager) saveTokenToFile() error {
	jsOut, err := json.Marshal(cm.localStorage.Token)
	if err != nil {
		return err
	}
	err = utils.SaveFile(cm.authFilePath, jsOut)
	if err != nil {
		return err
	}
	return nil
}

func (cm *CommandManager) importTokenFromFile() error {
	var jwtModel jwttoken.JWT

	data, err := os.ReadFile(cm.authFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &jwtModel)
	if err != nil {
		return err
	}

	cm.localStorage.SetToken(&jwtModel)

	return nil
}
