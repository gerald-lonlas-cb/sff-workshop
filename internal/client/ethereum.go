package client

import (
	"context"
	"fmt"

	"github.cbhq.net/engineering/sff-workshop/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

// NewEthClient returns an Ethereum client that uses a CoinbaseCloud Node
func NewEthClient(ctx context.Context, cfg *config.Config) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(
		ctx,
		fmt.Sprintf("https://%s:%s@goerli.ethereum.coinbasecloud.net", cfg.Username, cfg.Password),
	)

	return client, err
}
