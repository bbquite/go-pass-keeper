package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbquite/go-pass-keeper/internal/models"
	clitable "github.com/bbquite/go-pass-keeper/pkg/table"
)

func (cm *CommandManager) createCommand(dataType models.DataTypeEnum, params CommandParams) error {

	token := cm.localStorage.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only, run \"AUTH\"")
	}

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

	token := cm.localStorage.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only, run \"AUTH\"")
	}

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

	token := cm.localStorage.GetToken()
	if token == "" {
		return fmt.Errorf("authorization only, run \"AUTH\"")
	}

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

func (cm *CommandManager) getCommand() error {

	err := cm.dataService.GetData(context.Background())
	if err != nil {
		return err
	}

	cm.printPairs()
	cm.printCards()
	cm.printTexts()

	return nil
}

func (cm *CommandManager) printPairs() {
	pairs, _ := cm.localStorage.GetPairs()
	if len(pairs) == 0 {
		return
	}

	pairsTable := clitable.New([]string{"ID", "KEY", "PWD", "META"})
	pairsTable.Markdown = true

	for _, item := range pairs {
		pairsTable.AddRow(map[string]interface{}{"ID": item.ID, "KEY": item.Key, "PWD": item.Pwd, "META": item.Meta})
	}

	fmt.Printf("\nPAIR DATA: \n\n")
	pairsTable.Print()
}

func (cm *CommandManager) printTexts() {
	texts, _ := cm.localStorage.GetTexts()
	if len(texts) == 0 {
		return
	}

	fmt.Printf("\nTEXT DATA: \n\n")
	for _, item := range texts {
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Meta: %s\n", item.Meta)
		fmt.Printf("Text: %s\n", item.Text)
	}
}

func (cm *CommandManager) printCards() {
	cards, _ := cm.localStorage.GetCards()
	if len(cards) == 0 {
		return
	}

	cardsTable := clitable.New([]string{"ID", "NUM", "CVV", "EXP", "OWNER", "META"})
	cardsTable.Markdown = true

	for _, item := range cards {
		cardsTable.AddRow(map[string]interface{}{
			"ID": item.ID, "NUM": item.CardNum,
			"CVV": item.CardCvv, "EXP": item.CardExp,
			"OWNER": item.CardOwner, "META": item.Meta})
	}

	fmt.Printf("\nCARD DATA: \n\n")
	cardsTable.Print()
}
