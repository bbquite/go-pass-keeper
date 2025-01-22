package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli/validator"
	"github.com/bbquite/go-pass-keeper/internal/models"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var (
	ErrorNoExecution    = errors.New("no commands execution found")
	ErrorUnknownCommand = errors.New("unknown commands")
	ErrorGracefullyStop = errors.New("cli gracefully stop")
)

type (
	CommandActionWithTypeParams func(dataType models.DataTypeEnum, params CommandParams) error

	CommandParamsItem struct {
		validateFunc validator.ValidateFunc
		usage        string
		value        string
	}

	CommandParams  map[string]CommandParamsItem
	CommandExecute func() error
	CommandThree   map[string]Command
)

type Command struct {
	Desc        string
	Execute     CommandExecute
	Subcommands CommandThree
}

func (c *Command) GetSubCommandsNames() []string {
	if c.Subcommands != nil {
		var cNames []string
		for name := range c.Subcommands {
			cNames = append(cNames, name)
		}
		return cNames
	}
	return nil
}

func (c *Command) GetCommandHelp(lvl int, cmd Command) string {
	var resultHelp string

	if len(cmd.Subcommands) > 0 {

		lvlTab := strings.Repeat("    ", lvl)

		for k, v := range cmd.Subcommands {
			resultHelp += fmt.Sprintf("%s %s - %s\n", lvlTab, k, v.Desc)

			if len(v.Subcommands) > 0 {
				lvl += 1
				resultHelp += cmd.GetCommandHelp(lvl, v)
			}
		}
	}

	return resultHelp
}

type CommandManager struct {
	localStorage       *local.ClientStorage
	authService        *clientService.ClientAuthService
	dataService        *clientService.ClientDataService
	authFilePath       string
	pairExportFilePath string
	textExportFilePath string
	cardExportFilePath string
	helpInfo           string
	CommandRoot        CommandThree
}

func NewCommandManager(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *CommandManager {

	localStorage := local.NewClientStorage()
	authService := clientService.NewClientAuthService(grpcClient, localStorage, logger)
	dataService := clientService.NewClientDataService(grpcClient, localStorage, logger)

	cm := &CommandManager{
		localStorage:       localStorage,
		authService:        authService,
		dataService:        dataService,
		authFilePath:       "./auth.json",
		pairExportFilePath: "./pairExport.json",
		textExportFilePath: "./textExport.json",
		cardExportFilePath: "./cardExport.json",
	}

	cm.initCommandsThree()
	cm.generateHelpInfo()

	err := cm.importTokenFromFile()
	if err != nil {
		fmt.Printf("auth file not found: %v\n\n", err)
	}

	fmt.Print(cm.helpInfo)

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

func (cm *CommandManager) generateHelpInfo() {
	var help string
	for k, cmd := range cm.CommandRoot {
		help += fmt.Sprintf("\n%s - %s\n%s", k, cmd.Desc, cmd.GetCommandHelp(1, cmd))
	}
	cm.helpInfo = help + "\n"
}

func (cm *CommandManager) initCommandsThree() {
	commandRoot := CommandThree{
		"AUTH": {
			Desc: "Authorization in the system by login and password",
			Execute: func() error {
				return cm.accountAction(authParams, cm.authService.AuthUser)
			},
		},
		"REG": {
			Desc: "Registration in the system",
			Execute: func() error {
				return cm.accountAction(authParams, cm.authService.RegisterUser)
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
			Subcommands: cm.initExportCommands(),
		},
		"CREATE": {
			Desc:        "Creating a record in the system",
			Subcommands: cm.initCreateCommands(),
		},
		"UPDATE": {
			Desc:        "Updating a record in the system",
			Subcommands: cm.initUpdateCommands(),
		},
		"DEL": {
			Desc: "Deleting a record in the system",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeUNDEFINE, wrapIDParam(emptyParams), cm.deleteCommand)
			},
		},
		"DEBUG": {
			Desc: "Data output for the developer",
			Execute: func() error {
				return cm.dataService.Debug()
			},
		},
		"HELP": {
			Desc: "Show help information",
			Execute: func() error {
				fmt.Print(cm.helpInfo)
				return nil
			},
		},
		"EXIT": {
			Desc: "Close CLI",
			Execute: func() error {
				return ErrorGracefullyStop
			},
		},
	}
	cm.CommandRoot = commandRoot
}

func (cm *CommandManager) initExportCommands() CommandThree {
	var p CommandParams
	cmThree := CommandThree{
		"PAIR": {
			Desc: "Download a key value pair (ALL)",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypePAIR, p, cm.exportCommand)
			},
		},
		"TEXT": {
			Desc: "Download text data (ALL)",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeTEXT, p, cm.exportCommand)
			},
		},
		"FILE": {
			Desc: "Download binary data (by ID)",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeBINARY, wrapIDParam(emptyParams), cm.exportCommand)
			},
		},
		"CARD": {
			Desc: "Download card data (ALL)",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeCARD, p, cm.exportCommand)
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
				return cm.checkTokenWrapper(models.DataTypePAIR, pairParams, cm.createCommand)
			},
		},
		"TEXT": {
			Desc: "Creating text data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeTEXT, textParams, cm.createCommand)
			},
		},
		"FILE": {
			Desc: "Creating binary data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeBINARY, binaryParams, cm.createCommand)
			},
		},
		"CARD": {
			Desc: "Creating card data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeCARD, cardParams, cm.createCommand)
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
				return cm.checkTokenWrapper(models.DataTypePAIR, wrapIDParam(pairParams), cm.updateCommand)
			},
		},
		"TEXT": {
			Desc: "Updating text data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeTEXT, wrapIDParam(textParams), cm.updateCommand)
			},
		},
		"FILE": {
			Desc: "Updating binary data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeBINARY, wrapIDParam(binaryParams), cm.updateCommand)
			},
		},
		"CARD": {
			Desc: "Updating card data",
			Execute: func() error {
				return cm.checkTokenWrapper(models.DataTypeCARD, wrapIDParam(cardParams), cm.updateCommand)
			},
		},
	}

	return cmThree
}
