package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type ServerConfig struct {
	Host        string `json:"host"`
	DatabaseURI string `json:"db_host"`
	JWTSecret   string `json:"jwt_secret"`
}

func (c *ServerConfig) GetFromENV() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New(".env file not found")
	}

	if envRunAddress, ok := os.LookupEnv("HOST"); ok {
		c.Host = envRunAddress
	}

	if envDATABASE, ok := os.LookupEnv("DATABASE_URI"); ok {
		c.DatabaseURI = envDATABASE
	}

	if envJWTSECRET, ok := os.LookupEnv("JWTSECRET"); ok {
		c.JWTSecret = envJWTSECRET
	} else {
		return errors.New("jwt access token not found")
	}

	return nil
}

func (c *ServerConfig) PrintConfig() string {
	jsonConfig, _ := json.Marshal(c)
	return fmt.Sprintf("Server run with config: %s", jsonConfig)
}
