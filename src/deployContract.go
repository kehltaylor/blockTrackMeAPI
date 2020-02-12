package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/ethclient"

)


func DeployContract(w http.ResponseWriter, r *http.Request) {


		/*client, err := ethclient.Dial("https://rinkeby.infura.io")
		if err != nil {
			log.Fatal(err)
		}*/
		fmt.Println("Inside function")

		client := ReturnClient()

		//privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
		privateKey, err := crypto.HexToECDSA("565C207EBBF93700F788D0F2BF79C1F68D5EBC307995E53E6E7EF8C72E6DDD25")

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Read Private key")
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}
		fmt.Println("Public key good")

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal("Nounce shit went wrong")
			log.Fatal(err)
		}
		fmt.Println("Nounce good")

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth := bind.NewKeyedTransactor(privateKey)
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(1000)     // in wei
		auth.GasLimit = uint64(3000000) // in units
		auth.GasPrice = gasPrice

		address, tx, instance, err := DeployTheContract(auth, &client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
		fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

		SetContractAddress(address.Hex())
		var theContract ERC20Contract
		theContract.addr = address.Hex()
		res, err := json.Marshal(theContract)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = instance
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
}