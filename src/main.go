package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"math/big"
	"net/http"
	"time"
)

var DB *gorm.DB
var err error

type ERC20Contract struct {
	addr  	string `json:"addr"`

}

type ClientConnection struct {
	Status    string
}


type User struct {
	gorm.Model
	FistName  	string `json:"firstName"`
	LastName 	string `json:"lastName"`
	PubKey		string `json:"pubKey"`
}

func main() {
	_, err := dbConnect()
	if err != nil {
		log.Println("connection to DB failed, aborting...")
		log.Fatal(err)
	}

	log.Println("connected to DB")
	r := mux.NewRouter()
	r.HandleFunc("/", connection).Methods("POST")
	r.HandleFunc("/showUserHandler", showUserHandler).Methods("POST")
	r.HandleFunc("/createUserHeader", createUserHandler).Methods("POST")
	r.HandleFunc("/deployContract", DeployContract).Methods("POST")
	r.HandleFunc("/transaction", sendTransaction).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func dbConnect()(*gorm.DB, error){

	connectionParams := "dbname=goland user=postgres password=root sslmode=disable host=postgres"
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", connectionParams) // gorm checks Ping on Open
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return DB, err
	}

	DB.AutoMigrate(&User{})

	testPost := User{FistName: "Kojak", LastName: "Zoda", PubKey: "SexyBoy77"}
	DB.Create(&testPost)

	return DB, err
}

func connection(w http.ResponseWriter, r *http.Request) {

	client, err := ethclient.Dial("https://ropsten.infura.io/v3/511162a74a0c4a80a9fbab7b9d2718b8")
	if err != nil {
		log.Fatal(err)
	}
	_ = client // we'll use this in the upcoming sections


	fmt.Println("Inside function")

	//privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	privateKey, err := crypto.HexToECDSA("E6F9A2469E33F7808666ED49CB88C4AB6637E08AE6ADD0CFAC4681CB0D87B3F7")

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

	address, tx, instance, err := DeployTheContract(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	SetContractAddress(address.Hex())
	var theContract = ERC20Contract{address.Hex()}
	res, err := json.Marshal(theContract)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = instance

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}





