package server

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glonlas/sff-aicp/contract"
	"github.com/glonlas/sff-aicp/internal/client"
	"github.com/glonlas/sff-aicp/internal/config"
	"github.com/glonlas/sff-aicp/internal/wallet"
)

type Server struct {
	context           context.Context
	mainWallet      wallet.MainWallet
	client            *ethclient.Client
	contractAddress   common.Address
	contractABI       *abi.ABI
}

func NewServer(port *int) (*Server, error) {
	ctx := context.Background()
	cfg, err := config.NewConfig(port)
	if err != nil {
		return nil, err
	}

	client, err := client.NewEVMClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mainWallet, err := wallet.NewMainWallet(cfg)
	if err != nil {
		return nil, err
	}

	// Get the Contract ABI
	contractAbi, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	s := &Server{
		context: ctx,
		mainWallet: mainWallet,
		client: client,
		contractAddress: common.HexToAddress(cfg.ContractAddress),
		contractABI: contractAbi,
	}

	return s, nil
}
