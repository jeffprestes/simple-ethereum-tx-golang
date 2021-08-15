package main

import (
	"context"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jeffprestes/goethereumhelper"
)

// Example of Go Ethereum Helper library usage to send ether
// but calculating fees using EIP-1559
func main() {
	// Example: https://mainnet.infura.io/v3/0f0f0f0f0ff00ff
	ethereumClientURL := os.Getenv("ETH_CLIENT_URL")
	if len(ethereumClientURL) < 10 {
		log.Fatalln("Invalid Ethereum client URL")
	}

	pvtKeyHex := os.Getenv("PVT_KEY")
	if len(pvtKeyHex) < 64 {
		log.Fatalln("Invalid HEX private key")
	}

	senderPrivateKey, err := crypto.HexToECDSA(pvtKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %s", err.Error())
	}

	client, err := goethereumhelper.GetCustomNetworkClient(ethereumClientURL)
	if err != nil {
		log.Fatalf("Error connecting to Ethereum client: %s", err.Error())
	}

	to := common.HexToAddress("0x263C3Ab7E4832eDF623fBdD66ACee71c028Ff591")
	value := int64(100000000000000000)

	tx, err := goethereumhelper.SendEther(client, senderPrivateKey, to, value)
	if err != nil {
		log.Fatalf("Error submitting transaction: %s", err.Error())
	}

	log.Printf("Transaction submitted: %s\n\n", tx.Hash())
	log.Println("Waiting for the transaction being mined...")
	txReceipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Error mining transaction: %s", err.Error())
	}

	log.Printf("Transaction processed: %+v", txReceipt)
}
