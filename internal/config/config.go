package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NodeURI                 string
	ContractAddress         string
	Mnemonic                string
}

func NewConfig(port *int) (*Config, error) {
	if *port != -1 {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		NodeURI:             os.Getenv("NODE_URI"),
		ContractAddress:     os.Getenv("CONTRACT_ADDRESS"),
		Mnemonic:            os.Getenv("MNEMONIC"),
	}, nil
}
