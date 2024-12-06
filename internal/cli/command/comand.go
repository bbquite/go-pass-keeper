package command

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli/validator"
	"github.com/bbquite/go-pass-keeper/internal/models"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var (
	ErrorNoExecution    = errors.New("no command execution found")
	ErrorUnknownCommand = errors.New("unknown command")
)

type CommandExecute func() error
type CommandThree map[string]Command

type CommandParams map[string]struct {
	validateFunc validator.ValidateFunc
	usage        string
	value        string
}

type Command struct {
	Desc        string
	Usage       string
	Execute     CommandExecute
	Subcommands CommandThree
}

func (c *Command) GetSubCommandsNames() string {
	if c.Subcommands != nil {
		cNames := "| "
		for name, _ := range c.Subcommands {
			cNames += fmt.Sprintf("%s | ", name)
		}
		return cNames
	}
	return ""
}

type CommandManager struct {
	localStorage *local.ClientStorage
	authService  *clientService.ClientAuthService
	dataService  *clientService.ClientDataService
	CommandRoot  CommandThree
}

func NewCommandManager(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *CommandManager {

	localStorage := local.NewClientStorage()
	authService := clientService.NewClientAuthService(grpcClient, localStorage, logger)
	dataService := clientService.NewClientDataService(grpcClient, localStorage, logger)

	cm := &CommandManager{
		localStorage: localStorage,
		authService:  authService,
		dataService:  dataService,
	}

	cm.initCommandsThree()

	return cm
}

func (cm *CommandManager) validateParams(params CommandParams) CommandParams {

	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for name, body := range params {
		for {
			usage := ""
			if body.usage != "" {
				usage = fmt.Sprintf(" (%s)", body.usage)
			}
			fmt.Printf("Enter %s%s: ", name, usage)
			scanner.Scan()
			input = scanner.Text()

			err := body.validateFunc(input)
			if err == nil {
				break
			}
			color.Red("%v", err)
		}

		// https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
		if entry, ok := params[name]; ok {
			entry.value = input
			params[name] = entry
		}
	}

	return params
}

func (cm *CommandManager) initCommandsThree() {
	commandRoot := CommandThree{
		"AUTH": {
			Desc: "Authorization in the system by login and password",
			Execute: func() error {
				params := CommandParams{
					"username": {
						validateFunc: validator.StringValidation,
						value:        "",
					},
					"password": {
						validateFunc: validator.StringValidation,
						value:        "",
					},
				}
				err := cm.authCommand(params)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"REGISTER": {
			Desc: "Registration in the system",
			Execute: func() error {
				params := CommandParams{
					"username": {
						validateFunc: validator.StringValidation,
						value:        "",
					},
					"password": {
						validateFunc: validator.StringValidation,
						value:        "",
					},
				}
				err := cm.registerCommand(params)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"SHOW": {
			Desc: "Show records from remote server",
			Execute: func() error {
				err := cm.getCommand()
				if err != nil {
					return err
				}
				return nil
			},
		},
		"CREATE": {
			Desc: "Creating a record in the system",
			Subcommands: CommandThree{
				"PAIR": {
					Desc: "Create a key value pair",
					Execute: func() error {
						params := CommandParams{
							"key": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
							"pwd": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}

						err := cm.createCommand(models.DataTypePAIR, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"TEXT": {
					Desc: "Creating text data",
					Execute: func() error {
						params := CommandParams{
							"text": {
								validateFunc: validator.StringValidationUnlimit,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.createCommand(models.DataTypeTEXT, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"BINARY": {
					Desc: "Creating binary data",
					Execute: func() error {
						params := CommandParams{
							"filepath": {
								validateFunc: validator.FilePathValidation,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.createCommand(models.DataTypeBINARY, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"CARD": {
					Desc: "Creating card data",
					Execute: func() error {
						params := CommandParams{
							"number": {
								validateFunc: validator.CardNumberValidation,
								usage:        "4242 4242 4242 4242",
								value:        "",
							},
							"cvv": {
								validateFunc: validator.CardCvvValidation,
								usage:        "777",
								value:        "",
							},
							"owner": {
								validateFunc: validator.StringValidation,
								usage:        "IVAN IVANOV",
								value:        "",
							},
							"exp": {
								validateFunc: validator.DateValidation,
								usage:        "ex. 01.06",
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.createCommand(models.DataTypeCARD, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
			},
		},
		"UPDATE": {
			Desc: "Updating a record in the system",
			Subcommands: CommandThree{
				"PAIR": {
					Desc: "Updating a key value pair",
					Execute: func() error {
						params := CommandParams{
							"id": {
								validateFunc: validator.IntValidation,
								value:        "",
							},
							"key": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
							"pwd": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.updateCommand(models.DataTypePAIR, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"TEXT": {
					Desc: "Updating text data",
					Execute: func() error {
						params := CommandParams{
							"id": {
								validateFunc: validator.IntValidation,
								value:        "",
							},
							"text": {
								validateFunc: validator.StringValidationUnlimit,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.updateCommand(models.DataTypeTEXT, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"BINARY": {
					Desc: "Updating binary data",
					Execute: func() error {
						params := CommandParams{
							"id": {
								validateFunc: validator.IntValidation,
								value:        "",
							},
							"filepath": {
								validateFunc: validator.FilePathValidation,
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.updateCommand(models.DataTypeBINARY, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
				"CARD": {
					Desc: "Updating card data",
					Execute: func() error {
						params := CommandParams{
							"id": {
								validateFunc: validator.IntValidation,
								value:        "",
							},
							"number": {
								validateFunc: validator.CardNumberValidation,
								usage:        "4242 4242 4242 4242",
								value:        "",
							},
							"cvv": {
								validateFunc: validator.CardCvvValidation,
								usage:        "777",
								value:        "",
							},
							"owner": {
								validateFunc: validator.StringValidation,
								usage:        "IVAN IVANOV",
								value:        "",
							},
							"exp": {
								validateFunc: validator.DateValidation,
								usage:        "ex. 01.06",
								value:        "",
							},
							"meta": {
								validateFunc: validator.StringValidation,
								value:        "",
							},
						}
						err := cm.updateCommand(models.DataTypeCARD, params)
						if err != nil {
							return err
						}
						return nil
					},
				},
			},
		},
		"DELETE": {
			Desc: "Deleting a record in the system",
			Execute: func() error {
				params := CommandParams{
					"id": {
						validateFunc: validator.IntValidation,
						value:        "",
					},
				}
				err := cm.deleteCommand(params)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"DEBUG": {
			Desc: "Data output for the developer",
			Execute: func() error {
				err := cm.dataService.Debug()
				if err != nil {
					return err
				}
				return nil
			},
		},
		"HELP": {
			Desc: "Show information for help",
			Execute: func() error {
				log.Print("exec help")
				return nil
			},
		},
	}
	cm.CommandRoot = commandRoot
}
