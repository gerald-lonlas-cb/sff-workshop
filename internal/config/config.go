package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Username        string
	Password        string
	Mnemonic        string
	ContractAddress string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	return &Config{
		Username:        os.Getenv("USERNAME"),
		Password:        os.Getenv("PASSWORD"),
		Mnemonic:        os.Getenv("MNEMONIC"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
	}, nil
}
