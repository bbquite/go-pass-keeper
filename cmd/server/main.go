package main

import (
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/config"
	"log"

	"github.com/bbquite/go-pass-keeper/internal/app/server"
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

	cfg := new(config.ServerConfig)
	cfg.SetFlags()
	err = cfg.SetENV()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Server run with config: %s", cfg.PrintConfig())

	srv, err := server.NewGRPCServer(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	srv.RunGRPCServer()
}
