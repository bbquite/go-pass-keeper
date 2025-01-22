package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli/commands"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

type ClientCLI struct {
	commandManager *commands.CommandManager
	logger         *zap.SugaredLogger
}

func NewClientCLI(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *ClientCLI {
	commandManager := commands.NewCommandManager(grpcClient, logger)

	return &ClientCLI{
		commandManager: commandManager,
		logger:         logger.Named("CLI"),
	}
}

func (cli *ClientCLI) Run() error {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter root commands: ")

		scanner.Scan()
		input = strings.ToUpper(scanner.Text())

		cmd, exists := cli.commandManager.CommandRoot[input]
		if !exists {
			color.Red("Unknown commands. Run \"HELP\"\n")
			continue
		}

		err := cli.exec(cmd, scanner)
		if err != nil {
			if errors.Is(err, commands.ErrorGracefullyStop) {
				return err
			}
			color.Red("Error while executing commands: %v\n", err)
		} else {
			green := color.New(color.FgGreen)
			green.Printf("\n--- Successfully ---\n\n")
		}
	}
}

func (cli *ClientCLI) exec(cmd commands.Command, scanner *bufio.Scanner) error {
	if cmd.Subcommands == nil {
		if cmd.Execute != nil {
			err := cmd.Execute()
			if err != nil {
				return err
			}
			return nil
		}
		return commands.ErrorNoExecution
	}

	fmt.Printf("Enter one of: %s\n", strings.Join(cmd.GetSubCommandsNames(), " | "))

	scanner.Scan()
	input := strings.ToUpper(scanner.Text())

	cmd, exists := cmd.Subcommands[input]
	if !exists {
		return commands.ErrorUnknownCommand
	}

	err := cli.exec(cmd, scanner)
	if err != nil {
		return err
	}

	return nil
}
