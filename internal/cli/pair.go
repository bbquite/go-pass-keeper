package cli

import (
	"context"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"strconv"
)

func (cli *ClientCLI) createPairCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.dataService.CreatePairData(context.Background(), &models.PairData{
		Key:  paramsV["key"],
		Pwd:  paramsV["pwd"],
		Meta: paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) updatePairCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.UpdatePairData(context.Background(), &models.PairData{
		ID:   id,
		Key:  paramsV["key"],
		Pwd:  paramsV["pwd"],
		Meta: paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) deletePairCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.DeletePairData(context.Background(), &models.PairData{ID: id})
	if err != nil {
		return err
	}
	return nil
}
