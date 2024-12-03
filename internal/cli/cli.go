package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/app/client"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var (
	ErrorNoExecution    = errors.New("no command execution found")
	ErrorUnknownCommand = errors.New("unknown command")
)

type commandExecute func() error
type commandThree map[string]command
type commandParams map[string]string

type command struct {
	desc        string
	usage       string
	execute     commandExecute
	subcommands commandThree
}

func (c *command) getSubCommandsNames() string {
	if c.subcommands != nil {
		cNames := ""
		for name, _ := range c.subcommands {
			cNames += fmt.Sprintf("%s ", name)
		}
		return cNames
	}
	return ""
}

type ClientCLI struct {
	localStorage *local.ClientStorage
	authService  *clientService.ClientAuthService
	dataService  *clientService.ClientDataService
	commandsRoot commandThree
	logger       *zap.SugaredLogger
}

func NewClientCLI(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *ClientCLI {

	localStorage := local.NewClientStorage()
	authService := clientService.NewClientAuthService(grpcClient, localStorage, logger)
	dataService := clientService.NewClientDataService(grpcClient, localStorage, logger)

	return &ClientCLI{
		localStorage: localStorage,
		authService:  authService,
		dataService:  dataService,
		logger:       logger.Named("CLI"),
	}
}

func (cli *ClientCLI) InitCommandsThree() {
	commandsRoot := commandThree{
		"AUTH": {
			desc: "Authorization in the system by login and password",
			execute: func() error {
				params := commandParams{
					"username": "",
					"password": "",
				}
				err := cli.authCommand(params)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"REGISTER": {
			desc: "Registration in the system",
			execute: func() error {
				params := commandParams{
					"username": "",
					"password": "",
					"email":    "",
				}
				err := cli.registerCommand(params)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"CREATE": {
			desc: "Creating a record in the system",
			subcommands: commandThree{
				"PAIR": {
					desc: "Create a key value pair",
					execute: func() error {
						params := commandParams{
							"key":  "",
							"pwd":  "",
							"meta": "",
						}
						err := cli.createPairCommand(params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"TEXT": {
					desc: "Creating text data",
					execute: func() error {
						log.Print("exec text")
						return nil
					},
				},
				"BINARY": {
					desc: "Creating binary data",
					execute: func() error {
						log.Print("exec binary")
						return nil
					},
				},
				"CARD": {
					desc: "Creating card data",
					execute: func() error {
						log.Print("exec binary")
						return nil
					},
				},
			},
		},
		"DEBUG": {
			desc: "Data output for the developer",
			execute: func() error {
				err := cli.dataService.Debug()
				if err != nil {
					return err
				}
				return nil
			},
		},
		"HELP": {
			desc: "Show information for help",
			execute: func() error {
				log.Print("exec help")
				return nil
			},
		},
		"EXIT": {
			desc: "Exiting the program",
			execute: func() error {
				os.Exit(1)
				return nil
			},
		},
	}
	cli.commandsRoot = commandsRoot
}

func (cli *ClientCLI) Run() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter root command")

		scanner.Scan()
		input = strings.ToUpper(scanner.Text())

		cmd, exists := cli.commandsRoot[input]
		if !exists {
			fmt.Printf("Unknown command. Run \"HELP\"\n")
			continue
		}

		err := cli.exec(cmd, scanner)
		if err != nil {
			fmt.Printf("Error while executing command: %v\n", err)
		}
	}
}

func (cli *ClientCLI) exec(cmd command, scanner *bufio.Scanner) error {
	if cmd.subcommands == nil {
		if cmd.execute != nil {
			err := cmd.execute()
			if err != nil {
				return err
			}
			return nil
		}
		return ErrorNoExecution
	}

	fmt.Printf("Enter one of: %s\n", cmd.getSubCommandsNames())

	scanner.Scan()
	input := strings.ToUpper(scanner.Text())

	cmd, exists := cmd.subcommands[input]
	if !exists {
		return ErrorUnknownCommand
	}

	err := cli.exec(cmd, scanner)
	if err != nil {
		return err
	}

	return nil
}

func (cli *ClientCLI) validateParams(params commandParams) commandParams {
	scanner := bufio.NewScanner(os.Stdin)

	for i, _ := range params {
		var input string

		for valid := false; !valid; {
			fmt.Printf("Enter %s: ", i)
			scanner.Scan()
			input = scanner.Text()
			if input != "" {
				valid = true
			}
		}

		params[i] = input
	}

	return params
}
