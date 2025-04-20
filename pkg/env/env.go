package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type ServerConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type ProxyConfig struct {
	Addr     string
	KeyPath  string
	CertPath string
}
type Config struct {
	ProxyConfig  ProxyConfig
	MongoConfig  MongoConfig
	ServerConfig ServerConfig
}

func NewConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("can't create config: %v", err)
	}

	cfg := &Config{
		ProxyConfig: ProxyConfig{
			Addr:     os.Getenv("PROXY_ADDR"),
			KeyPath:  os.Getenv("PROXY_KEY_PATH"),
			CertPath: os.Getenv("PROXY_CERT_PATH"),
		},
		ServerConfig: ServerConfig{
			Addr:            os.Getenv("SERVER_ADDR"),
			ReadTimeout:     10,
			WriteTimeout:    10,
			IdleTimeout:     100,
			ShutdownTimeout: 10,
		},
		MongoConfig: MongoConfig{
			Host:     os.Getenv("MONGO_HOST"),
			Port:     os.Getenv("MONGO_PORT"),
			Database: os.Getenv("MONGO_DATABASE"),
		},
	}

	return cfg, nil
}
