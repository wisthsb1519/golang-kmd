package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

// These constants represent the kmd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
const kmdAddress = "http://localhost:7833"
const kmdToken = "e4a0aa0518da921bae3f27f40f443a4e154a1c80f25f7759db7767e04d73c06c"

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	fmt.Println("Made a kmd client")

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
	// account, err := kmdClient.ListKeys(exampleWalletHandleToken)
	// fmt.Printf("Accounts", account)
	// Extract the account sk
	accountKeyResponse, err := kmdClient.ExportKey(exampleWalletHandleToken, "testpassword", "EMJ637L2UXZGGZLP4MTZST6R626XU3HAN264MUVGKKUH2OOCPNYCYK6MNU")
	accountKey := accountKeyResponse.PrivateKey
	fmt.Printf("Account Key Response: %v ", accountKey, "\n ")
	fmt.Printf("Account Key: \n", accountKey, "\n ")

	fmt.Printf("length of the key is %v \n", len(accountKey))
	// Convert sk to mnemonic
	mn, err := mnemonic.FromPrivateKey(accountKey)
	if err != nil {
		fmt.Printf("Error getting backup phrase: %s\n", err)
		return
	}
	fmt.Printf(mn)

}
