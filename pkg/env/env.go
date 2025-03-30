package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr     string
	KeyPath  string
	CertPath string
}

func NewConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("can't create config: %v", err)
	}

	cfg := &Config{
		Addr:     os.Getenv("ADDR"),
		KeyPath:  os.Getenv("KEY_PATH"),
		CertPath: os.Getenv("CERT_PATH"),
	}

	return cfg, nil
}
