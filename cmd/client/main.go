package main

import (
	"flag"
	"fmt"
	clientApp "github.com/bbquite/go-pass-keeper/internal/app/client"
	"github.com/bbquite/go-pass-keeper/internal/cli"
	"github.com/bbquite/go-pass-keeper/internal/config"
	"go.uber.org/zap"
	"log"
)

const (
	defServerHost = "localhost:8080"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func showBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
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
	flag.StringVar(&cfg.Host, "h", defServerHost, "remote host")
	flag.StringVar(&cfg.RootCertPath, "ca", "./cert/ca.pem", "root cert path")
	flag.Parse()

	err = cfg.GetFromENV()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Client run with config: %s", cfg.PrintConfig())

	grpcClient, err := clientApp.NewGRPCClient(cfg.Host, cfg.RootCertPath)
	if err != nil {
		logger.Fatal(err)
	}

	cliClient := cli.NewClientCLI(grpcClient, logger)
	err = cliClient.Run()
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		err := grpcClient.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()
}
