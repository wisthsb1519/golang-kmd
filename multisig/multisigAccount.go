package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// Change these values for your own algod.token and algod.net values
const algodAddress = "http://127.0.0.1:8080"
const algodToken = "a967f42b017cd4c5c95a633e87b5ff14226ae60609e174bf5832722631946e13"

func main() {

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %v\n", err)
		return
	}

	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	// Generate Accounts
	acct1 := crypto.GenerateAccount()
	acct2 := crypto.GenerateAccount()
	acct3 := crypto.GenerateAccount()

	// Decode the account addresses
	addr1, _ := types.DecodeAddress(acct1.Address.String())
	addr2, _ := types.DecodeAddress(acct2.Address.String())
	addr3, _ := types.DecodeAddress(acct3.Address.String())

	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})
	if err != nil {
		panic("invalid multisig parameters")
	}

	// declare txn parameters
	fee := txParams.Fee
	firstRound := txParams.LastRound
	lastRound := txParams.LastRound + 1000
	genesisID := txParams.GenesisID     // replace me
	genesisHash := txParams.GenesisHash // replace me
	const amount1 = 2000
	const amount2 = 1500
	var note []byte
	closeRemainderTo := ""

	fromAddr, _ := ma.Address()
	toAddr := "WICXIYCKG672UGFCCUPBAJ7UYZ2X7GZCNBLSAPBXW7M6DZJ5YY6SCXML4A"

	// Create the transaction
	txn, err := transaction.MakePaymentTxn(fromAddr.String(), toAddr, fee, amount1, firstRound, lastRound, note, closeRemainderTo, genesisID, genesisHash)

	// First signature on PST
	txid, preStxBytes, err := crypto.SignMultisigTransaction(acct1.PrivateKey, ma, txn)
	if err != nil {
		panic("could not sign multisig transaction")
	}
	fmt.Printf("Made partially-signed multisig transaction with TxID %s \n", txid)

	// Second signature on PST
	txid2, stxBytes, err := crypto.AppendMultisigTransaction(acct2.PrivateKey, ma, preStxBytes)
	if err != nil {
		panic("could not sign multisig transaction")
	}
	fmt.Printf("Made partially-signed multisig transaction with TxID %s \n", txid2)

	// Print multisig account
	fmt.Printf("Here is your multisig address : %s \n", fromAddr.String())

	fmt.Println("Please go to: https://bank.testnet.algorand.network/ to fund your multisig account.")
	fmt.Scanln() // wait for Enter Key

	// Send transaction to the network
	sendResponse, err := algodClient.SendRawTransaction(stxBytes)
	if err != nil {
		fmt.Printf("Failed to create payment transaction: %v\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
}
