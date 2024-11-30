package command

import (
	"github.com/bbquite/go-pass-keeper/internal/service/client"
)

type DebugCommand struct {
	authService *client.ClientAuthService
}

func NewDebugCommand(authService *client.ClientAuthService) *DebugCommand {
	return &DebugCommand{
		authService: authService,
	}
}

func (c *DebugCommand) Name() string {
	return "DEBUG"
}

func (c *DebugCommand) Desc() string {
	return "TO BE ADDED"
}

func (c *DebugCommand) Usage() string {
	return "TO BE ADDED"
}

func (c *DebugCommand) Execute() error {
	err := c.authService.Debug()
	if err != nil {
		return err
	}
	return nil
}
