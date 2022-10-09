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

// TODO: move to .env file
const privateKeyString = "7473e2f801b7a5f50fd96a049375e36551ad52d97ef8171d794627e941c7308a"
const contractAddressString = "0xEB20BcB4689FA795443ddA3a5360282AFF7bFF75"

type TransactionHandler struct {
	ctx    context.Context
	client *ethclient.Client
}

func newTransactionHandler(_ *ethclient.Client, ctx context.Context) TransactionHandler {
	// TODO: use passed in ETH client
	client, _ := ethClient()
	return TransactionHandler{
		ctx:    ctx,
		client: client,
	}
}

// constructUnsigned construct the unsigned transaction
func (h TransactionHandler) erc1155Transfer(
	to string,
	id int64,
	quantity int64,
) (string, error) {
	toAddr := common.HexToAddress(to)

	privateKey, fromAddr, err := privateKeyAddress()
	if err != nil {
		return "", fmt.Errorf("error getting private key: %v", err)
	}

	nonce, err := h.client.PendingNonceAt(h.ctx, *fromAddr)
	if err != nil {
		return "", fmt.Errorf("error getting nonce: %v", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(h.ctx)
	if err != nil {
		return "", fmt.Errorf("error suggesting gas price: %v", err)
	}

	chainID, err := h.client.ChainID(h.ctx)
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

	contractAddr := common.HexToAddress(contractAddressString)
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

	return tx.Hash().Hex(), nil
}

func privateKeyAddress() (*ecdsa.PrivateKey, *common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
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

func ethClient() (*ethclient.Client, error) {
	ctx := context.Background()
	client, err := ethclient.DialContext(
		ctx,
		"https://goerli.infura.io/v3/212009411b2846588674ff677abf0fa5",
		//"https://MTSJPQSCJQU6YWKIGJVG:S3E4XACGWSK26DORYXVUASRSKI2GSZBQQ7CKMKO2@https://goerli.ethereum.coinbasecloud.net",
	)

	return client, err
}
