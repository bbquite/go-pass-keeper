package config

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	defServerHost   = "localhost:8080"
	defDatabaseHost = "host=localhost user=postgres password=123 dbname=gopasskeeper sslmode=disable"
	defSecretKey    = "ytrfedghjfgjkhk123"
	defCryptoKey    = "01234567890123456789012345678901"
	defServerKey    = "./cert/server.key"
	defServerCrt    = "./cert/server.crt"
)

type ServerConfig struct {
	Host          string `json:"host" env:"HOST"`
	DatabaseURI   string `json:"db_host" env:"DH_HOST"`
	JWTSecret     string `json:"jwt_secret" env:"SECRET_KEY"`
	CryptoKey     string `json:"crypto_key" env:"CRYPTO_KEY"`
	ServerKeyPath string `json:"server_key" env:"SERVER_KEY_PATH"`
	ServerCrtPath string `json:"server_crt" env:"SERVER_CRT_PATH"`
}

func (c *ServerConfig) SetENV() error {
	err := env.Parse(c)
	if err != nil {
		return fmt.Errorf("env not found: %w", err)
	}

	return nil
}

func (c *ServerConfig) SetFlags() {
	flag.StringVar(&c.Host, "h", defServerHost, "server host")
	flag.StringVar(&c.DatabaseURI, "d", defDatabaseHost, "db host uri")
	flag.StringVar(&c.JWTSecret, "secret-key", defSecretKey, "secret key for hash")
	flag.StringVar(&c.CryptoKey, "crypto-key", defCryptoKey, "crypto key")
	flag.StringVar(&c.ServerKeyPath, "server-key", defServerKey, "path to server key")
	flag.StringVar(&c.ServerCrtPath, "server-crt", defServerCrt, "path to server crt")
	flag.Parse()
}

func (c *ServerConfig) PrintConfig() []byte {
	jsonConfig, _ := json.Marshal(c)
	return jsonConfig
}

func (c *ServerConfig) GetHost() string {
	return c.Host
}

func (c *ServerConfig) GetDatabaseURI() string {
	return c.DatabaseURI
}

func (c *ServerConfig) GetSecretKey() string {
	return c.JWTSecret
}

func (c *ServerConfig) GetCryptoKey() string {
	return c.CryptoKey
}

func (c *ServerConfig) GetServerKeyPath() string {
	return c.ServerKeyPath
}

func (c *ServerConfig) GetServerCrtPath() string {
	return c.ServerCrtPath
}
