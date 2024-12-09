package main

import (
	"fmt"
	"log"

	clientApp "github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli"
	"github.com/bbquite/go-pass-keeper/internal/config"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func showBuildInfo() {
	fmt.Printf("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n\n", buildVersion, buildDate, buildCommit)
}

func initServerLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	defer logger.Sync()
	return sugar, nil
}

func main() {
	showBuildInfo()

	logger, err := initServerLogger()
	if err != nil {
		log.Fatal(err)
	}

	cfg := new(config.ClientConfig)
	cfg.SetFlags()
	cfg.SetENV()

	logger.Infof("Client run with config: %s", cfg.PrintConfig())

	grpcClient, err := clientApp.NewGRPCClient(cfg.Host, cfg.RootCertPath)
	if err != nil {
		logger.Fatal(err)
	}

	cliClient := cli.NewClientCLI(grpcClient, logger)
	cliClient.Run()

	defer func() {
		err := grpcClient.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()
}
