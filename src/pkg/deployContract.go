package features

import (
	"../utils"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/ethclient"

	store "../Contracts" // for demo
)

func DeployContract(w http.ResponseWriter, r *http.Request) {


		/*client, err := ethclient.Dial("https://rinkeby.infura.io")
		if err != nil {
			log.Fatal(err)
		}*/

		client := utils.ReturnClient()

		//privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
		privateKey, err := crypto.HexToECDSA("D62B4115FE4E6063282DE99D9C50A6C7A3E66F75FA1E2A5014C3BA8B93160687")

		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth := bind.NewKeyedTransactor(privateKey)
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(1000)     // in wei
		auth.GasLimit = uint64(3000000) // in units
		auth.GasPrice = gasPrice

		address, tx, instance, err := store.DeployContract(auth, &client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
		fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

		utils.SetContractAddress(address.Hex())

		_ = instance
}