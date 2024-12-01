package command

import (
	"github.com/bbquite/go-pass-keeper/internal/service/client"
)

type DebugCommand struct {
	dataService *client.ClientDataService
}

func NewDebugCommand(dataService *client.ClientDataService) *DebugCommand {
	return &DebugCommand{
		dataService: dataService,
	}
}

func (c *DebugCommand) Name() string {
	return "DEBUG"
}

func (c *DebugCommand) Desc() string {
	return "TO BE ADDED"
}

func (c *DebugCommand) Usage() string {
	return "\t\tTO BE ADDED\n"
}

func (c *DebugCommand) Execute() error {
	err := c.dataService.Debug()
	if err != nil {
		return err
	}
	return nil
}
