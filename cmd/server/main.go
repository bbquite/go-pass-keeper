package main

import (
	"flag"
	"fmt"
	app "github.com/bbquite/go-pass-keeper/internal/app/server"
	"log"
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

func main() {
	showBuildInfo()

	cfg := new(app.ServerConfig)
	flag.StringVar(&cfg.Host, "h", defServerHost, "HOST")
	flag.StringVar(&cfg.DatabaseURI, "d", defDatabaseURI, "DB HOST")
	flag.Parse()

	srv, err := app.NewGRPCServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv.RunGRPCServer()
}
