package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Username        string
	Password        string
	NodeURI         string
	Mnemonic        string
	ContractAddress string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	return &Config{
		Username:        os.Getenv("USERNAME"),
		Password:        os.Getenv("PASSWORD"),
		NodeURI:         os.Getenv("NODE_URI"),
		Mnemonic:        os.Getenv("MNEMONIC"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
	}, nil
}
