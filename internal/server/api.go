package server

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/glonlas/sff-aicp/contract"
	"github.com/glonlas/sff-aicp/internal/wallet"
)

// GetCustomerAvailableBalance handles the GET API request to retrieve the number (balance) of token the customer has.
// Required Parameters: customer_id, token_id
type GetCustomerAvailableBalanceResponse = struct {
    Balance *big.Int `json:"balance"`
}
func (s *Server) GetCustomerAvailableBalance(w http.ResponseWriter, r *http.Request) {
    log.Println("Received GetCustomerAvailableBalance request")

    // Define the required parameter names for this GET API
    requiredParams := []string{"customer_id", "token_id"}
    
    // Receive GET parameters and verify them
    query := r.URL.Query()

    // Verify all the required parameters are present
    err := s.verifyQueryParams(query, requiredParams)
    if err != nil {
        s.responseWithError(w, err)
        return
    }

    customerID := query.Get("customer_id")
    tokenID, err := s.getInt64(&query, "token_id")
    if err != nil {
        s.responseWithError(w, err)
        return
    }

    // Load customer wallet from customer ID
    customerWallet:= s.getCustomerWallet(customerID)
    customerWalletAddress := customerWallet.Address()

    // Call the smart contract
    instance, err := contract.NewContractCaller(s.contractAddress, s.client)
    if err != nil {
        s.responseWithError(w, err)
        return
    }

    balance, err := instance.BalanceOf(nil, *customerWalletAddress, big.NewInt(tokenID))
    if err != nil {
        s.responseWithError(w, err)
        return
    }

    // Send the response as JSON
    s.respondWithJSON(
        w, 
        http.StatusOK, 
        &GetCustomerAvailableBalanceResponse{
            Balance: balance,
        },
    )
}

// Generic response type for any call of the smart contract Write function 
type WriteTransactionResponse = struct {
	Transaction string `json:"transaction"`
}

// CreateOrder handles the GET API request to create an order.
// Required Parameters: customer_id, order_id, amount, merchant_address
func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received CreateOrder request")

	// Define the required parameter names for this GET API
	requiredParams := []string{"customer_id", "order_id", "amount", "merchant_address"}
	
	// Receive GET parameters and verify them
	query := r.URL.Query()

	// Verify all the required parameters are present
	err := s.verifyQueryParams(query, requiredParams)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	customerID := query.Get("customer_id")
	orderID := query.Get("order_id")
	merchantAddress := common.HexToAddress(query.Get("merchant_address"))
	amount, err := s.getInt64(&query, "amount")
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	// Load customer wallet from customer ID
	customerWallet:= s.getCustomerWallet(customerID)

	// Generating the txData (OFFLINE)
	txData, err := s.contractABI.Pack(
		"createOrder",
		orderID,
		big.NewInt(amount),
		merchantAddress,
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("error generating txData: %v", err))
		return
	}

	// Execute the signed transaction
	signedTx, err := s.constructAndSignTx(
		customerWallet,
		txData,
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("Failed to execute Sign transaction: %v", err))
		return
	}

	// Send the response as JSON
	s.respondWithJSON(
		w, 
		http.StatusOK, 
		&WriteTransactionResponse{
			Transaction: signedTx.Hash().Hex(),
		},
	)
}

// Refund handles the API request to refund an order.
// Required Parameters: customer_id, order_id
func (s *Server) CancelOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received CancelOrder request")
	
	// Define the required parameter names for this GET API
	requiredParams := []string{"customer_id", "order_id",}
	
	// Receive GET parameters and verify them
	query := r.URL.Query()

	// Verify all the required parameters are present
	err := s.verifyQueryParams(query, requiredParams)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	customerID := query.Get("customer_id")
	orderID := query.Get("order_id")

	// Load the customer wallet
	// Note, ideally No one except admins should be able to cancel an order. 
	// Given the current smart contract does not have any constraint of who can call it, 
	// Here we are using any wallet we control to execute that command.
	customerWallet:= s.getCustomerWallet(customerID)

	// Generating the txData (OFFLINE)
	txData, err := s.contractABI.Pack(
		"cancelOrder",
		orderID,
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("error generating txData: %v", err))
		return
	}

	// Execute the signed transaction
	signedTx, err := s.constructAndSignTx(
		customerWallet,
		txData,
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("Failed to execute Sign transaction: %v", err))
		return
	}

	// Send the response as JSON
	s.respondWithJSON(
		w, 
		http.StatusOK, 
		&WriteTransactionResponse{
			Transaction: signedTx.Hash().Hex(),
		},
	)
}

// GetOrderStatus returns the status of a specific order
// Required Parameters: order_id
type GetOrderStatusResponse = struct {
    Status uint8 `json:"status"`
}
func (s *Server) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
    log.Println("Received GetOrderStatus request")

    	// Define the required parameter names for this GET API
	requiredParams := []string{"order_id",}
	
	// Receive GET parameters and verify them
	query := r.URL.Query()

	// Verify all the required parameters are present
	err := s.verifyQueryParams(query, requiredParams)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	orderID := query.Get("order_id")

	// Call the smart contract
	instance, err := contract.NewContractCaller(s.contractAddress, s.client)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	status, err := instance.GetOrderStatus(nil, orderID)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	// Send the response as JSON
	s.respondWithJSON(
		w, 
		http.StatusOK, 
		&GetOrderStatusResponse{
			Status: status,
		},
	)
}

// GetTransactionStatus returns the status of a on chain transaction.
// You can use it to know if the transaction is complete
// Required Parameters: transaction_hash
type GetTransactionStatusResponse = struct {
    Receipt *types.Receipt `json:"Receipt"`
}
func (s *Server) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
    log.Println("Received GetTransactionStatus request")

    // Define the required parameter names for this GET API
    requiredParams := []string{"transaction_hash"}
    
    // Receive GET parameters and verify them
    query := r.URL.Query()

    // Verify all the required parameters are present
    err := s.verifyQueryParams(query, requiredParams)
    if err != nil {
        s.responseWithError(w, err)
        return
    }

    transactionHash := common.HexToHash(query.Get("transaction_hash"))

	receipt, err := s.client.TransactionReceipt(s.context, transactionHash)
	if err != nil {
        s.responseWithError(w, err)
        return
    }

    // Send the response as JSON
    s.respondWithJSON(
        w, 
        http.StatusOK, 
        &GetTransactionStatusResponse{
            Receipt: receipt,
        },
    )
}


// Allows you to create and send to the customer XXX of the ERC-1155 tokens uses a voucher
// Required Parameters: customer_id, token_id, quantity
func (s *Server) MintTokens(w http.ResponseWriter, r *http.Request) {
	log.Println("Received MintTokens request")
	
	// Define the required parameter names for this GET API
	requiredParams := []string{"customer_id", "quantity",}
	
	// Receive GET parameters and verify them
	query := r.URL.Query()

	// Verify all the required parameters are present
	err := s.verifyQueryParams(query, requiredParams)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	customerID := query.Get("customer_id")
	tokenID, err := s.getInt64(&query, "token_id")
	if err != nil {
		s.responseWithError(w, err)
		return
	}
	quantity, err := s.getInt64(&query, "quantity")
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	// Load the customer wallet
	// Note, ideally No one except admins should be able to mint tokens
	// But this temporary smart contract allows anyone to mint the tokens 
	// we will use this function to fund our customer's wallets when needed
	customerWallet:= s.getCustomerWallet(customerID)

	// Generating the txData (OFFLINE)
	txData, err := s.contractABI.Pack(
		"mint",
		customerWallet.Address(),
		big.NewInt(tokenID),
		big.NewInt(quantity),
		[]byte{},
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("error generating txData: %v", err))
		return
	}

	// Execute the signed transaction
	signedTx, err := s.constructAndSignTx(
		customerWallet,
		txData,
	)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("Failed to execute Sign transaction: %v", err))
		return
	}

	// Send the response as JSON
	s.respondWithJSON(
		w, 
		http.StatusOK, 
		&WriteTransactionResponse{
			Transaction: signedTx.Hash().Hex(),
		},
	)
}

// GetCustomerPublicAddress handles the GET API request to display the public wallet address of a CustomerID.
// It does not interact with the blockchain. Only compute the wallet from your private key.
// Required Parameters: customer_id
type GetCustomerPublicAddressResponse = struct {
	PublicAddress *common.Address `json:"public_address"`
}
func (s *Server) GetCustomerPublicAddress(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GetCustomerPublicAddress request")

	// Define the required parameter names for this GET API
	requiredParams := []string{"customer_id"}
	
	// Receive GET parameters and verify them
	query := r.URL.Query()

	// Verify all the required parameters are present
	err := s.verifyQueryParams(query, requiredParams)
	if err != nil {
		s.responseWithError(w, err)
		return
	}

	customerID := query.Get("customer_id")

	// Load customer wallet from customer ID
	customerWallet:= s.getCustomerWallet(customerID)
	customerWalletAddress := customerWallet.Address()

	// Send the response as JSON
	s.respondWithJSON(
		w, 
		http.StatusOK, 
		&GetCustomerPublicAddressResponse{
			PublicAddress: customerWalletAddress,
		},
	)
}

func (s *Server) getCustomerWallet(customerID string) wallet.CustomerWallet {
	customerWallet, err := wallet.NewCustomerWallet(s.mainWallet, customerID)
	if err != nil {
		log.Printf("Failed to load customer wallet: %v", err)
		return nil
	}

	return customerWallet
}

func (s *Server) getInt64(query *url.Values, field string) (int64, error) {
	val := query.Get(field)
	return strconv.ParseInt(val, 10, 64)
}
