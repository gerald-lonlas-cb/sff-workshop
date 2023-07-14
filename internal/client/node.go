package client

import (
	"context"
	"fmt"
	"log"

	"github.com/glonlas/sff-aicp/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

// NewEVMClient returns an EVM client that uses a CoinbaseCloud Node
func NewEVMClient(ctx context.Context, cfg *config.Config) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(
		ctx,
		fmt.Sprintf("https://%s", cfg.NodeURI),
	)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return nil, err
	}

	fmt.Println("Connected to EVM node successfully!")

	return client, err
}
