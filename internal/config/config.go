package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Username        string
	Password        string
	PrivateKey      string
	ContractAddress string
	Mnemonic        string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	return &Config{
		Username:        os.Getenv("USERNAME"),
		Password:        os.Getenv("PASSWORD"),
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		Mnemonic:        os.Getenv("MNEMONIC"),
	}, nil
}
