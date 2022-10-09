package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.cbhq.net/engineering/sff-workshop/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionHandler struct {
	cfg    *Config
	client *ethclient.Client
}

func newEthClient(ctx context.Context, cfg *Config) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(
		ctx,
		fmt.Sprintf("https://%s:%s@goerli.ethereum.coinbasecloud.net", cfg.Username, cfg.Password),
	)

	return client, err
}

func NewTransactionHandler(ctx context.Context, cfg *Config) (*TransactionHandler, error) {
	client, err := newEthClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &TransactionHandler{
		cfg:    cfg,
		client: client,
	}, nil
}

// constructUnsigned construct the unsigned transaction
func (h *TransactionHandler) Erc1155Transfer(
	ctx context.Context,
	to string,
	id int64,
	quantity int64,
) (string, error) {
	toAddr := common.HexToAddress(to)

	privateKey, fromAddr, err := h.privateKeyAddress()
	if err != nil {
		return "", fmt.Errorf("error getting private key: %v", err)
	}

	nonce, err := h.client.PendingNonceAt(ctx, *fromAddr)
	if err != nil {
		return "", fmt.Errorf("error getting nonce: %v", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("error suggesting gas price: %v", err)
	}

	chainID, err := h.client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("error creating transactor : %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	contractAddr := common.HexToAddress(h.cfg.ContractAddress)
	contractInstance, err := contract.NewContract(contractAddr, h.client)
	if err != nil {
		return "", fmt.Errorf("error loading contract: %v", err)
	}

	tx, err := contractInstance.SafeTransferFrom(
		auth,
		*fromAddr,
		toAddr,
		big.NewInt(id),
		big.NewInt(quantity),
		nil,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (h *TransactionHandler) privateKeyAddress() (*ecdsa.PrivateKey, *common.Address, error) {
	privateKey, err := crypto.HexToECDSA(h.cfg.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, &address, nil
}
