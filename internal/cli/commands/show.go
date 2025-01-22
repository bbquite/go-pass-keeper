package commands

import (
	"context"
	"fmt"
	clitable "github.com/bbquite/go-pass-keeper/pkg/table"
)

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
