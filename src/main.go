package main

import (
	"./pkg"
	"./utils"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var client = ethclient.Client{}





func handler(w http.ResponseWriter, r *http.Request) {
	title := "Test"

	from := ""
	if r.URL != nil {
		from = r.URL.String()
	}
	if from != "/favicon.ico" {
		log.Printf("title: %s\n", title)
	}

	client, err := ethclient.Dial("https://ropsten.infura.io/v3/511162a74a0c4a80a9fbab7b9d2718b8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	utils.SetClient(*client)
	_ = client // we'll use this in the upcoming sections
}


func main() {
	println("STARTING...")
	//dat, err := os.Getwd()


	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/deployContract", features.DeployContract).Methods("POST")
	r.HandleFunc("/transaction", features.GenerateTransaction).Methods("POST")
	http.ListenAndServe(":8080", r)

}
