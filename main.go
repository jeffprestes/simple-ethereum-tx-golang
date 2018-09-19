package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		fmt.Printf("Houve falha ao conectar com Rinkeby via infuria: %+v", err)
	}
	pvtkey, err := crypto.HexToECDSA("yourprivatekey")
	if err != nil {
		fmt.Printf("Houve falha ao gerar a chave privada: %+v", err)
	}
	pubkey := pvtkey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Printf("Houve falha fazer o casting da chave publica para o padrão ECDSA")
	}
	origem := crypto.PubkeyToAddress(*pubkeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), origem)
	if err != nil {
		fmt.Printf("Houve falha ao gerar o nonce da conta na rede: %+v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("Houve falha ao obter o preco sugerido de gas da rede: %+v", err)
	}

	auth := bind.NewKeyedTransactor(pvtkey)
	auth.GasLimit = uint64(30000000)
	//
	//  Increase here the gas price suggested by Ethereum network in order to have your transaction mined faster
	//
	auth.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	auth.Value = big.NewInt(0)
	//
	// If you want to resend the transaction with higher gas price do not change the previous nonce number
	//
	auth.Nonce = big.NewInt(int64(nonce))

	key, _ := crypto.GenerateKey()
	destinationAccount := bind.NewKeyedTransactor(key)
	//Carrega conta da empresa com weis

	tx := types.NewTransaction(uint64(auth.Nonce.Int64()), destinationAccount.From, big.NewInt(101000000000000), auth.GasLimit, auth.GasPrice, nil)
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, tx)
	if err != nil {
		fmt.Printf("Houve falha ao assinar a transação para enviar ether a nova conta de um comprador de registro: %+v", err)
	}
	client.SendTransaction(context.Background(), signedTx)
}
