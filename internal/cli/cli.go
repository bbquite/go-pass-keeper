package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli/command"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var ErrorCLIGracefullyStop = errors.New("cli gracefully stop")

type ClientCLI struct {
	commandManager *command.CommandManager
	logger         *zap.SugaredLogger
}

func NewClientCLI(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *ClientCLI {
	commandManager := command.NewCommandManager(grpcClient, logger)

	return &ClientCLI{
		commandManager: commandManager,
		logger:         logger.Named("CLI"),
	}
}

func (cli *ClientCLI) Run() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter root command: ")

		scanner.Scan()
		input = strings.ToUpper(scanner.Text())

		cmd, exists := cli.commandManager.CommandRoot[input]
		if !exists {
			color.Red("Unknown command. Run \"HELP\"\n")
			continue
		}

		err := cli.exec(cmd, scanner)
		if err != nil {
			color.Red("Error while executing command: %v\n", err)
		} else {
			green := color.New(color.FgGreen)
			green.Printf("\n--- Successfully ---\n\n")
		}
	}
}

func (cli *ClientCLI) exec(cmd command.Command, scanner *bufio.Scanner) error {
	if cmd.Subcommands == nil {
		if cmd.Execute != nil {
			err := cmd.Execute()
			if err != nil {
				return err
			}
			return nil
		}
		return command.ErrorNoExecution
	}

	fmt.Printf("Enter one of: %s\n", cmd.GetSubCommandsNames())

	scanner.Scan()
	input := strings.ToUpper(scanner.Text())

	cmd, exists := cmd.Subcommands[input]
	if !exists {
		return command.ErrorUnknownCommand
	}

	err := cli.exec(cmd, scanner)
	if err != nil {
		return err
	}

	return nil
}
