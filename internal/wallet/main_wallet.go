package wallet

import (
	"fmt"

	"github.com/glonlas/sff-aicp/internal/config"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

type MainWallet interface {
	GetWalletRootPath() *hdkeychain.ExtendedKey
}

type mainWallet struct {
	walletRootPath 	*hdkeychain.ExtendedKey
}

func NewMainWallet(cfg *config.Config) (MainWallet, error) {
	seed := bip39.NewSeed(cfg.Mnemonic, "")

	mainKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("error getting master key: %v", err)
	}

	// Drill down to the HD wallet level in which we can assign a wallet per customer
	// This gives the BIP32 Derivation path: m/44H
	acc44H, err := mainKey.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	// This gives the BIP32 Derivation path: m/44H/60H
	acc44H60H, err := acc44H.Child(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		return nil, err
	}

	// This gives the BIP32 Derivation path: m/44H/60H/0H
	acc44H60H0H, err := acc44H60H.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	// This gives the BIP32 Derivation path: m/44H/60H/0H/0
	acc44H60H0H0, err := acc44H60H0H.Child(0)
	if err != nil {
		return nil, err
	}

	return &mainWallet{
		walletRootPath: acc44H60H0H0,
	}, nil
}

func (mw *mainWallet) GetWalletRootPath() *hdkeychain.ExtendedKey {
	return mw.walletRootPath
}
