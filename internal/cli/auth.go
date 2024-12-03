package cli

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/models"
)

func (cli *ClientCLI) authCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.authService.AuthUser(context.Background(), &models.UserLoginData{
		Username: paramsV["username"],
		Password: paramsV["password"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) registerCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.authService.RegisterUser(context.Background(), &models.UserRegisterData{
		Username: paramsV["username"],
		Password: paramsV["password"],
		Email:    paramsV["email"],
	})
	if err != nil {
		return err
	}
	return nil
}
