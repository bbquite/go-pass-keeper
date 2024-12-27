package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
		filePath := paramsValidated["filepath"].value
		fileContent, fileName, fileSize, ferr := cm.getFileInfo(filePath)
		if ferr != nil {
			return err
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
		filePath := paramsValidated["filepath"].value
		fileContent, fileName, fileSize, err := cm.getFileInfo(filePath)
		if err != nil {
			return err
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

func (cm *CommandManager) showCommand() error {

	err := cm.dataService.GetData(context.Background())
	if err != nil {
		return err
	}

	cm.printPairs()
	cm.printCards()
	cm.printBinary()
	cm.printTexts()

	return nil
}

func (cm *CommandManager) downloadCommand(dataType models.DataTypeEnum, params ...CommandParams) error {

	err := cm.dataService.GetData(context.Background())
	if err != nil {
		return err
	}

	switch dataType {
	case models.DataTypePAIR:
		jsOut, err := json.MarshalIndent(cm.localStorage.PairsList, "", "	")
		if err != nil {
			return err
		}
		err = cm.saveFile("./pair export.json", jsOut)
		if err != nil {
			return err
		}

	case models.DataTypeTEXT:
		jsOut, err := json.MarshalIndent(cm.localStorage.TextsList, "", "	")
		if err != nil {
			return err
		}
		err = cm.saveFile("./text export.json", jsOut)
		if err != nil {
			return err
		}

	case models.DataTypeBINARY:
		var dataID uint64

		if len(params) > 0 {
			paramsValidated := cm.validateParams(params[0])

			if _, ok := paramsValidated["id"]; ok {
				dataID, err = strconv.ParseUint(paramsValidated["id"].value, 10, 32)
				if err != nil {
					return err
				}

				dataItem, err := cm.dataService.GetDataByID(context.Background(), uint32(dataID))
				if err != nil {
					return err
				}

				m := models.BinaryData{}
				d := []byte(dataItem.DataInfo)
				err = json.Unmarshal(d, &m)
				if err != nil {
					return err
				}

				err = cm.saveFile("./"+m.FileName, m.Binary)
				if err != nil {
					return err
				}

				return nil
			}
		}

		return fmt.Errorf("cant get binary data")

	case models.DataTypeCARD:
		jsOut, err := json.MarshalIndent(cm.localStorage.CardsList, "", "	")
		if err != nil {
			return err
		}
		err = cm.saveFile("./card export.json", jsOut)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cm *CommandManager) saveFile(filePath string, fileData []byte) error {
	err := os.WriteFile(filePath, fileData, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (cm *CommandManager) getFileInfo(filePath string) ([]byte, string, int64, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", 0, err
	}
	_, fileName := filepath.Split(filePath)

	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, "", 0, err
	}

	fileSize := fi.Size()

	return fileContent, fileName, fileSize, nil
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
		fmt.Println("--------------------")
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

func (cm *CommandManager) printBinary() {
	bin, _ := cm.localStorage.GetBinary()
	if len(bin) == 0 {
		return
	}

	binTable := clitable.New([]string{"ID", "NAME", "SIZE", "META"})
	binTable.Markdown = true

	for _, item := range bin {
		binTable.AddRow(map[string]interface{}{"ID": item.ID, "NAME": item.FileName, "SIZE": item.FileSize, "META": item.Meta})
	}

	fmt.Printf("\nBINARY DATA: \n\n")
	binTable.Print()
}
