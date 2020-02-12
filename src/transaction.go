package main

//https://goethereumbook.org/en/transfer-tokens/

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"


	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/ethereum/go-ethereum/ethclient"
	"net/http"
)

type Transaction struct{
	senderPrivKey	string 	`json:"senderPrivKey"`
	receiverPubKey	string	`json:"receiverPubKey"`
}



func sendTransaction(w http.ResponseWriter, r *http.Request){

	client, error := ethclient.Dial("https://ropsten.infura.io/v3/511162a74a0c4a80a9fbab7b9d2718b8")
	if error != nil {
		log.Fatal(error)
	}
	_ = client // we'll use this in the upcoming sections

	var trans Transaction

	err := json.NewDecoder(r.Body).Decode(&trans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println("Addrs priv")
	fmt.Println(trans.senderPrivKey)
	//privateKey, err := crypto.HexToECDSA("E6F9A2469E33F7808666ED49CB88C4AB6637E08AE6ADD0CFAC4681CB0D87B3F7")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//publicKey := privateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	//}
	//fmt.Println("Addrs pub")
	//fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//value := big.NewInt(0) // in wei (0 eth)
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//toAddress := common.HexToAddress("0x634222A6A7Fa4eFB4744cac8A1970A910f0373E2")
	//tokenAddress := common.HexToAddress("0xDE01C179E980b5e5939AE3246400E4bD646bE2C1")
	//
	//fmt.Println("Addrs")
	//transferFnSignature := []byte("transfer(address,uint256)")
	//hash := sha3.NewLegacyKeccak256()
	//hash.Write(transferFnSignature)
	//methodID := hash.Sum(nil)[:4]
	//fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb
	//
	//paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
	//
	//amount := new(big.Int)
	//amount.SetString("1000000000000000000", 10) // sets the value to 1 tokens, in the token denomination
	//
	//paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	//
	//var data []byte
	//data = append(data, methodID...)
	//data = append(data, paddedAddress...)
	//data = append(data, paddedAmount...)
	//
	//gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &tokenAddress,
	//	Data: data,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(gasLimit) // 23256
	//
	//tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	//
	//chainID, err := client.NetworkID(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Sign begin")
	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Signing")
	//err = client.SendTransaction(context.Background(), signedTx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(signedTx.Hash().Hex())
	//fmt.Printf("tx sent: %s", signedTx.Hash().Hex())



	privateKey, err := crypto.HexToECDSA("565C207EBBF93700F788D0F2BF79C1F68D5EBC307995E53E6E7EF8C72E6DDD25")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xe199aB06903C2a1F812dD90d3FA921ca390555c0")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signedTx.Hash().Hex())
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())



	/*value := big.NewInt(0)
	toAddress := common.HexToAddress(_toAddress)
	tokenAddress := common.HexToAddress(_tokenAddress)
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	client := returnClient()

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)*/




	//client := utils.ReturnClient()
	//
	//tokenAddress := common.HexToAddress(utils.GetContractAddress())
	//instance, err := token.NewToken(tokenAddress, client)
	/*privateKey, err := crypto.HexToECDSA(_hexPK)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(_toAddress)
	tokenAddress := common.HexToAddress(_tokenAddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // sets the value to 1000 tokens, in the token denomination

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
*/
	 
}