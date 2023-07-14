package server

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/glonlas/sff-aicp/internal/wallet"
)

// Use this method to create a transaction for a Write smart contract function
// txData is the data you get from a contractAbi.Pack()
func (s *Server) constructAndSignTx(
	customerWallet wallet.CustomerWallet,
	txData []byte,
) (*types.Transaction, error) {
	
	// Prepare unsigned content
	unsignedTx, err := s.constructUnsignedTx(
		customerWallet.Address(),
		&s.contractAddress,
		txData,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare Unsigned content: %v", err)
	}
	
	signedTx, err := s.signTx(customerWallet, unsignedTx)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	// Submit the transaction on chain
	err = s.client.SendTransaction(s.context, signedTx)
	if err != nil {
		return nil, fmt.Errorf("error submitting transaction: %v", err)
	}

	return signedTx, nil
}

func (s *Server) constructUnsignedTx(
	customerWalletAddress *common.Address,
	smartContractAddress *common.Address,
	txData []byte,
) (*types.Transaction, error) {
	// Get Wallet nonce
	nonce, err := s.GetCustomerWalletNonce(customerWalletAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to load customer wallet nonce: %v", err)
	}

	// Get Gas Price
	gasPrice, err := s.GetGasPrice()
	if err != nil {
		return nil, fmt.Errorf("Failed to load gas price: %v", err)
	}

	// Construct Transaction (OFFLINE)
	baseTx := &types.LegacyTx{
		To:       &s.contractAddress,
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
// The signer here represent a customer wallet.
//Ensure the given customer wallet used, has enough utility tokens (ex: Matic) to pay the gas fee
func (s *Server) signTx(
	customerWallet wallet.CustomerWallet,
	unsignedTx *types.Transaction,
) (*types.Transaction, error) {
	// Getting ChainID (ONLINE)
	chainId, err := s.client.ChainID(s.context)
	if err != nil {
		return nil, fmt.Errorf("error getting ChainID")
	}

	signedTx, err := customerWallet.Sign(chainId, unsignedTx)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	return signedTx, nil
}

func (s *Server) GetCustomerWalletNonce(walletAddress *common.Address) (uint64, error) {
	// Retrieve nonce for fromAddress (ONLINE)
	nonce, err := s.client.PendingNonceAt(s.context, *walletAddress)
	if err != nil {
		return 0, fmt.Errorf("error getting nonce: %v", err)
	}

	log.Printf("Retrieved nonce %d", nonce)

	return nonce, nil
}

func (s *Server) GetGasPrice() (*big.Int, error) {
	// Estimate Gas Price (ONLINE)
	suggestedGasPrice, err := s.client.SuggestGasPrice(s.context)
	if err != nil {
		return nil, fmt.Errorf("error suggesting gas price: %v", err)
	}
	gasPrice := big.NewInt(int64(float64(suggestedGasPrice.Int64()) * 1.5))

	log.Printf("Suggested gas %d, used gas %d", suggestedGasPrice.Int64(), gasPrice.Int64())

	return gasPrice, nil
}

// Function to sign the message used for the redeem code
func (s *Server) SignMessage() error {
	return nil
}
