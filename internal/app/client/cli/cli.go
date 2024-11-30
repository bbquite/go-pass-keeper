package cli

import (
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/app/client/command"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	"go.uber.org/zap"
	"os"
	"strings"
)

func RunCLI(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) error {
	localStorage := local.NewClientStorage()
	authService := clientService.NewClientAuthService(grpcClient, localStorage, logger)

	commandsRoot := []command.ClientCommand{
		command.NewRegisterCommand(authService, os.Stdin, os.Stdout),
	}

	commandNames := make([]string, len(commandsRoot))
	commandMap := make(map[string]command.ClientCommand)
	for i, cmd := range commandsRoot {
		commandMap[cmd.Name()] = cmd
		commandNames[i] = cmd.Name()
	}

	fmt.Println("Available commands: ", strings.Join(commandNames, "\n"))

	for {
		fmt.Print("Enter the command: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			logger.Info("Command input error: %v", err)
		}
		cmd, exists := commandMap[strings.ToUpper(input)]
		if !exists {
			logger.Info("Unknown command. Run \"HELP\" command", input)
			continue
		}

		err = cmd.Execute()
		if err != nil {
			logger.Errorf("error while executing command: %v", err)
		}
	}
}
