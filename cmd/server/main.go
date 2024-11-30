package main

import (
	"flag"
	"fmt"
	"github.com/bbquite/go-pass-keeper/internal/config"
	"log"

	"github.com/bbquite/go-pass-keeper/internal/app/server"
	"go.uber.org/zap"
)

const (
	defServerHost  = "localhost:8080"
	defDatabaseURI = ""
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

	cfg := new(config.ServerConfig)
	flag.StringVar(&cfg.Host, "h", defServerHost, "HOST")
	flag.StringVar(&cfg.DatabaseURI, "d", defDatabaseURI, "DB HOST")
	flag.Parse()

	err = cfg.GetFromENV()
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
