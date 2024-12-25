package command

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli/validator"
	"github.com/bbquite/go-pass-keeper/internal/models"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var (
	ErrorNoExecution    = errors.New("no command execution found")
	ErrorUnknownCommand = errors.New("unknown command")
)

type CommandExecute func() error
type CommandActionWithParams func(params CommandParams) error
type CommandThree map[string]Command

type CommandParamsItem struct {
	validateFunc validator.ValidateFunc
	usage        string
	value        string
}

type CommandParams map[string]CommandParamsItem

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
	authFilePath string
	helpInfo     string
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
		authFilePath: "./auth.json",
		helpInfo:     "",
	}

	cm.initCommandsThree()
	err := cm.authFromFile()
	if err != nil {
		fmt.Printf("auth file not found: %v\n", err)
	}

	return cm
}

func (cm *CommandManager) authFromFile() error {
	var jwtModel jwttoken.JWT

	data, err := os.ReadFile(cm.authFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &jwtModel)
	if err != nil {
		return err
	}

	cm.localStorage.SetToken(&jwtModel)

	return nil
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

func (cm *CommandManager) executeWithParams(params CommandParams, action CommandActionWithParams) error {
	validatedParams := cm.validateParams(params)
	return action(validatedParams)
}

func (cm *CommandManager) initCommandsThree() {
	commandRoot := CommandThree{
		"AUTH": {
			Desc: "Authorization in the system by login and password",
			Execute: func() error {
				return cm.executeWithParams(authParams, cm.authCommand)
			},
		},
		"REGISTER": {
			Desc: "Registration in the system",
			Execute: func() error {
				return cm.executeWithParams(authParams, cm.registerCommand)
			},
		},
		"SHOW": {
			Desc: "Show records from remote server",
			Execute: func() error {
				return cm.showCommand()
			},
		},
		"GET": {
			Desc:        "Download data from remote server",
			Subcommands: cm.initGetCommands(),
		},
		"CREATE": {
			Desc:        "Creating a record in the system",
			Subcommands: cm.initCreateCommands(),
		},
		"UPDATE": {
			Desc:        "Updating a record in the system",
			Subcommands: cm.initUpdateCommands(),
		},
		"DELETE": {
			Desc: "Deleting a record in the system",
			Execute: func() error {
				return cm.executeWithParams(onlyIDParams, cm.deleteCommand)
			},
		},
		"DEBUG": {
			Desc: "Data output for the developer",
			Execute: func() error {
				return cm.dataService.Debug()
			},
		},
		"HELP": {
			Desc: "Show information for help",
			Execute: func() error {
				log.Print(cm.helpInfo)
				return nil
			},
		},
	}
	cm.CommandRoot = commandRoot
}

func (cm *CommandManager) initGetCommands() CommandThree {
	cmThree := CommandThree{
		"PAIR": {
			Desc: "Download a key value pair (ALL)",
			Execute: func() error {
				return cm.downloadCommand(models.DataTypePAIR)
			},
		},
		"TEXT": {
			Desc: "Download text data (ALL)",
			Execute: func() error {
				return cm.downloadCommand(models.DataTypeTEXT)
			},
		},
		"BINARY": {
			Desc: "Download binary data (by ID)",
			Execute: func() error {
				return cm.downloadCommand(models.DataTypeBINARY, onlyIDParams)
			},
		},
		"CARD": {
			Desc: "Creating card data",
			Execute: func() error {
				return cm.downloadCommand(models.DataTypeCARD)
			},
		},
	}

	return cmThree
}

func (cm *CommandManager) initCreateCommands() CommandThree {
	cmThree := CommandThree{
		"PAIR": {
			Desc: "Create a key value pair",
			Execute: func() error {
				return cm.createCommand(models.DataTypePAIR, pairParams)
			},
		},
		"TEXT": {
			Desc: "Creating text data",
			Execute: func() error {
				return cm.createCommand(models.DataTypeTEXT, textParams)
			},
		},
		"BINARY": {
			Desc: "Creating binary data",
			Execute: func() error {
				return cm.createCommand(models.DataTypeBINARY, binaryParams)
			},
		},
		"CARD": {
			Desc: "Creating card data",
			Execute: func() error {
				return cm.createCommand(models.DataTypeCARD, cardParams)
			},
		},
	}

	return cmThree
}

func (cm *CommandManager) initUpdateCommands() CommandThree {
	cmThree := CommandThree{
		"PAIR": {
			Desc: "Updating a key value pair",
			Execute: func() error {
				return cm.updateCommand(models.DataTypePAIR, wrapIDParam(pairParams))
			},
		},
		"TEXT": {
			Desc: "Updating text data",
			Execute: func() error {
				return cm.updateCommand(models.DataTypeTEXT, wrapIDParam(textParams))
			},
		},
		"BINARY": {
			Desc: "Updating binary data",
			Execute: func() error {
				return cm.updateCommand(models.DataTypeBINARY, wrapIDParam(binaryParams))
			},
		},
		"CARD": {
			Desc: "Updating card data",
			Execute: func() error {
				return cm.updateCommand(models.DataTypeCARD, wrapIDParam(cardParams))
			},
		},
	}

	return cmThree
}
