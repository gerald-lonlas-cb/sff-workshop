package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.cbhq.net/engineering/sff-workshop/contract"
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

func (h *TransactionHandler) GetTransactOpts(
	ctx context.Context,
) (*bind.TransactOpts, error) {
	privateKey, fromAddr, err := h.privateKeyAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting private key: %v", err)
	}

	nonce, err := h.client.PendingNonceAt(ctx, *fromAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting nonce: %v", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error suggesting gas price: %v", err)
	}

	chainID, err := h.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting chain id: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("error creating transactor : %v", err)
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0)     // in wei
	opts.GasLimit = uint64(300000) // in units
	opts.GasPrice = gasPrice

	return opts, nil
}

// constructUnsigned construct the unsigned transaction
func (h *TransactionHandler) Erc1155Transfer(
	ctx context.Context,
	to string,
	id int64,
	quantity int64,
) (string, error) {
	toAddr := common.HexToAddress(to)

	// Get the wallet address from which we are sending the NFT
	_, fromAddr, err := h.privateKeyAddress()

	// Fill in some standard transaction options (e.g. gas limit, auth etc)
	transactOpts, err := h.GetTransactOpts(ctx)
	if err != nil {
		return "", fmt.Errorf("error generating transact opts", err)
	}

	// Instantiate an instance of ERC1155 contract which defines our tokens and NFTs
	contractAddr := common.HexToAddress(h.cfg.ContractAddress)
	contractInstance, err := contract.NewContract(contractAddr, h.client)
	if err != nil {
		return "", fmt.Errorf("error loading contract: %v", err)
	}

	// Send the "Golden Badge" to the user's wallet
	tx, err := contractInstance.SafeTransferFrom(
		transactOpts,
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
