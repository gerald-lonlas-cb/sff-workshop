package handler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.cbhq.net/engineering/sff-workshop/contract"
	"github.cbhq.net/engineering/sff-workshop/internal/config"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (h *TransactionHandler) getTransactOpts(
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

	fmt.Println(fmt.Sprintf("Retrieved nonce %d", nonce))

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

// ERC1155Transfer handles ERC1155 transfer that sends the the pre-minted tokens
func (h *TransactionHandler) ERC1155Transfer(
	ctx context.Context,
	to string,
	id int64,
	quantity int64,
) (string, error) {
	toAddr := common.HexToAddress(to)

	// Get the wallet address from which we are sending the NFT (OFFLINE)
	_, fromAddr, err := h.privateKeyAddress()
	if err != nil {
		return "", fmt.Errorf("error generating generating address: %v", err)
	}

	// Fill in some standard transaction options (e.g. gas limit, auth etc) (ONLINE)
	transactOpts, err := h.getTransactOpts(ctx)
	if err != nil {
		return "", fmt.Errorf("error generating transact opts: %v", err)
	}

	// Instantiate an instance of ERC1155 contract which defines our tokens and NFTs (OFFLINE)
	contractAddr := common.HexToAddress(h.cfg.ContractAddress)

	// Getting the ABI (OFFLINE)
	abi, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return "", fmt.Errorf("error getting ABI: %v", err)
	}

	// Generating the txData (OFFLINE)
	var data []byte = nil
	txData, err := abi.Pack(
		"safeTransferFrom",
		*fromAddr,
		toAddr,
		big.NewInt(id),
		big.NewInt(quantity),
		data,
	)
	if err != nil {
		return "", fmt.Errorf("error generating txData: %v", err)
	}

	// Create a legacy transaction using gas price (prior to London upgrade) (OFFLINE)
	nonce := transactOpts.Nonce.Uint64()
	baseTx := &types.LegacyTx{
		To:       &contractAddr,
		Nonce:    nonce,
		GasPrice: transactOpts.GasPrice,
		Gas:      transactOpts.GasLimit,
		Value:    transactOpts.Value,
		Data:     txData,
	}
	rawTx := types.NewTx(baseTx)

	// Sign the transaction (OFFLINE)
	signedTx, err := transactOpts.Signer(transactOpts.From, rawTx)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	// Submit transaction to Cloud Node (ONLINE)
	err = h.client.SendTransaction(transactOpts.Context, signedTx)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
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
