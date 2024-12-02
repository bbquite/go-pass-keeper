package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/app/client/command"
	clientService "github.com/bbquite/go-pass-keeper/internal/service/client"
	"github.com/bbquite/go-pass-keeper/internal/storage/local"
	"go.uber.org/zap"
)

type commandsNames []string
type commandsMap map[string]command.ClientCommand
type commandUsage string

type ClientCLI struct {
	localStorage *local.ClientStorage
	authService  *clientService.ClientAuthService
	dataService  *clientService.ClientDataService
	commandsRoot []command.ClientCommand
	cNames       commandsNames
	cMap         commandsMap
	cUsage       commandUsage
	logger       *zap.SugaredLogger
}

func NewClientCLI(grpcClient *client.GRPCClient, logger *zap.SugaredLogger) *ClientCLI {

	localStorage := local.NewClientStorage()
	authService := clientService.NewClientAuthService(grpcClient, localStorage, logger)
	dataService := clientService.NewClientDataService(grpcClient, localStorage, logger)

	commandsRoot := []command.ClientCommand{
		command.NewRegisterCommand(authService, os.Stdin, os.Stdout),
		command.NewAuthCommand(authService, os.Stdin, os.Stdout),
		command.NewDebugCommand(dataService),
	}

	cNames, cMap, cUsage := buildCommandsInfo(commandsRoot)

	return &ClientCLI{
		localStorage: localStorage,
		authService:  authService,
		dataService:  dataService,
		commandsRoot: commandsRoot,
		cNames:       cNames,
		cMap:         cMap,
		cUsage:       cUsage,
		logger:       logger.Named("CLI"),
	}
}

func buildCommandsInfo(commandsRoot []command.ClientCommand) (commandsNames, commandsMap, commandUsage) {

	var cNames commandsNames
	cUsage := "\n\nHELP INFO: \n"
	cMap := make(map[string]command.ClientCommand)

	for _, cmd := range commandsRoot {
		cMap[cmd.Name()] = cmd
		cNames = append(cNames, cmd.Name())

		cUsage += fmt.Sprintf("\t%s - %s\n%s", cmd.Name(), cmd.Desc(), cmd.Usage())
	}

	log.Print(cUsage)

	return cNames, cMap, commandUsage(cUsage)
}

func (cli *ClientCLI) Run() error {
	for {
		fmt.Print("Enter the command: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			cli.logger.Infof("Command input error: %v", err)
			continue
		}
		cmd, exists := cli.cMap[strings.ToUpper(input)]
		if !exists {
			cli.logger.Infof("Unknown command. Run \"HELP\" command %s", input)
			continue
		}

		err = cmd.Execute()
		if err != nil {
			cli.logger.Infof("error while executing command: %v", err)
		}
	}
}
