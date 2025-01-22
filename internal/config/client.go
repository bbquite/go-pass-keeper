package config

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	defRemoteHost     = "localhost:8080"
	defClientCertPath = "./cert/ca.pem"
)

type ClientConfig struct {
	Host         string `json:"host" env:"HOST"`
	RootCertPath string `json:"root_cert_path" env:"CLIENT_CRT_PATH"`
}

func (c *ClientConfig) SetENV() error {
	err := env.Parse(c)
	if err != nil {
		return fmt.Errorf("env not found: %w", err)
	}

	return nil
}

func (c *ClientConfig) SetFlags() {
	flag.StringVar(&c.Host, "h", defRemoteHost, "remote host")
	flag.StringVar(&c.RootCertPath, "ca", defClientCertPath, "root cert path")
	flag.Parse()
}

func (c *ClientConfig) PrintConfig() []byte {
	jsonConfig, _ := json.Marshal(c)
	return jsonConfig
}
