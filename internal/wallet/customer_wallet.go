package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type CustomerWallet interface {
	Sign(chainId *big.Int, unsignedTx *types.Transaction) (*types.Transaction, error)
	Address() *common.Address
	PrivateKey() *ecdsa.PrivateKey
}

type customerWallet struct {
	privateKey 		*ecdsa.PrivateKey
	address    		*common.Address
}

func NewCustomerWallet(mainWallet MainWallet, customerID string) (CustomerWallet, error) {
	// Convert customerID to integer
	index := StringToInt(customerID)
	log.Printf("converted customerID %s to %d", customerID, index)

	// Derive the Customer Wallet
	// This gives the path: m/44H/60H/0H/0/{index}
	// This provides 1 wallet per int {index} value
	walletRootPath := mainWallet.GetWalletRootPath()
	acc44H60H0H00, err := walletRootPath.Child(uint32(index))
	if err != nil {
		return nil, err
	}

	// Determine the Customer wallet private key
	btcecPrivKey, err := acc44H60H0H00.ECPrivKey()
	if err != nil {
		return nil, err
	}

	privateKey := btcecPrivKey.ToECDSA()
	
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("customer wallet address is %s", address)
	// It is strongly discouraged to uncomment the Private key log below. Use it only for test.
	//log.Printf("customer wallet private key is %s", PrivateKeyToHexString(privateKey))

	return &customerWallet {
		privateKey: privateKey,
		address:    &address,
	}, nil
}

// Sign the transaction (OFFLINE)
// This does not need to prompt customer to sign the message.
func (s *customerWallet) Sign(chainId *big.Int, unsignedTx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainId), s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	return signedTx, nil
}

func (s *customerWallet) Address() *common.Address {
	return s.address
}

func (s *customerWallet) PrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}

func StringToInt(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// /!\ For debug only /!\ 
// Use this function is you want to display the private of a specific wallet.
// So you can import it in a self-custody wallet
func PrivateKeyToHexString(privateKey *ecdsa.PrivateKey) string {
    return hex.EncodeToString(privateKey.D.Bytes())
}

