package cli

import (
	"context"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"strconv"
)

func (cli *ClientCLI) createCardCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)
	err := cli.dataService.CreateCardData(context.Background(), &models.CardData{
		CardNum:   paramsV["card number"],
		CardCvv:   paramsV["card cvv"],
		CardOwner: paramsV["card owner"],
		CardExp:   paramsV["card exp"],
		Meta:      paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) updateCardCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.UpdateCardData(context.Background(), &models.CardData{
		ID:        id,
		CardNum:   paramsV["card number"],
		CardCvv:   paramsV["card cvv"],
		CardOwner: paramsV["card owner"],
		CardExp:   paramsV["card exp"],
		Meta:      paramsV["meta"],
	})
	if err != nil {
		return err
	}
	return nil
}

func (cli *ClientCLI) deleteCardCommand(params map[string]string) error {
	paramsV := cli.validateParams(params)

	idInt, err := strconv.ParseInt(paramsV["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	id := uint32(idInt)

	err = cli.dataService.DeleteCardData(context.Background(), &models.CardData{ID: id})
	if err != nil {
		return err
	}
	return nil
}
