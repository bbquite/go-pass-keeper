package cli

import (
	"context"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"strconv"
)

func (cli *ClientCLI) createBinaryCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.dataService.CreateBinaryData(context.Background(), &models.BinaryData{
		Binary: []byte{},
		Meta:   paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) updateBinaryCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.UpdateBinaryData(context.Background(), &models.BinaryData{
		ID:     id,
		Binary: []byte{},
		Meta:   paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) deleteBinaryCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.DeleteBinaryData(context.Background(), &models.BinaryData{ID: id})
	if err != nil {
		return err
	}
	return nil
}
