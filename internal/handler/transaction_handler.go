package handler

import (
	"context"
	"fmt"
	"math/big"

	"github.cbhq.net/engineering/sff-workshop/contract"
	"github.cbhq.net/engineering/sff-workshop/internal/config"
	"github.cbhq.net/engineering/sff-workshop/internal/keystore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionHandler struct {
	cfg    *config.Config
	client *ethclient.Client
	signer keystore.Signer
}

func NewTransactionHandler(
	ctx context.Context,
	ethClient *ethclient.Client,
	cfg *config.Config,
	signer keystore.Signer,
) (*TransactionHandler, error) {
	return &TransactionHandler{
		cfg:    cfg,
		client: ethClient,
		signer: signer,
	}, nil
}

// ERC1155Transfer handles ERC1155 transfer that sends the pre-minted tokens
func (h *TransactionHandler) ERC1155Transfer(
	ctx context.Context,
	to string,
	id int64,
	quantity int64,
) (string, error) {
	unsignedTx, err := h.constructUnsignedTx(ctx, to, id, quantity)
	if err != nil {
		return "", fmt.Errorf("error constructing transaction: %v", err)
	}

	signedTx, err := h.signTx(ctx, unsignedTx)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	// Submit transaction to Cloud Node (ONLINE)
	err = h.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// constructUnsignedTx takes in input params and construct a raw unsigned transaction
func (h *TransactionHandler) constructUnsignedTx(
	ctx context.Context,
	to string,
	id int64,
	quantity int64,
) (*types.Transaction, error) {
	fromAddr := *h.signer.Address()
	toAddr := common.HexToAddress(to)
	contractAddr := common.HexToAddress(h.cfg.ContractAddress)

	// Retrieve nonce for fromAddress (ONLINE)
	nonce, err := h.client.NonceAt(ctx, fromAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting nonce: %v", err)
	}

	// Estimate Gas Price (ONLINE)
	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error suggesting gas price: %v", err)
	}

	// Getting the Contract ABI (OFFLINE)
	contractAbi, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("error getting ABI: %v", err)
	}

	// Generating the txData (OFFLINE)
	var data []byte = nil
	txData, err := contractAbi.Pack(
		"safeTransferFrom",
		fromAddr,
		toAddr,
		big.NewInt(id),
		big.NewInt(quantity),
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating txData: %v", err)
	}

	// Construct Transaction (OFFLINE)
	baseTx := &types.LegacyTx{
		To:       &contractAddr,
		Nonce:    nonce,
		GasPrice: gasPrice,       // in wei
		Gas:      uint64(300000), // in unit
		Value:    big.NewInt(0),
		Data:     txData,
	}
	unsignedTx := types.NewTx(baseTx)

	return unsignedTx, nil
}

// signTx signs the unsigned transaction using the signer
func (h *TransactionHandler) signTx(
	ctx context.Context,
	unsignedTx *types.Transaction,
) (*types.Transaction, error) {
	// Getting ChainID (ONLINE)
	chainId, err := h.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting ChainID")
	}

	signedTx, err := h.signer.Sign(chainId, unsignedTx)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	return signedTx, nil
}
