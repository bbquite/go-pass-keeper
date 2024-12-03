package cli

import (
	"context"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"strconv"
)

func (cli *ClientCLI) createTextCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.dataService.CreateTextData(context.Background(), &models.TextData{
		Text: paramsV["text"],
		Meta: paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) updateTextCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.UpdateTextData(context.Background(), &models.TextData{
		ID:   id,
		Text: paramsV["text"],
		Meta: paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) deleteTextCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.DeleteTextData(context.Background(), &models.TextData{ID: id})
	if err != nil {
		return err
	}
	return nil
}
