package command

import (
	"context"

	"github.com/bbquite/go-pass-keeper/internal/models"
)

func (cm *CommandManager) authCommand(params CommandParams) error {

	paramsValidated := cm.validateParams(params)

	err := cm.authService.AuthUser(context.Background(), &models.UserLoginData{
		Username: paramsValidated["username"].value,
		Password: paramsValidated["password"].value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (cm *CommandManager) registerCommand(params CommandParams) error {

	paramsValidated := cm.validateParams(params)

	err := cm.authService.RegisterUser(context.Background(), &models.UserRegisterData{
		Username: paramsValidated["username"].value,
		Password: paramsValidated["password"].value,
	})
	if err != nil {
		return err
	}
	return nil
}
