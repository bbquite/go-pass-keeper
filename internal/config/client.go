package config

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type ClientConfig struct {
	Host         string `json:"host"`
	RootCertPath string `json:"root_cert_path"`
}

func (c *ClientConfig) GetFromENV() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New(".env file not found")
	}

	if envRunAddress, ok := os.LookupEnv("HOST"); ok {
		c.Host = envRunAddress
	}

	return nil
}

func (c *ClientConfig) PrintConfig() []byte {
	jsonConfig, _ := json.Marshal(c)
	return jsonConfig
}
