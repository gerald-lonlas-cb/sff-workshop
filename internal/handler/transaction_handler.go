package handler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"math/big"

	"github.cbhq.net/engineering/sff-workshop/contract"
	"github.cbhq.net/engineering/sff-workshop/internal/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip39"
)

type TransactionHandler struct {
	cfg    *config.Config
	client *ethclient.Client
}

func NewTransactionHandler(ctx context.Context, ethClient *ethclient.Client, cfg *config.Config) (*TransactionHandler, error) {
	return &TransactionHandler{
		cfg:    cfg,
		client: ethClient,
	}, nil
}

func (h *TransactionHandler) GetTransactOpts(
	ctx context.Context,
) (*bind.TransactOpts, error) {
	privateKey, fromAddr, err := h.privateKeyAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting private key: %v", err)
	}

	nonce, err := h.client.NonceAt(ctx, *fromAddr, nil)
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
	seed := bip39.NewSeed(h.cfg.Mnemonic, "")

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting master key: %v", err)
	}

	// This gives the path: m/44H
	acc44H, err := masterKey.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, nil, err
	}

	// This gives the path: m/44H/60H
	acc44H60H, err := acc44H.Child(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		return nil, nil, err
	}

	// This gives the path: m/44H/60H/0H
	acc44H60H0H, err := acc44H60H.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, nil, err
	}

	// This gives the path: m/44H/60H/0H/0
	acc44H60H0H0, err := acc44H60H0H.Child(0)
	if err != nil {
		return nil, nil, err
	}

	// This gives the path: m/44H/60H/0H/0/0
	acc44H60H0H00, err := acc44H60H0H0.Child(0)
	if err != nil {
		return nil, nil, err
	}

	btcecPrivKey, err := acc44H60H0H00.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	privateKey := btcecPrivKey.ToECDSA()

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, &address, nil
}
