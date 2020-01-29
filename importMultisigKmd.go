package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

const kmdAddress = "http://localhost:7833"
const kmdToken = "e4a0aa0518da921bae3f27f40f443a4e154a1c80f25f7759db7767e04d73c06c"

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}

	// Get the list of wallets
	listResponse, err := kmdClient.ListWallets()
	if err != nil {
		fmt.Printf("error listing wallets: %s\n", err)
		return
	}

	// Find our wallet name in the list
	var exampleWalletID string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == "wallet22" {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			exampleWalletID = wallet.ID
		}
	}

	// Get a wallet handle
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "testpassword")
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// create account1
	account1 := crypto.GenerateAccount()
	fmt.Println(account1.Address)
	fmt.Println(account1.PrivateKey)

	// create account1
	account2 := crypto.GenerateAccount()
	fmt.Println(account2.Address)
	fmt.Println(account2.PrivateKey)

	// create msig object
	msig, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		account1.Address,
		account2.Address,
	})
	if err != nil {
		panic("invalid multisig parameters")
	}

	// importedMultiSig, err := kmdClient.ImportMultisig(exampleWalletHandleToken, 1, 2, msig)
	importedMultiSig, err := kmdClient.ImportMultisig(exampleWalletHandleToken, 1, 2, account1.Address)
	fmt.Println(importedMultiSig)
}
